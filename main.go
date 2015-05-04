package main

import (
	"flag"
	"fmt"
	"github.com/mfycheng/name-dyndns/api"
	"github.com/mfycheng/name-dyndns/dyndns"
	"os"
)

func filterConfigs(configs []api.Config, dev bool) []api.Config {
	for i := 0; i < len(configs); i++ {
		if configs[i].Dev != dev {
			configs = append(configs[:i], configs[i+1:]...)
		}
	}

	return configs
}

func main() {
	configPath := flag.String("config", "./config.json", "Specify the configuration file")
	daemon := flag.Bool("daemon", false, "Operate in daemon mode.")
	dev := flag.Bool("dev", false, "Use development configurations instead.")
	flag.Parse()

	configs, err := api.LoadConfigs(*configPath)
	if err != nil {
		fmt.Println("Error loading config:", err)
		os.Exit(1)
		return
	}

	configs = filterConfigs(configs, *dev)

	fmt.Printf("Successfully loaded %d configs\n", len(configs))
	dyndns.Run(configs, *daemon)
}
