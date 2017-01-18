// Package dyndns provides a tool for running a
// dynamic dns updating service.
package dyndns

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/mfycheng/name-dyndns/api"
	"github.com/mfycheng/name-dyndns/log"
)

var wg sync.WaitGroup

func contains(c api.Config, val string) bool {
	for _, v := range c.Hostnames {
		// We have a special case where an empty hostname
		// is equivalent to the domain (i.e. val == domain).
		if val == c.Domain && v == "" {
			return true
		} else if fmt.Sprintf("%s.%s", v, c.Domain) == val {
			return true
		}
	}
	return false
}

func updateDNSRecord(a api.API, domain, recordID string, newRecord api.DNSRecord) error {
	log.Logger.Printf("Deleting DNS record for %s.\n", newRecord.Name)
	err := a.DeleteDNSRecord(domain, newRecord.RecordID)
	if err != nil {
		return err
	}

	log.Logger.Printf("Creating DNS record for %s: %s\n", newRecord.Name, newRecord)

	// Remove the domain from the DNSRecord name.
	// This is an unfortunate inconsistency from the API
	// implementation (returns full name, but only requires host)
	if newRecord.Name == domain {
		newRecord.Name = ""
	} else {
		newRecord.Name = strings.TrimSuffix(newRecord.Name, fmt.Sprintf(".%s", domain))
	}

	return a.CreateDNSRecord(domain, newRecord)
}

func runConfig(c api.Config, daemon bool) {
	defer wg.Done()

	a := api.NewAPIFromConfig(c)
	for {
		ip, err := GetExternalIP()
		if err != nil {
			log.Logger.Print("Failed to retreive IP: ")
			if daemon {
				log.Logger.Printf("Will retry in %d seconds...\n", c.Interval)
				time.Sleep(time.Duration(c.Interval) * time.Second)
				continue
			} else {
				log.Logger.Println("Giving up.")
				break
			}
		}

		// GetRecords retrieves a list of DNSRecords,
		// 1 per hostname with the associated domain.
		// If the content is not the current IP, then
		// update it.
		records, err := a.GetDNSRecords(c.Domain)
		if err != nil {
			log.Logger.Printf("Failed to retreive records for %s:\n\t%s\n", c.Domain, err)
			if daemon {
				log.Logger.Printf("Will retry in %d seconds...\n", c.Interval)
				time.Sleep(time.Duration(c.Interval) * time.Second)
				continue
			} else {
				log.Logger.Print("Giving up.")
				break
			}
		}

		for _, r := range records {
			if !contains(c, r.Name) {
				continue
			}

			// Only A records should be mapped to an IP.
			// TODO: Support AAAA records.
			if r.Type != "A" {
				continue
			}

			log.Logger.Printf("Running update check for %s.", r.Name)
			if r.Content != ip {
				r.Content = ip
				err = updateDNSRecord(a, c.Domain, r.RecordID, r)
				if err != nil {
					log.Logger.Printf("Failed to update record %s [%s] with IP: %s\n\t%s\n", r.RecordID, r.Name, ip, err)
				} else {
					log.Logger.Printf("Updated record %s [%s] with IP: %s\n", r.RecordID, r.Name, ip)
				}
			}
		}

		log.Logger.Println("Update complete.")
		if !daemon {
			log.Logger.Println("Non daemon mode, stopping.")
			return
		}
		log.Logger.Printf("Will update again in %d seconds.\n", c.Interval)

		time.Sleep(time.Duration(c.Interval) * time.Second)
	}
}

// Run will process each configuration in configs.
// If daemon is true, then Run will run forever,
// processing each configuration at its specified
// interval.
//
// Each configuration represents a domain with
// multiple hostnames.
func Run(configs []api.Config, daemon bool) {
	for _, config := range configs {
		wg.Add(1)
		go runConfig(config, daemon)
	}

	wg.Wait()
}
