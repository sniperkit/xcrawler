package operation

import (
	"context"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/robinle/etcdweb/api"

	"github.com/coreos/etcd/client"
)

// InitClient init a etcd client
func (c *EtcdClient) InitClient(config *api.Config) {
	cfg := client.Config{
		Endpoints: []string{config.Endpoint},
		Transport: client.DefaultTransport,
	}
	if strings.Split(config.Endpoint, ":")[0] == "https" {
		httpsTransport := &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			Dial: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).Dial,
			TLSHandshakeTimeout: 10 * time.Second,
		}
		cfg.Transport = httpsTransport
	}

	cl, err := client.New(cfg)
	if err != nil {
		panic("etcd client generate fail")
	}
	c.Endpoint = config.Endpoint
	c.KeysAPI = client.NewKeysAPI(cl)
}

// Set save data in etcd
func (c *EtcdClient) Set(key string, value string, dir bool) error {
	_, err := c.KeysAPI.Set(
		context.Background(),
		key,
		value,
		&client.SetOptions{Dir: dir},
	)
	if err != nil {
		return err
	}
	return nil
}

// GetDirKeys get all keys of dir
func (c *EtcdClient) GetDirKeys(dir string) ([]Element, error) {
	elements := []Element{}
	resp, err := c.KeysAPI.Get(context.Background(), dir, nil)
	if err != nil {
		return elements, err
	}
	for _, node := range resp.Node.Nodes {
		element := Element{}
		element.Key = node.Key
		element.Value = node.Value
		element.Dir = node.Dir
		element.CreateIndex = node.CreatedIndex
		element.ModifyIndex = node.CreatedIndex
		elements = append(elements, element)
	}
	return elements, nil
}

// GetValue get value of key
func (c *EtcdClient) GetValue(key string) (string, error) {
	value := ""
	resp, err := c.KeysAPI.Get(context.Background(), key, nil)
	if err != nil {
		return value, err
	}
	value = resp.Node.Value
	return value, nil
}

// GetKeyValue get key value from etcd
func (c *EtcdClient) GetKeyValue(key string) (Element, error) {
	element := Element{}
	resp, err := c.KeysAPI.Get(context.Background(), key, nil)
	if err != nil {
		return element, err
	}
	// baseKey := filepath.Base(key)
	element.Key = resp.Node.Key
	element.Value = resp.Node.Value
	element.Dir = resp.Node.Dir
	element.CreateIndex = resp.Node.CreatedIndex
	element.ModifyIndex = resp.Node.CreatedIndex
	return element, nil
}

// Update update specfic key's value
func (c *EtcdClient) Update(key string, value string) error {
	_, err := c.KeysAPI.Update(context.Background(), key, value)
	if err != nil {
		return err
	}
	return nil
}

// Delete delete a key from etcd
func (c *EtcdClient) Delete(key string) error {
	_, err := c.KeysAPI.Delete(
		context.Background(),
		key,
		&client.DeleteOptions{Dir: true, Recursive: true},
	)
	if err != nil {
		_, err := c.KeysAPI.Delete(
			context.Background(),
			key,
			&client.DeleteOptions{},
		)
		if err != nil {
			return err
		}
	}
	return nil
}

// CheckKey check if key exist in etcd
func (c *EtcdClient) CheckKey(key string) (bool, error) {
	_, err := c.GetKeyValue(key)
	if err != nil {
		return false, err
	}
	return true, nil
}

//func SetDir() {}
