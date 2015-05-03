package main

import (
	"flag"
	"fmt"
	"github.com/mfycheng/name-dyndns/api"
	"github.com/mfycheng/name-dyndns/dyndns"
	"os"
)

func main() {
	configPath := flag.String("config", "./config.json", "Specify the configuration file")
	daemon := flag.Bool("daemon", false, "operate in daemon mode")
	flag.Parse()

	configs, err := api.LoadConfigs(*configPath)
	if err != nil {
		fmt.Println("Error loading config:", err)
		os.Exit(1)
		return
	}

	fmt.Printf("Successfully loaded %d configs\n", len(configs))
	dyndns.Run(configs, *daemon)
}
