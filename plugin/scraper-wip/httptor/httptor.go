package httptor

import (
	"net/http"

	"golang.org/x/net/proxy"
)

func init() {
	dialer, err := proxy.SOCKS5("tcp", "localhost:9050", nil, proxy.Direct)
	if err != nil {
		panic(err)
	}
	http.DefaultClient.Transport = &http.Transport{Dial: dialer.Dial}
}
