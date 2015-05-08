package api

import (
	"encoding/json"
	"io/ioutil"
)

// Config represents the configuration for a
// specific domain. Each domain can have multiple
// hostnames, including the root domain, where
// hostname is an empty string.
//
// The interval is the polling time (in seconds) for
// daemon mode.
type Config struct {
	Dev       bool     `json:"dev"`
	Domain    string   `json:"domain"`
	Hostnames []string `json:"hostnames"`
	Interval  int      `json:"interval"`
	Token     string   `json:"token"`
	Username  string   `json:"username"`
}

// LoadConfigs loads configurations from a file. The configuration
// is stored as an array of JSON serialized Config structs.
func LoadConfigs(path string) ([]Config, error) {
	var configs struct {
		Configs []Config
	}

	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(file, &configs)
	if err != nil {
		return nil, err
	}

	return configs.Configs, nil
}
