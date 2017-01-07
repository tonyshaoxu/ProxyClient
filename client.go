package proxyclient

import (
	"net"
	"net/url"
	"strings"
)

type Client interface {
	Dial(network, address string) (net.Conn, error)
	DialTCP(net string, localAddr, remoteAddr *net.TCPAddr) (net.Conn, error)
	DialUDP(net string, localAddr, remoteAddr *net.UDPAddr) (net.Conn, error)
}

type ClientBuilder func(*url.URL) (Client, error)

var schemes = map[string]ClientBuilder{
	// DIRECT
	"direct":    NewDirectProxyClient,
	// SOCKS
	"socks":     NewSocksProxyClient,
	"socks4":    NewSocksProxyClient,
	"socks4a":   NewSocksProxyClient,
	"socks5":    NewSocksProxyClient,
	// HTTP
	"http":      NewHTTPProxyClient,
	"https":     NewHTTPProxyClient,
	// Shadowsocks
	"ss":        NewShadowsocksProxyClient,
	// SSH
	"ssh":       NewSSHAgentProxyClient,
}

func NewProxyClient(proxy string) (Client, error) {
	link, err := url.Parse(proxy)
	if err != nil {
		return nil, err
	}
	link = normalizeLink(*link)
	if factory, ok := schemes[link.Scheme]; ok {
		return factory(link)
	}
	return nil, ErrUnsupportedProtocol
}

func RegisterScheme(schemeName string, factory ClientBuilder) {
	schemes[strings.ToLower(schemeName)] = factory
}

func SupportedSchemes() []string {
	schemeNames := make([]string, 0, len(schemes))
	for schemeName := range schemes {
		schemeNames = append(schemeNames, schemeName)
	}
	return schemeNames
}
