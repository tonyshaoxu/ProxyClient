package proxyclient

import (
	"crypto/tls"
	"errors"
	"net"
	"net/url"

	httpProxy "github.com/RouterScript/HTTPProxy"
	socksProxy "github.com/RouterScript/SOCKSProxy"
)

func newDirectProxyClient(_ *url.URL, _ Dial) (Dial, error) {
	return net.Dial, nil
}

func newRejectProxyClient(_ *url.URL, _ Dial) (Dial, error) {
	dial := func(network, address string) (net.Conn, error) {
		return nil, errors.New("reject dial")
	}
	return dial, nil
}

func newHTTPProxyClient(proxy *url.URL, upstreamDial Dial) (dial Dial, err error) {
	client := httpProxy.Client{
		Proxy:       *proxy,
		TLSConfig:   &tls.Config{
			ServerName:        proxy.Query().Get("tls-domain"),
			InsecureSkipVerify:proxy.Query().Get("tls-insecure-skip-verify") == "true",
		},
		UpstreamDial:upstreamDial,
	}
	if client.TLSConfig.ServerName == "" {
		client.TLSConfig.ServerName = proxy.Host
	}
	dial = client.Dial
	dial = dial.TCPOnly()
	return
}

func newSocksProxyClient(proxy *url.URL, upstreamDial Dial) (dial Dial, err error) {
	client, err := socksProxy.NewClient(proxy, &socksProxy.SOCKSConf{Dial:upstreamDial})
	if err != nil {
		return
	}
	dial = client.Dial
	dial = dial.TCPOnly()
	return
}
