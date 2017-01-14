package proxyclient

import (
	"net"
	"net/url"
	"strings"
	"time"
)

func limitSchemes(proxy *url.URL, names ...string) bool {
	if proxy == nil {
		return false
	}
	schemeName := strings.ToUpper(proxy.Scheme)
	for _, name := range names {
		if strings.EqualFold(schemeName, strings.ToUpper(name)) {
			return true
		}
	}
	return false
}

func normalizeLink(link url.URL) *url.URL {
	switch strings.ToUpper(link.Path) {
	case "DIRECT", "REJECT", "DROP", "BLACKHOLE":
		link.Scheme = link.Path
		link.Path = ""
	}
	link.Scheme = strings.ToUpper(link.Scheme)
	query := link.Query()
	for name, value := range query {
		query[strings.ToLower(name)] = value
	}
	return &link
}

func DialWithTimeout(timeout time.Duration) Dial {
	return func(network, address string) (net.Conn, error) {
		return net.DialTimeout(network, address, timeout)
	}
}
