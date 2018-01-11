package operation

import (
	"github.com/coreos/etcd/client"
)

// EtcdClient etcd client
type EtcdClient struct {
	Endpoint string
	KeysAPI  client.KeysAPI
}

// Element object of a etcd value
type Element struct {
	Key         string
	Dir         bool
	Value       string
	ModifyIndex uint64
	CreateIndex uint64
}
