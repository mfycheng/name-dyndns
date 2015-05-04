package api

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Dev       bool     `json:"dev"`
	Domain    string   `json:"domain"`
	Hostnames []string `json:"hostnames"`
	Interval  int      `json:"interval"`
	Token     string   `json:"token"`
	Username  string   `json:"username"`
}

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
