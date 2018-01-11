package ajii

import (
	"sync"
)

type BaseConfig struct {
	sync.Mutex
	v2KeysUrl  string
	serviceUrl string
	foo        string
	Port       int
}

type SimpleNode struct {
	Key   string
	Value string
}

func (c *BaseConfig) GetFoo() string {
	return c.foo
}

func (c *BaseConfig) SetFoo(f string) {
	c.Lock()
	defer c.Unlock()
	c.foo = f
}

func (c *BaseConfig) Set(key string, value string) (string, error) {
	return "", nil
}

func (c *BaseConfig) V2KeysUrl() string {
	return c.v2KeysUrl
}

func (c *BaseConfig) ServiceUrl() string {
	return c.serviceUrl
}

type Config interface {
	Get(key string) (SimpleNode, error)
	Set(key string, value string) (string, error)
	Delete(key string) error
	Dump() ([]SimpleNode, error)

	GetFoo() string
	SetFoo(f string)
	ServiceUrl() string
	V2KeysUrl() string
}
