package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

const (
	productionUrl = "https://api.name.com/"
	devUrl        = "https://api.dev.name.com/"
)

type API struct {
	baseUrl  string
	username string
	token    string
}

type DNSRecord struct {
	RecordId   string `json:"record_id"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	Content    string `json:"content"`
	Ttl        string `json:"ttl"`
	CreateDate string `json:"create_date"`
}

type resultResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Constructs a new Name.com API
//
// If dev = true, then the API uses the development
// API, instead of the production API.
func NewNameAPI(username, token string, dev bool) API {
	a := API{username: username, token: token}

	if dev {
		a.baseUrl = devUrl
	} else {
		a.baseUrl = productionUrl
	}

	return a
}

// Constructs a new Name.com API from a configuration
func NewAPIFromConfig(c Config) API {
	return NewNameAPI(c.Username, c.Token, c.Dev)
}

func (api API) performRequest(method, url string, body io.Reader) (response []byte, err error) {
	var client http.Client
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("api-username", api.username)
	req.Header.Add("api-token", api.token)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

// Returns a slice of DNS records associated with a given hostname.
func (api API) GetRecords(hostname string) (records []DNSRecord, err error) {
	resp, err := api.performRequest("GET", fmt.Sprintf("%s%s%s", api.baseUrl, "api/dns/list/", hostname), nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Result  resultResponse
		Records []DNSRecord
	}

	err = json.Unmarshal(resp, &result)
	if err != nil {
		return nil, err
	}

	// Make sure the API call was successful
	if result.Result.Code != 100 {
		return nil, errors.New(result.Result.Message)
	}

	return result.Records, nil
}
