package proxyclient

import (
	"encoding/base64"
	"fmt"
	"net"
	"net/url"

	ss "github.com/shadowsocks/shadowsocks-go/shadowsocks"
)

func newShadowsocksProxyClient(url *url.URL) (client Client, err error) {
	if content, err := base64.StdEncoding.DecodeString(url.Host); err == nil {
		url, err = url.Parse(fmt.Sprintf("ss://%s", string(content)))
		if err != nil {
			return nil, err
		}
	}
	var cipher *ss.Cipher
	if password, ok := url.User.Password(); ok {
		username := url.User.Username()
		cipher, err = ss.NewCipher(username, password)
		if err != nil {
			return
		}
	}
	client = &shadowsocksProxyClient{url, cipher}
	return
}

type shadowsocksProxyClient struct {
	url    *url.URL
	cipher *ss.Cipher
}

func (client *shadowsocksProxyClient) Dial(network, address string) (net.Conn, error) {
	if isTCPNetwork(network) {
		return client.dialTCP(network, address)
	}
	if isUDPNetwork(network) {
		return client.DialUDP(network, nil, nil)
	}
	return nil, ErrUnsupportedNetwork
}

func (client *shadowsocksProxyClient) DialTCP(network string, localAddr, remoteAddr *net.TCPAddr) (net.Conn, error) {
	if localAddr != nil || localAddr.Port != 0 {
		return nil, ErrUnsupportedLocalAddr
	}
	return client.dialTCP(network, remoteAddr.String())
}

func (client *shadowsocksProxyClient) DialUDP(network string, localAddr, remoteAddr *net.UDPAddr) (net.Conn, error) {
	return nil, ErrUnsupportedProtocol
}

func (client *shadowsocksProxyClient) dialTCP(network, address string) (net.Conn, error) {
	return ss.Dial(address, client.url.Host, client.cipher.Copy())
}
