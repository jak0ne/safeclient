package main

import (
	"log"
	"os"

	"github.com/jak0ne/safeclient"
)

func main() {
	networkPolicy, err := safeclient.DefaultNetworkPolicy()
	if err != nil {
		log.Fatalf("Could not create network policy: %v", err)
	}

	safeClient := safeclient.New(networkPolicy, 5)

	resp, err := safeClient.Get(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%d\n", resp.StatusCode)
}
