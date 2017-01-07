package proxyclient

import (
	"net"
	"net/url"
	"strconv"

	"github.com/fangdingjun/socks-go"
)

type socksProxyClient struct {
	url *url.URL
}

func NewSocksProxyClient(url *url.URL) (Client, error) {
	client := &socksProxyClient{url}
	return client, nil
}

func (client *socksProxyClient) Dial(network, address string) (net.Conn, error) {
	if isTCPNetwork(network) {
		return client.dialTCP(network, address)
	}
	if isUDPNetwork(network) {
		return client.DialUDP(network, nil, nil)
	}
	return nil, ErrUnsupportedNetwork
}

func (client *socksProxyClient) DialTCP(network string, localAddr, remoteAddr *net.TCPAddr) (net.Conn, error) {
	if localAddr != nil || localAddr.Port != 0 {
		return nil, ErrUnsupportedLocalAddr
	}
	return client.dialTCP(network, remoteAddr.String())
}

func (client *socksProxyClient) DialUDP(network string, localAddr, remoteAddr *net.UDPAddr) (net.Conn, error) {
	return nil, ErrUnsupportedProtocol
}

func (client *socksProxyClient) dialTCP(network, address string) (net.Conn, error) {
	conn, err := net.Dial(network, address)
	if err != nil {
		return nil, err
	}
	username := client.url.User.Username()
	password, _ := client.url.User.Password()
	socksClient := &socks.Client{Conn:conn, Username:username, Password:password}
	hostname, hostPort, err := net.SplitHostPort(address)
	if err != nil {
		return nil, err
	}
	port, err := strconv.Atoi(hostPort)
	if err != nil {
		return nil, err
	}
	if err := socksClient.Connect(hostname, uint16(port)); err != nil {
		return nil, err
	}
	return socksClient, nil
}
