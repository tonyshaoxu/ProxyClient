package proxyclient

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net"
	"net/url"

	ss "github.com/shadowsocks/shadowsocks-go/shadowsocks"
)

func NewShadowsocksProxyClient(proxy *url.URL, upstreamDial Dial) (dial Dial, err error) {
	if !limitSchemes(proxy, "ss") {
		return nil, errors.New("scheme is not SS")
	}
	if content, err := base64.StdEncoding.DecodeString(proxy.Host); err == nil {
		proxy, err = proxy.Parse(fmt.Sprintf("ss://%s", string(content)))
		if err != nil {
			return nil, err
		}
	}
	var cipher *ss.Cipher
	if password, ok := proxy.User.Password(); ok {
		username := proxy.User.Username()
		cipher, err = ss.NewCipher(username, password)
		if err != nil {
			return
		}
	}
	dial = func(network, address string) (ssConn net.Conn, err error) {
		conn, err := upstreamDial("tcp", proxy.Host)
		if err != nil {
			return
		}
		rawAddr, err := ss.RawAddr(address)
		if err != nil {
			return
		}
		ssConn = ss.NewConn(conn, cipher.Copy())
		if _, err = ssConn.Write(rawAddr); err != nil {
			ssConn.Close()
			return
		}
		return
	}
	dial = dial.TCPOnly()
	return
}
