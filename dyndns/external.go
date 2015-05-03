package dyndns

import (
	"errors"
	"io/ioutil"
	"net/http"
)

var (
	Urls = []string{"http://myexternalip.com/raw"}
)

func tryMirror(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(contents), nil
}

// Retrieves the external facing IP Address
func GetExternalIP() (string, error) {
	for _, url := range Urls {
		resp, err := tryMirror(url)
		if err == nil {
			return resp, err
		}
	}

	return "", errors.New("Could not retreive external IP")
}
