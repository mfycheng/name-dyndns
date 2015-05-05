package main

import (
	"flag"
	"fmt"
	"github.com/mfycheng/name-dyndns/api"
	"github.com/mfycheng/name-dyndns/dyndns"
	"github.com/mfycheng/name-dyndns/log"
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
	logFile := flag.String("log", "", "Specify a logfile. If no file is provided, uses stdout.")
	flag.Parse()

	var file *os.File
	defer file.Close()

	if *logFile == "" {
		file = os.Stdout
	} else {
		var err error
		file, err = os.OpenFile(*logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			fmt.Println("Could not open log for reading")
			os.Exit(1)
		}
	}

	log.Init(file)

	configs, err := api.LoadConfigs(*configPath)
	if err != nil {
		log.Logger.Fatalln("Error loading config:", err)
	}

	for _, config := range configs {
		if config.Domain == "" || len(config.Hostnames) == 0 {
			log.Logger.Fatalf("Empty configuration detected. Exiting.")
		}
	}

	configs = filterConfigs(configs, *dev)

	log.Logger.Printf("Successfully loaded %d configs\n", len(configs))
	dyndns.Run(configs, *daemon)
}
