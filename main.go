package main

import (
	"fmt"
	"github.com/mfycheng/name-dyndns/api"
	"os"
)

const (
	defaultConfigPath = "./config.json"
)

func main() {
	configs, err := api.LoadConfigs(defaultConfigPath)
	if err != nil {
		fmt.Println("Error loading config:", err)
		os.Exit(1)
		return
	}

	fmt.Printf("Successfully loaded %d configs\n", len(configs))
	for _, config := range configs {
		fmt.Println("Processing config", config)
	}
}
