# safeclient

Minimal base Go HTTP client to interact with untrusted URLs in a "safe" manner preventing SSRF attacks.

Uses [networkpolicy](https://github.com/projectdiscovery/networkpolicy) for enforcing allow/deny lists of hosts, network ranges, ports and URL schemes.

Features:

- Fully configurable allow/deny lists with [networkpolicy](https://github.com/projectdiscovery/networkpolicy).
- IPv4 and IPv6 support.
- Checks against malicious URL redirects.
- Protects from DNS rebinding (TOC/TOU) by hooking the HTTP client's transport dial context.
- Protects against alternate encodings such as 0x7F000001. 

## Usage

Example - Create a client that denies connections to localhost only
```go
package main

import (
	"log"

	"github.com/projectdiscovery/networkpolicy"
	"github.com/jak0ne/safeclient"
)

func main() {
	var npOptions networkpolicy.Options
	npOptions.DenyList = append(npOptions.DenyList, "127.0.0.0/8")

	networkPolicy, err := networkpolicy.New(npOptions)
	if err != nil {
		log.Fatal(err)
	}

	safeClient := safeclient.New(networkPolicy, 5)

	if _, err = safeClient.Get("https://127.0.0.1"); err != nil {
		log.Fatal(err)
	}
}
```

Example - Using default network policy
```go
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
```

Example - Override HTTP client options, disable certificate validation
```go
package main

import (
	"log"

	"github.com/jak0ne/safeclient"
)

func main() {
	networkPolicy, err := safeclient.DefaultNetworkPolicy()
	if err != nil {
		log.Fatalf("Could not create network policy: %v", err)
	}

	safeClient := safeclient.New(networkPolicy, 5)
	safeClient.Transport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	if resp, err = safeClient.Get("https://untrusted-url"); err != nil {
		log.Fatal(err)
	}

	log.Printf("%+v\n", resp)
}
```
