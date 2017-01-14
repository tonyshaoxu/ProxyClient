package proxyclient

import (
	"crypto/tls"
	"errors"
	"net"
	"net/url"
	"time"

	httpProxy "github.com/RouterScript/HTTPProxy"
	socksProxy "github.com/RouterScript/SOCKSProxy"
)

func newDirectProxyClient(proxy *url.URL, _ Dial) (dial Dial, err error) {
	dial = net.Dial
	if timeout := proxy.Query().Get("timeout"); timeout != "" {
		dialTimeout, err := time.ParseDuration(timeout)
		if err != nil {
			return nil, err
		}
		dial = DialWithTimeout(dialTimeout)
	}
	return
}

func newRejectProxyClient(_ *url.URL, _ Dial) (dial Dial, err error) {
	dial = func(network, address string) (net.Conn, error) {
		return nil, errors.New("reject dial")
	}
	return
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
