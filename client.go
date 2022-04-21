package safeclient

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/projectdiscovery/networkpolicy"
)

func DefaultNetworkPolicy() (*networkpolicy.NetworkPolicy, error) {
	var npOptions networkpolicy.Options
	// https://github.com/projectdiscovery/networkpolicy/blob/main/cidr.go
	npOptions.DenyList = append(npOptions.DenyList, networkpolicy.DefaultIPv4DenylistRanges...)
	npOptions.DenyList = append(npOptions.DenyList, networkpolicy.DefaultIPv6DenylistRanges...)
	// Allow only HTTP and HTTPS schemes
	npOptions.AllowSchemeList = networkpolicy.DefaultSchemeAllowList
	np, err := networkpolicy.New(npOptions)
	if err != nil {
		return nil, err
	}

	return np, nil
}

func New(np *networkpolicy.NetworkPolicy, maxRedirects uint) *http.Client {
	dialer := &net.Dialer{
		Timeout: 15 * time.Second,
	}
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		host, port, err := net.SplitHostPort(addr)
		if err != nil {
			return nil, err
		}

		ip, ok := np.ValidateHost(host)
		if !ok {
			return nil, errors.New("Could not find a valid IP to dial to")
		}

		addr = net.JoinHostPort(ip, port)
		con, err := dialer.DialContext(ctx, network, addr)
		if err != nil {
			return nil, err
		}

		return con, nil
	}

	return &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= int(maxRedirects) {
				return errors.New("Too many redirects")
			}

			if req.Response.StatusCode == 307 || req.Response.StatusCode == 308 {
				return errors.New("Unsupported redirect (Status code 307/308)")
			}

			_, ok := np.ValidateHost(req.URL.Host)
			if req.URL.Host != "" && !ok {
				return errors.New("Redirects to forbidden target")
			}

			return nil
		},
		Transport: transport,
	}
}
