package safeclient

import (
	"testing"
)

func TestHTTPClient(t *testing.T) {
	networkPolicy, _ := DefaultNetworkPolicy()
	client := New(networkPolicy, 5)

	if _, err := client.Get("http://localhost"); err == nil {
		t.Errorf("Does not return an error when making requests to localhost")
	}

	if _, err := client.Get("http://aws.interactions.hounder.io"); err == nil {
		t.Errorf("Does not return an error when DNS points to AWS metadata")
	}

	if _, err := client.Get("http://www.google.com"); err != nil {
		t.Errorf("Returns an error when making requests to public addresses")
	}
}
