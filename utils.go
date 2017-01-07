package proxyclient

import (
	"net/url"
	"strings"
)

func isTCPNetwork(network string) bool {
	switch strings.ToLower(network) {
	case "tcp", "tcp4", "tcp6":
		return true
	}
	return false
}

func isUDPNetwork(network string) bool {
	switch strings.ToLower(network) {
	case "udp", "udp4", "udp6":
		return true
	}
	return false
}

func normalizeLink(link url.URL) *url.URL {
	if link.Path == "direct" {
		link.Scheme = link.Path
		link.Path = ""
	}
	link.Scheme = strings.ToLower(link.Scheme)
	query := link.Query()
	for name, value := range query {
		query[strings.ToLower(name)] = value
	}
	return &link
}
