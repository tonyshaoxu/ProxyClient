package proxyclient

import (
	"io/ioutil"
	"net"
	"net/url"
	"golang.org/x/crypto/ssh"
)

type sshagentProxyClient struct {
	*ssh.Client
}

func newSSHAgentProxyClient(url *url.URL) (Client, error) {
	conf := &ssh.ClientConfig{
		User: url.User.Username(),
		Auth: sshagentAuth(url),
	}
	sshClient, err := ssh.Dial("tcp", url.Host, conf)
	if err != nil {
		return nil, err
	}
	client := &sshagentProxyClient{sshClient}
	return client, nil
}

func (client *sshagentProxyClient) DialUDP(network string, localAddr, remoteAddr *net.UDPAddr) (net.Conn, error) {
	return nil, ErrUnsupportedProtocol
}

func sshagentAuth(url *url.URL) []ssh.AuthMethod {
	methods := []ssh.AuthMethod{}
	publicKey := url.Query().Get("public-key")
	if publicKey != "" {
		buffer, err := ioutil.ReadFile(publicKey)
		if err != nil {
			return nil
		}
		key, err := ssh.ParsePrivateKey(buffer)
		if err != nil {
			return nil
		}
		method := ssh.PublicKeys(key)
		methods = append(methods, method)
	}
	if password, ok := url.User.Password(); ok {
		method := ssh.Password(password)
		methods = append(methods, method)
	}
	return methods
}
