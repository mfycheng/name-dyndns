// Package api provides a basic interface for dealing
// with Name.com DNS API's.
package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

const (
	productionURL = "https://api.name.com/"
	devURL        = "https://api.dev.name.com/"
)

// API Contains details required to access the Name.com API.
type API struct {
	baseURL  string
	username string
	token    string
}

// DNSRecord contains information about a Name.com DNS record.
type DNSRecord struct {
	RecordID   string `json:"record_id"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	Content    string `json:"content"`
	TTL        string `json:"ttl"`
	CreateDate string `json:"create_date"`
}

type resultResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// NewNameAPI constructs a new Name.com API. If dev is true, then
// the API uses the development API, instead of the production API.
func NewNameAPI(username, token string, dev bool) API {
	a := API{username: username, token: token}

	if dev {
		a.baseURL = devURL
	} else {
		a.baseURL = productionURL
	}

	return a
}

// NewAPIFromConfig constructs a new Name.com API from a configuration.
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

// CreateDNSRecord creates a DNS record for a given domain. The name
// field in DNSRecord is in the format [hostname].[domainname]
func (api API) CreateDNSRecord(domain string, record DNSRecord) error {
	// We need to transform name -> hostname for JSON.
	var body struct {
		Hostname string `json:"hostname"`
		Type     string `json:"type"`
		Content  string `json:"content"`
		TTL      string `json:"ttl"`
	}

	body.Hostname = record.Name
	body.Type = record.Type
	body.Content = record.Content
	body.TTL = record.TTL

	b, err := json.Marshal(body)
	if err != nil {
		return err
	}

	resp, err := api.performRequest(
		"POST",
		fmt.Sprintf("%s%s%s", api.baseURL, "api/dns/create/", domain),
		bytes.NewBuffer(b),
	)
	if err != nil {
		return err
	}

	var result struct {
		Result resultResponse
	}

	err = json.Unmarshal(resp, &result)
	if err != nil {
		return err
	}
	if result.Result.Code != 100 {
		return errors.New(result.Result.Message)
	}

	return nil
}

// DeleteDNSRecord deletes a DNS record for a given domain. The recordID can
// be retreived from GetDNSRecords.
func (api API) DeleteDNSRecord(domain, recordID string) error {
	var body struct {
		RecordID string `json:"record_id"`
	}
	body.RecordID = recordID

	b, err := json.Marshal(body)
	if err != nil {
		return err
	}

	resp, err := api.performRequest(
		"POST",
		fmt.Sprintf("%s%s%s", api.baseURL, "api/dns/delete/", domain),
		bytes.NewBuffer(b),
	)
	if err != nil {
		return err
	}

	var result struct {
		Result resultResponse
	}

	err = json.Unmarshal(resp, &result)
	if err != nil {
		return err
	}
	if result.Result.Code != 100 {
		return errors.New(result.Result.Message)
	}

	return nil
}

// GetDNSRecords returns a slice of DNS records associated with a given domain.
func (api API) GetDNSRecords(domain string) (records []DNSRecord, err error) {
	resp, err := api.performRequest(
		"GET",
		fmt.Sprintf("%s%s%s", api.baseURL, "api/dns/list/", domain),
		nil,
	)

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
	if result.Result.Code != 100 {
		return nil, errors.New(result.Result.Message)
	}

	return result.Records, nil
}
