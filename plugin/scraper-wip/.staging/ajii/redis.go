package ajii

import (
	"gopkg.in/redis.v5"
)

type RedisConfig struct {
	BaseConfig
	Client *redis.Client
}

func NewRedisConfig() *RedisConfig {
	host := "localhost:6379"
	client := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: "",
		DB:       0,
	})
	return &RedisConfig{
		BaseConfig: BaseConfig{
			v2KeysUrl:  host,
			foo:        "bar",
			serviceUrl: "http://localhost:8020",
		},
		Client: client,
	}
}

func (c *RedisConfig) Set(key string, value string) (string, error) {
	err := c.Client.Set(key, value, 0).Err()
	return value, err
}

func (c *RedisConfig) Delete(key string) error {
	err := c.Client.Del(key).Err()

	return err
}

func (c *RedisConfig) Dump() ([]SimpleNode, error) {
	keys := c.Client.Keys("*").Val()
	var nodes []SimpleNode
	for _, key := range keys {
		value := c.Client.Get(key).Val()
		newNode := SimpleNode{
			Key:   key,
			Value: value,
		}
		nodes = append(nodes, newNode)
	}
	return nodes, nil
}
