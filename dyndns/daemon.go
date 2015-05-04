package dyndns

import (
	"fmt"
	"github.com/mfycheng/name-dyndns/api"
	"sync"
	"time"
)

var wg sync.WaitGroup

func updateDNSRecord(a api.API, domain, recordId string, newRecord api.DNSRecord) error {
	fmt.Println("Deleting record...")
	err := a.DeleteDNSRecord(domain, newRecord.RecordId)
	if err != nil {
		return err
	}

	// Does a /create/ overwrite? or do we have to delete first?
	fmt.Println("Creating record")
	return a.CreateDNSRecord(domain, newRecord)
}

func runConfig(c api.Config, daemon bool) {
	defer wg.Done()

	a := api.NewAPIFromConfig(c)
	for {
		ip, err := GetExternalIP()
		if err != nil {
			fmt.Print("Fail to retreive IP: ")
			if daemon {
				fmt.Println("Will retry...")
				continue
			} else {
				fmt.Println("Giving up.")
				break
			}
		}

		// GetRecords retrieves a list of DNSRecords,
		// 1 per hostname with the associated domain.
		// If the content is not the current IP, then
		// update it.
		records, err := a.GetRecords(c.Domain)
		if err != nil {
			fmt.Print("Failed to retreive records for:%s\n", c.Domain)
			if daemon {
				fmt.Println("Will retry...")
				continue
			} else {
				fmt.Println("Giving up.")
				break
			}
		}

		for _, r := range records {
			if r.Content != ip {
				fmt.Printf("Updating record %s for %s\n", r.RecordId, c.Domain)
				r.Content = ip
				err = updateDNSRecord(a, c.Domain, r.RecordId, r)
				if err != nil {
					fmt.Println("Failed to update record.", err)
				}
			}
		}

		if !daemon {
			return
		}

		time.Sleep(time.Duration(c.Interval) * time.Second)
	}
}

// For each domain, check if the host record matches
// the current external IP. If it does not, it updates.
// If daemon is true, then Run will run forever, polling at
// an interval specified in each config.
func Run(configs []api.Config, daemon bool) {
	for _, config := range configs {
		wg.Add(1)
		go runConfig(config, daemon)
	}

	wg.Wait()
}
