package main

import (
	// "github.com/coreos/go-etcd/etcd/client"
	"github.com/jinuljt/getcds"
)

func main() {
	// 定义struct
	var S struct {
		I32 int    `etcd:"i32"`
		I64 int    `etcd:"i64"`
		Str string `etcd:"str"`

		S struct {
			I32 int    `etcd:"i32"`
			I64 int    `etcd:"i64"`
			Str string `etcd:"str"`
		} `etcd:"s"`
	}

	machines := []string{"http://192.168.1.58:2379"}
	client := getcds.NewClient(machines)
	defer client.Close()

	if err := client.Unmarshal("/test", &S); err != nil {
		fmt.Println("getcd unmarshal error due to", err)
	}
}
