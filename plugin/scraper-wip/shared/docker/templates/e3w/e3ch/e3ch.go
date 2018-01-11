package e3ch

import (
	// "fmt"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/pkg/transport"
	"github.com/roscopecoltran/e3w/conf"
	"github.com/soyking/e3ch"
)

func NewE3chClient(config *conf.Config) (*client.EtcdHRCHYClient, error) {

	etcdConfig := clientv3.Config{
		Endpoints: config.EtcdEndPoints,
		Username:  config.EtcdUsername,
		Password:  config.EtcdPassword,
	}

	if config.CertFile != "" || config.KeyFile != "" || config.TrustedCAFile != "" {
		tlsInfo := transport.TLSInfo{
			CertFile:      config.CertFile,
			KeyFile:       config.KeyFile,
			TrustedCAFile: config.TrustedCAFile,
		}
		tlsConfig, err := tlsInfo.ClientConfig()
		if err != nil {
			panic(err)
		}
		etcdConfig.TLS = tlsConfig
	}

	clt, err := clientv3.New(etcdConfig)
	if err != nil {
		return nil, err
	}

	client, err := client.New(clt, config.EtcdRootKey, config.DirValue)
	if err != nil {
		return nil, err
	}
	return client, client.FormatRootKey()
}

func CloneE3chClient(username, password string, client *client.EtcdHRCHYClient) (*client.EtcdHRCHYClient, error) {
	clt, err := clientv3.New(clientv3.Config{
		Endpoints: client.EtcdClient().Endpoints(),
		Username:  username,
		Password:  password,
	})
	if err != nil {
		return nil, err
	}
	return client.Clone(clt), nil
}
