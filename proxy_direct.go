package proxyclient

import (
	"net"
	"net/url"
)

func NewDirectProxyClient(url *url.URL) (Client, error) {
	localAddr := ":0"
	if url != nil {
		localAddr = url.Host
	}
	tcpAddr, err := net.ResolveTCPAddr("tcp", localAddr)
	if err != nil {
		return nil, err
	}
	udpAddr, err := net.ResolveUDPAddr("udp", localAddr)
	if err != nil {
		return nil, err
	}
	client := &directProxyClient{tcpAddr, udpAddr}
	return client, nil
}

type directProxyClient struct {
	tcpLocalAddr *net.TCPAddr
	udpLocalAddr *net.UDPAddr
}

func (client *directProxyClient) Dial(network, address string) (net.Conn, error) {
	if isTCPNetwork(network) {
		addr, err := net.ResolveTCPAddr(network, address)
		if err != nil {
			return nil, err
		}
		return client.DialTCP(network, client.tcpLocalAddr, addr)
	}
	if isUDPNetwork(network) {
		addr, err := net.ResolveUDPAddr(network, address)
		if err != nil {
			return nil, err
		}
		return client.DialUDP(network, client.udpLocalAddr, addr)
	}
	return nil, ErrUnsupportedNetwork
}

func (client *directProxyClient) DialTCP(network string, localAddr, remoteAddr *net.TCPAddr) (net.Conn, error) {
	if localAddr == nil {
		localAddr = client.tcpLocalAddr
	}
	return net.DialTCP(network, localAddr, remoteAddr)
}

func (client *directProxyClient) DialUDP(network string, localAddr, remoteAddr *net.UDPAddr) (net.Conn, error) {
	if localAddr == nil {
		localAddr = client.udpLocalAddr
	}
	return net.DialUDP(network, localAddr, remoteAddr)
}
