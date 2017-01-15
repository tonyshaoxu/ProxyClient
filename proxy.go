package proxyclient

import (
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
		Proxy:        *proxy,
		TLSConfig:    tlsConfigByProxyURL(proxy),
		UpstreamDial: upstreamDial,
	}
	dial = Dial(client.Dial).TCPOnly()
	return
}

func newSocksProxyClient(proxy *url.URL, upstreamDial Dial) (dial Dial, err error) {
	conf := &socksProxy.SOCKSConf{
		TLSConfig: tlsConfigByProxyURL(proxy),
		Dial:      upstreamDial,
	}
	client, err := socksProxy.NewClient(proxy, conf)
	if err != nil {
		return
	}
	dial = Dial(client.Dial).TCPOnly()
	return
}
