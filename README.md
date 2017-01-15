# ProxyClient

the proxy client library

supported SOCK4, SOCKS4A, SOCKS5, HTTP, HTTPS etc proxy protocols

## Supported Protocols
- [x] HTTP
- [x] HTTPS
- [x] SOCKS4
- [x] SOCKS4A
- [x] SOCKS5
- [x] SOCKS5 with TLS
- [x] ShadowSocks
- [x] SSH Agent
- [ ] VMess

# Documentation

The full documentation is available on [Godoc](//godoc.org/github.com/RouterScript/ProxyClient).

# Example
```go
package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"github.com/RouterScript/ProxyClient"
)

func main() {
	dial, _ := proxyclient.NewProxyClient("http://localhost:8080")
	client := &http.Client{
		Transport: &http.Transport{
			DialContext: dial.WrappedContext(),
		},
	}
	request, err := client.Get("http://www.example.com")
	if err != nil {
		panic(err)
	}
	content, err := ioutil.ReadAll(request.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(content))
}
```

# Reference

see http://github.com/GameXG/ProxyClient
