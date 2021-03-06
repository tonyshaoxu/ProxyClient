package proxyclient

import (
	"bytes"
	"net"
	"net/url"
)

func NewSplitPacketProxyClient(_ *url.URL, upstreamDial Dial) (dial Dial, err error) {
	dial = func(network, address string) (conn net.Conn, err error) {
		conn, err = upstreamDial(network, address)
		if err != nil {
			return
		}
		return splitPacketConn{conn}, err
	}
	dial = dial.TCPOnly()
	return
}

type splitPacketConn struct{ net.Conn }

func (conn splitPacketConn) Write(packet []byte) (n int, err error) {
	for _, spittedPacket := range splitHTTPPacket(packet) {
		writeLen, err := conn.Conn.Write(spittedPacket)
		n += writeLen
		if err != nil {
			return n, err
		}
	}
	return
}

func splitHTTPPacket(buffer []byte) (response [][]byte) {
	splitPacket := func(buffer []byte, pos int) [][]byte {
		if len(buffer) > pos {
			return [][]byte{buffer[:pos], buffer[pos:]}
		}
		return [][]byte{buffer}
	}
	compose := func(cursor, pos int, prefix []byte) [][]byte {
		if !bytes.HasPrefix(buffer[cursor + 1:], prefix) {
			return nil
		}
		packets := splitPacket(buffer, cursor+1)
		packets = append([][]byte{packets[0]}, splitPacket(packets[1], pos)...)
		return append(packets[:len(packets) - 1], splitHTTPPacket(packets[len(packets) - 1])...)
	}
	response = [][]byte{buffer}
	for cursor, ch := range buffer {
		switch ch {
		case 'G':
			response = compose(cursor, 3, []byte("ET "))
		case 'P':
			response = compose(cursor, 5, []byte("OST "))
		case 'C':
			response = compose(cursor, 8, []byte("ONNECT "))
		case 'H':
			response = compose(cursor, 8, []byte("OST "))
			if response != nil {
				return
			}
			response = compose(cursor, 9, []byte("TTP "))
		}
	}
	return response
}
