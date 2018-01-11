package model

type CrawlerConfig struct {
	MaxDepth int `default:"1" etcd:"max_depth" json:"max_depth,omitempty" yaml:"max_depth,omitempty" toml:"max_depth,omitempty"`
}
