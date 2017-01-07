package proxyclient

import "errors"

var (
	ErrUnsupportedNetwork   = errors.New("Unsupported network.")
	ErrUnsupportedProtocol  = errors.New("Unsupported protocol.")
	ErrUnsupportedLocalAddr = errors.New("Unsupported local address.")
)
