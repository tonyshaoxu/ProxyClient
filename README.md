# ProxyClient

the proxy client library

supported SOCK4, SOCKS4A, SOCKS5, HTTP, HTTPS, etc proxy protocols

# Documentation

The full documentation is available on [Godoc](//godoc.org/github.com/RouterScript/ProxyClient).

# Example
```go
package main

import (
	"fmt"
	"github.com/RouterScript/ProxyClient"
)

func main() {
	client, err := proxyclient.NewProxyClient("http://localhost:1080")
	if err != nil {
		panic(err)
	}
	conn, err := client.Dial("tcp", "www.google.com:80")
	if err != nil {
		panic(err)
	}
	if _, err := conn.Write([]byte("GET / HTTP/1.0\r\nHOST:www.google.com\r\n\r\n")); err != nil {
		panic(err)
	}
	buffer := make([]byte, 2048)
	if n, err := conn.Read(buffer); err != nil {
		panic(err)
	} else {
		fmt.Print(string(buffer[:n]))
	}
	conn.Close()
}
```

# Reference

see http://github.com/GameXG/ProxyClient
