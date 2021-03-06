package proxyclient

import (
	"errors"
	"io/ioutil"
	"net/url"

	"golang.org/x/crypto/ssh"
)

func NewSSHAgentProxyClient(proxy *url.URL, upstreamDial Dial) (dial Dial, err error) {
	if !limitSchemes(proxy, "ssh") {
		err = errors.New("scheme is not SSH")
		return
	}
	conf := &ssh.ClientConfig{
		User: proxy.User.Username(),
		Auth: sshagentAuth(proxy),
	}
	conn, err := upstreamDial("tcp", proxy.Host)
	if err != nil {
		return
	}
	sshConn, sshChans, sshRequests, err := ssh.NewClientConn(conn, proxy.Host, conf)
	if err != nil {
		return
	}
	sshClient := ssh.NewClient(sshConn, sshChans, sshRequests)
	dial = sshClient.Dial
	dial = dial.TCPOnly()
	return
}

func sshagentAuth(proxy *url.URL) []ssh.AuthMethod {
	methods := []ssh.AuthMethod{}
	publicKey := proxy.Query().Get("public-key")
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
	if password, ok := proxy.User.Password(); ok {
		method := ssh.Password(password)
		methods = append(methods, method)
	}
	return methods
}
