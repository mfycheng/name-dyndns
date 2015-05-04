package api

import (
	"reflect"
	"testing"
)

var expectedConfigs []Config

func init() {
	expectedConfigs = []Config{
		Config{
			Username:  "dev-account",
			Token:     "asdasdasdasdasdad",
			Interval:  60,
			Dev:       true,
			Domain:    "test.com",
			Hostnames: []string{"mail", "chat"},
		},
		Config{
			Username:  "production-account",
			Token:     "123123123123",
			Interval:  3600,
			Domain:    "live.com",
			Hostnames: []string{"mail", "support"},
		},
	}
}

func TestLoadConfigs(t *testing.T) {
	configs, err := LoadConfigs("./config_test.json")

	if err != nil {
		t.Fatalf("Failed to load configs: %s\n", err)
	}

	if !reflect.DeepEqual(expectedConfigs, configs) {
		t.Fatalf("Unexpected configuration")
	}
}
