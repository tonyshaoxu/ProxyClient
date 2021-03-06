package proxyclient

import (
	"errors"
	"net"
	"net/url"
	"time"
)

func newBlackholeProxyClient(_ *url.URL, _ Dial) (dial Dial, err error) {
	dial = func(network, address string) (net.Conn, error) {
		return blackholeConn{}, nil
	}
	return
}

type blackholeConn struct{}

func (conn blackholeConn) Read([]byte) (int, error)           { return 0, nil }
func (conn blackholeConn) Write(buffer []byte) (int, error)   { return len(buffer), nil }
func (conn blackholeConn) Close() error                       { return nil }
func (conn blackholeConn) LocalAddr() net.Addr                { return nil }
func (conn blackholeConn) RemoteAddr() net.Addr               { return nil }
func (conn blackholeConn) SetDeadline(t time.Time) error      { return errors.New("unsupported") }
func (conn blackholeConn) SetReadDeadline(t time.Time) error  { return errors.New("unsupported") }
func (conn blackholeConn) SetWriteDeadline(t time.Time) error { return errors.New("unsupported") }
