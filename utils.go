package proxyclient

import (
	"crypto/tls"
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
	dialer := net.Dialer{Timeout: timeout}
	return dialer.Dial
}

func tlsConfigByProxyURL(proxy *url.URL) *tls.Config {
	query := proxy.Query()
	conf := &tls.Config{
		ServerName:         query.Get("tls-domain"),
		InsecureSkipVerify: query.Get("tls-insecure-skip-verify") == "true",
	}
	if conf.ServerName == "" {
		conf.ServerName = proxy.Host
	}
	return conf
}
