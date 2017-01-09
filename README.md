# ProxyClient

the proxy c library

supported SOCK4, SOCKS4A, SOCKS5, HTTP, HTTPS, etc proxy protocols

# Documentation

The full documentation is available on [Godoc](//godoc.org/github.com/RouterScript/ProxyClient).

# Example
```go
package main

import (
	"fmt"
	"net/http"
	"github.com/RouterScript/ProxyClient"
)

func main() {
	proxy, _ := proxyclient.NewProxyClient("http://localhost:1080")
	client := &http.Client{
		Transport: &http.Transport{
			Dial: proxy.Dial,
		},
	}
	response, err := client.Head("http://www.google.com")
	if err != nil {
		panic(err)
	}
	fmt.Println(response.Body)
}
```

# Reference

see http://github.com/GameXG/ProxyClient
