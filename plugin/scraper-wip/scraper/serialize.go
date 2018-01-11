package scraper

import (
	"fmt"

	"github.com/vmihailenco/msgpack"
)

/*
	Refs:
	- https://godoc.org/github.com/vmihailenco/msgpack#pkg-examples
	-
*/

func testMsgPack() {
	type Item struct {
		Foo string
	}

	b, err := msgpack.Marshal(&Item{Foo: "bar"})
	if err != nil {
		panic(err)
	}

	var item Item
	err = msgpack.Unmarshal(b, &item)
	if err != nil {
		panic(err)
	}
	fmt.Println(item.Foo)
	// Output: bar
}
