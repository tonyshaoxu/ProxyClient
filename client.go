package proxyclient

import (
	"context"
	"errors"
	"net"
	"net/url"
	"strings"
)

type Dial func(network, address string) (net.Conn, error)

type DialBuilder func(*url.URL, Dial) (Dial, error)

var schemes = map[string]DialBuilder{
	// DIRECT
	"DIRECT":    newDirectProxyClient,
	// REJECT
	"REJECT":    newRejectProxyClient,
	// BLACKHOLE
	"DROP":      newBlackholeProxyClient,
	"BLACKHOLE": newBlackholeProxyClient,
	// SOCKS
	"SOCKS":     newSocksProxyClient,
	"SOCKS4":    newSocksProxyClient,
	"SOCKS4A":   newSocksProxyClient,
	"SOCKS5":    newSocksProxyClient,
	// HTTP
	"HTTP":      newHTTPProxyClient,
	"HTTPS":     newHTTPProxyClient,
}

func NewProxyClient(proxy string) (Dial, error)           { return NewProxyClientWithDial(proxy, net.Dial) }
func NewProxyClientChain(proxies ...string) (Dial, error) { return NewProxyClientChainWithDial(proxies, net.Dial) }

func NewProxyClientWithDial(proxy string, dial Dial) (_ Dial, err error) {
	link, err := url.Parse(proxy)
	if err != nil {
		return
	}
	link = normalizeLink(*link)
	if factory, ok := schemes[link.Scheme]; ok {
		return factory(link, dial)
	}
	err = errors.New("Unsupported proxy client.")
	return
}

func NewProxyClientChainWithDial(proxies []string, upstreamDial Dial) (dial Dial, err error) {
	dial = upstreamDial
	for _, proxy := range proxies {
		dial, err = NewProxyClientWithDial(proxy, dial)
		if err != nil {
			return
		}
	}
	return
}

func RegisterScheme(schemeName string, factory DialBuilder) {
	schemes[strings.ToUpper(schemeName)] = factory
}

func SupportedSchemes() []string {
	schemeNames := make([]string, 0, len(schemes))
	for schemeName := range schemes {
		schemeNames = append(schemeNames, schemeName)
	}
	return schemeNames
}

func (dial Dial) WrappedContext() func(ctx context.Context, network, address string) (net.Conn, error) {
	return func(ctx context.Context, network, address string) (net.Conn, error) {
		return dial(network, address)
	}
}

func (dial Dial) TCPOnly() Dial {
	return func(network, address string) (net.Conn, error) {
		switch strings.ToUpper(network) {
		case "TCP", "TCP4", "TCP6":
			return dial(network, address)
		default:
			return nil, errors.New("Unsupported network type.")
		}
	}
}
