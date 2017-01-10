package proxyclient

import (
	"errors"
	"net"
	"net/url"
	"strings"
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

func dialTCPOnly(dial Dial) Dial {
	return func(network, address string) (net.Conn, error) {
		switch strings.ToUpper(network) {
		case "TCP", "TCP4", "TCP6":
			return dial(network, address)
		default:
			return nil, errors.New("Unsupported network type.")
		}
	}
}

func normalizeLink(link url.URL) *url.URL {
	if strings.ToUpper(link.Path) == "DIRECT" {
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
