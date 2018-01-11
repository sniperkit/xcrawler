package scraper

import (
	"bytes"
	"crypto/rsa"
	"crypto/sha1"
	"encoding/asn1"
	"encoding/base32"
	// "log"
	"net/http"
	// "os"
	// "sort"

	// "github.com/ak1t0/flame/crawler"
	// "github.com/ak1t0/flame/format"
	// "github.com/ak1t0/flame/reader"
	"golang.org/x/net/proxy"
)

/*
	Refs:
	- https://github.com/ak1t0/flame/blob/master/main.go
	- https://github.com/ak1t0/flame/blob/master/test.json
	- https://github.com/hanjm/fasthttpmiddleware
	- https://github.com/mischief/goniongen/blob/master/main.go
	- https://github.com/subgraph/onioncfg/blob/master/onioncfg.go
	- https://github.com/DonnchaC/oniongateway/blob/master/entry_proxy/host2onion.yaml
	- https://github.com/Abdullah2993/toker/blob/master/service/alpine/1/main.go
	- https://github.com/Abdullah2993/toker/blob/master/proxy/alpine/Dockerfile
	- https://github.com/prsolucoes/go-tor-crawler
	- https://github.com/nogoegst/onionutil
	- https://github.com/codekoala/torotator
	- https://github.com/codekoala/torotator
	- https://github.com/miolini/metasocks
	- https://github.com/gearmover/web2tor
	- https://github.com/phoenix1342/TOR-Browser/blob/master/relay/relay.go
*/

var onionencoding = base32.NewEncoding("abcdefghijklmnopqrstuvwxyz234567")

func init() {
	dialer, err := proxy.SOCKS5("tcp", "localhost:9050", nil, proxy.Direct)
	if err != nil {
		panic(err)
	}
	http.DefaultClient.Transport = &http.Transport{Dial: dialer.Dial}
}

// creates an onion address given an rsa public key component
func address(pub *rsa.PublicKey) string {
	derbytes, _ := asn1.Marshal(*pub)

	// 1. Let H = H(PK).
	hash := sha1.New()
	hash.Write(derbytes)
	sum := hash.Sum(nil)

	// 2. Let H' = the first 80 bits of H, considering each octet from
	//    most significant bit to least significant bit.
	sum = sum[:10]

	// 3. Generate a 16-character encoding of H', using base32 as defined
	//    in RFC 4648.
	var buf32 bytes.Buffer
	b32enc := base32.NewEncoder(onionencoding, &buf32)
	b32enc.Write(sum)
	b32enc.Close()

	return buf32.String()
}
