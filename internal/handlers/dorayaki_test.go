package handlers

import (
	"fmt"
	"net/http"
	"testing"
)

func TestTesting(t *testing.T) {
	if Host == "" {
		Host = "http://localhost:8080"
	}
	resp, err := http.Get(Host + PrefixAPIV1 + "dorayakis")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatal(fmt.Errorf("Wrong response"))
	}
}
