package proxyclient

import (
	"crypto/tls"
	"net"
	"net/url"
	"github.com/RouterScript/HTTPProxy"
)

func newHTTPProxyClient(url *url.URL) (Client, error) {
	client := proxy.Client{
		Proxy:       *url,
		TLSConfig:   &tls.Config{
			ServerName:        url.Query().Get("tls-domain"),
			InsecureSkipVerify:url.Query().Get("tls-insecure-skip-verify") == "true",
		},
		UpstreamDial:net.Dial,
	}
	return &httpProxyClient{client}, nil
}

type httpProxyClient struct {
	client proxy.Client
}

func (c *httpProxyClient) Dial(network, address string) (net.Conn, error) {
	if isTCPNetwork(network) {
		return c.client.Dial(network, address)
	}
	if isUDPNetwork(network) {
		return c.DialUDP(network, nil, nil)
	}
	return nil, ErrUnsupportedNetwork
}

func (c *httpProxyClient) DialTCP(network string, localAddr, remoteAddr *net.TCPAddr) (net.Conn, error) {
	if localAddr != nil || localAddr.Port != 0 {
		return nil, ErrUnsupportedLocalAddr
	}
	return c.client.Dial(network, remoteAddr.String())
}

func (c *httpProxyClient) DialUDP(network string, localAddr, remoteAddr *net.UDPAddr) (net.Conn, error) {
	return nil, ErrUnsupportedProtocol
}
