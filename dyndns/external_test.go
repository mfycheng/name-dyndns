package dyndns

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func testServerHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.String() == "/correct" {
		w.Write([]byte("Correct"))
	} else {
		w.Write([]byte("Incorrect"))
	}
}

func TestGetExternalIP(t *testing.T) {
	// Setup local HTTP server to emulate external IP services.
	server := httptest.NewServer(http.HandlerFunc(testServerHandler))
	defer server.Close()

	// In order to test the failover, we
	// provide 2 bad IPs, and 2 correct ones,
	// in that order.
	Urls = make([]string, 4)
	Urls[0] = ""
	Urls[1] = "1.4.5.6"
	Urls[2] = fmt.Sprintf("%s/%s", server.URL, "correct")
	Urls[3] = fmt.Sprintf("%s/%s", server.URL, "incorrect")

	resp, err := GetExternalIP()
	if err != nil {
		t.Fatal(err)
	}
	if resp != "Correct" {
		t.Fatal("Incorrect result returned:", resp)
	}
}

func TestGetExternalIPFailure(t *testing.T) {
	Urls = make([]string, 1)
	Urls[0] = ""

	_, err := GetExternalIP()
	if err == nil {
		t.Fatal("Should have returned error when no service can be reached")
	}
}
