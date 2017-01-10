package proxyclient

import (
	"crypto/tls"
	"net"
	"net/url"

	httpProxy "github.com/RouterScript/HTTPProxy"
	socksProxy "github.com/RouterScript/SOCKSProxy"
)

func newDirectProxyClient(_ *url.URL, _ Dial) (Dial, error) {
	return net.Dial, nil
}

func newHTTPProxyClient(proxy *url.URL, upstreamDial Dial) (Dial, error) {
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
	return dialTCPOnly(client.Dial), nil
}

func newSocksProxyClient(proxy *url.URL, upstreamDial Dial) (dial Dial, err error) {
	client, err := socksProxy.NewClient(proxy, &socksProxy.SOCKSConf{Dial:upstreamDial})
	if err != nil {
		return
	}
	dial = dialTCPOnly(client.Dial)
	return
}
