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
	productionUrl = "https://api.name.com/"
	devUrl        = "https://api.dev.name.com/"
)

// Contains details required to access the Name.com API.
type API struct {
	baseUrl  string
	username string
	token    string
}

// A Name.com DNS record.
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

// Constructs a new Name.com API. If dev is true, then
// the API uses the development API, instead of the production API.
func NewNameAPI(username, token string, dev bool) API {
	a := API{username: username, token: token}

	if dev {
		a.baseUrl = devUrl
	} else {
		a.baseUrl = productionUrl
	}

	return a
}

// Constructs a new Name.com API from a configuration.
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

// Create a DNS record for a given domain. The hostname is
// specified in the DNSRecord under name, and should not
// include the domain.
func (api API) CreateDNSRecord(domain string, record DNSRecord) error {
	b, err := json.Marshal(record)
	if err != nil {
		return err
	}

	resp, err := api.performRequest(
		"POST",
		fmt.Sprintf("%s%s%s", api.baseUrl, "api/dns/create/", domain),
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

// Deletes a DNS record for a given domain. The recordId can
// be retreived from GetDNSRecords.
func (api API) DeleteDNSRecord(domain, recordId string) error {
	var body struct {
		RecordId string `json:"record_id"`
	}
	body.RecordId = recordId

	b, err := json.Marshal(body)
	if err != nil {
		return err
	}

	resp, err := api.performRequest(
		"POST",
		fmt.Sprintf("%s%s%s", api.baseUrl, "api/dns/delete/", domain),
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

// Returns a slice of DNS records associated with a given domain.
func (api API) GetDNSRecords(domain string) (records []DNSRecord, err error) {
	resp, err := api.performRequest(
		"GET",
		fmt.Sprintf("%s%s%s", api.baseUrl, "api/dns/list/", domain),
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
