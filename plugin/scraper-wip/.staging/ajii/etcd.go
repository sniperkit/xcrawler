package ajii

import (
	"encoding/json"
	// "fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

type EtcdConfig struct {
	sync.Mutex
	V2KeysUrl  string
	ServiceUrl string
	foo        string
	Port       int
	// BaseConfig
	Client *http.Client
}

type EtcdResponse struct {
	Action string `json:"action"`
	Node   struct {
		Key   string `json:"key"`
		Dir   bool   `json:"dir"`
		Nodes []Node `json:"nodes"`
	} `json:"node"`
}

type GetResponse struct {
	Action string `json:"action"`
	Node   Node   `json:"node"`
}

type Node struct {
	Key           string `json:"key"`
	Value         string `json:"value"`
	ModifiedIndex int    `json:"modifiedIndex"`
	CreatedIndex  int    `json:"createdIndex"`
}

func NewConfig() *EtcdConfig {
	client := &http.Client{}
	return &EtcdConfig{
		V2KeysUrl:  "http://localhost:4001/v2/keys/",
		foo:        "bar",
		ServiceUrl: "http://localhost",
		Client:     client,
	}
}

func (c *EtcdConfig) Set(key string, value string) (string, error) {
	u := c.V2KeysUrl + key

	form := url.Values{}
	form.Add("value", value)
	body := strings.NewReader(form.Encode())

	req, _ := http.NewRequest("PUT", u, body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	_, err := c.Client.Do(req)
	return value, err
}

func (c *EtcdConfig) Delete(key string) error {
	u := c.V2KeysUrl + key
	req, _ := http.NewRequest("DELETE", u, nil)

	_, err := c.Client.Do(req)
	return err
}

func (c *EtcdConfig) Dump() ([]SimpleNode, error) {
	u := c.V2KeysUrl + "?recursive=true"
	resp, _ := http.Get(u)

	var r EtcdResponse
	json.NewDecoder(resp.Body).Decode(&r)
	var nodes []SimpleNode
	for _, node := range r.Node.Nodes {
		newNode := SimpleNode{
			Key:   node.Key,
			Value: node.Value,
		}
		nodes = append(nodes, newNode)
	}
	return nodes, nil
}

func (c *EtcdConfig) Get(key string) (SimpleNode, error) {
	u := c.V2KeysUrl + key
	resp, _ := http.Get(u)

	var r GetResponse
	json.NewDecoder(resp.Body).Decode(&r)
	newNode := SimpleNode{
		Key:   r.Node.Key,
		Value: r.Node.Value,
	}
	return newNode, nil
}
