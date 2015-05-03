package dyndns

import (
	"fmt"
	"github.com/mfycheng/name-dyndns/api"
	"sync"
	"time"
)

var wg sync.WaitGroup

func updateDomain(a api.API, currentIP, domain string) error {
	records, err := a.GetRecords(domain)
	if err != nil {
		return err
	}

	if len(records) == 0 {
		return nil
	}

	for _, record := range records {
		if record.Content != currentIP {
			// TODO: Update
		}
	}

	return nil
}

func runConfig(c api.Config, daemon bool) {
	defer wg.Done()

	a := api.NewAPIFromConfig(c)
	for {
		ip, err := GetExternalIP()
		if err != nil {
			fmt.Print("Failed to retreive IP: ")
			if daemon {
				fmt.Println("Will retry...")
				continue
			} else {
				fmt.Println("Giving up.")
				break
			}
		}

		for _, domain := range c.Domains {
			updateDomain(a, ip, domain)
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
