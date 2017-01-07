package proxyclient

import (
	"bufio"
	"crypto/tls"
	"net"
	"net/http"
	"net/url"
)

func newHTTPProxyClient(url *url.URL) (Client, error) {
	tlsConfig := &tls.Config{
		ServerName:        url.Query().Get("tls-domain"),
		InsecureSkipVerify:url.Query().Get("tls-insecure-skip-verify") == "true",
	}
	client := &httpProxyClient{url, tlsConfig}
	return client, nil
}

type httpProxyClient struct {
	url       *url.URL
	tlsConfig *tls.Config
}

func (client *httpProxyClient) Dial(network, address string) (net.Conn, error) {
	if isTCPNetwork(network) {
		return client.dialTCP(network, address)
	}
	if isUDPNetwork(network) {
		return client.DialUDP(network, nil, nil)
	}
	return nil, ErrUnsupportedNetwork
}

func (client *httpProxyClient) DialTCP(network string, localAddr, remoteAddr *net.TCPAddr) (net.Conn, error) {
	if localAddr != nil || localAddr.Port != 0 {
		return nil, ErrUnsupportedLocalAddr
	}
	return client.dialTCP(network, remoteAddr.String())
}

func (client *httpProxyClient) DialUDP(network string, localAddr, remoteAddr *net.UDPAddr) (net.Conn, error) {
	return nil, ErrUnsupportedProtocol
}

func (client *httpProxyClient) dialTCP(network, address string) (net.Conn, error) {
	conn, err := net.Dial(network, client.url.Host)
	if err != nil {
		return nil, err
	}
	if client.url.Scheme == "https" {
		tlsConn := tls.Client(conn, client.tlsConfig)
		if err := tlsConn.Handshake(); err != nil {
			tlsConn.Close()
			return nil, err
		}
		verifyHostname := tlsConn.VerifyHostname(client.tlsConfig.ServerName)
		verifyInsecureSkip := client.tlsConfig.InsecureSkipVerify
		if !verifyInsecureSkip && verifyHostname != nil {
			tlsConn.Close()
			return nil, verifyHostname
		}
		conn = tlsConn
	}
	remoteHost, _, err := net.SplitHostPort(address)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest(http.MethodConnect, remoteHost, nil)
	if err != nil {
		conn.Close()
		return nil, err
	}
	request.URL.Host = address
	request.Host = address
	if err := request.Write(conn); err != nil {
		conn.Close()
		return nil, err
	}
	reader := bufio.NewReader(conn)
	response, err := http.ReadResponse(reader, request)
	if err != nil {
		conn.Close()
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		conn.Close()
	}
	return conn, nil
}
