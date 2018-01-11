package conf

import (
	"fmt"

	"github.com/roscopecoltran/configor"
	"gopkg.in/ini.v1"
)

type Config struct {
	// E3W - WEBUI
	Host string `default:"0.0.0.0" json:"host,omitempty" yaml:"host,omitempty" toml:"host,omitempty"`
	Port string `default:"8086" json:"port,omitempty" yaml:"port,omitempty" toml:"port,omitempty"`
	Auth bool   `default:"false" json:"auth,omitempty" yaml:"auth,omitempty" toml:"auth,omitempty"`
	// ETCD v3
	EtcdRootKey   string   `json:"etcd_root_key,omitempty" yaml:"etcd_root_key,omitempty" toml:"etcd_root_key,omitempty"`
	EtcdEndPoints []string `default:"http://etcd-1:2379,http://etcd1:2379" json:"etcd_endpoints,omitempty" yaml:"etcd_endpoints,omitempty" toml:"etcd_endpoints,omitempty"`
	EtcdUsername  string   `json:"etcd_username,omitempty" yaml:"etcd_username,omitempty" toml:"etcd_username,omitempty"`
	EtcdPassword  string   `json:"etcd_password,omitempty" yaml:"etcd_password,omitempty" toml:"etcd_password,omitempty"`
	DirValue      string   `json:"dir_value,omitempty" yaml:"dir_value,omitempty" toml:"dir_value,omitempty"`
	// ETCD v3 - TLS
	CertFile      string `json:"cert_file,omitempty" yaml:"cert_file,omitempty" toml:"cert_file,omitempty"`
	KeyFile       string `json:"key_file,omitempty" yaml:"key_file,omitempty" toml:"key_file,omitempty"`
	TrustedCAFile string `json:"trusted_ca_file,omitempty" yaml:"trusted_ca_file,omitempty" toml:"trusted_ca_file,omitempty"`

	// E3W - Server
	Server struct {
		Host string `default:"0.0.0.0" json:"host,omitempty" yaml:"host,omitempty" toml:"host,omitempty"`
		Port string `default:"8086" json:"port,omitempty" yaml:"port,omitempty" toml:"port,omitempty"`
		Auth bool   `default:"false" json:"auth,omitempty" yaml:"auth,omitempty" toml:"auth,omitempty"`
	} `json:"server,omitempty" yaml:"server,omitempty" toml:"server,omitempty"`

	// E3W - Front-End
	Front struct {
		Disabled bool   `default:"false" json:"disabled,omitempty" yaml:"disabled,omitempty" toml:"disabled,omitempty"`
		Debug    bool   `default:"false" json:"debug,omitempty" yaml:"debug,omitempty" toml:"debug,omitempty"`
		Host     string `default:"0.0.0.0" json:"host,omitempty" yaml:"host,omitempty" toml:"host,omitempty"`
		Port     int    `default:"3002" json:"port,omitempty" yaml:"port,omitempty" toml:"port,omitempty"`
		Dev      struct {
			Disabled bool   `default:"false" json:"disabled,omitempty" yaml:"disabled,omitempty" toml:"disabled,omitempty"`
			Debug    bool   `default:"false" json:"debug,omitempty" yaml:"debug,omitempty" toml:"debug,omitempty"`
			Dir      string `default:"/data/static/dist" json:"dir,omitempty" yaml:"dir,omitempty" toml:"dir,omitempty"`
		} `json:"dev,omitempty" yaml:"dev,omitempty" toml:"dev,omitempty"`
		Dist struct {
			Disabled bool   `default:"false" json:"disabled,omitempty" yaml:"disabled,omitempty" toml:"disabled,omitempty"`
			Debug    bool   `default:"false" json:"debug,omitempty" yaml:"debug,omitempty" toml:"debug,omitempty"`
			Dir      string `default:"/data/static/dist" json:"dir,omitempty" yaml:"dir,omitempty" toml:"dir,omitempty"`
		} `json:"dist,omitempty" yaml:"dist,omitempty" toml:"dist,omitempty"`
	} `json:"front,omitempty" yaml:"front,omitempty" toml:"front,omitempty"`

	// ETCD v2/v3 (optional: TLS)
	Etcd struct {
		Disabled  bool     `default:"false" json:"disabled,omitempty" yaml:"disabled,omitempty" toml:"disabled,omitempty"`
		Debug     bool     `default:"false" json:"debug,omitempty" yaml:"debug,omitempty" toml:"debug,omitempty"`
		Version   int      `json:"version,omitempty" yaml:"version,omitempty" toml:"version,omitempty"`
		RootKey   string   `json:"root_key,omitempty" yaml:"root_key,omitempty" toml:"root_key,omitempty"`
		EndPoints []string `default:"http://etcd-1:2379,http://etcd1:2379" json:"endpoints,omitempty" yaml:"endpoints,omitempty" toml:"endpoints,omitempty"`
		Username  string   `json:"username,omitempty" yaml:"username,omitempty" toml:"username,omitempty"`
		Password  string   `json:"password,omitempty" yaml:"password,omitempty" toml:"password,omitempty"`
		DirValue  string   `json:"dir_value,omitempty" yaml:"dir_value,omitempty" toml:"dir_value,omitempty"`
		TLS       struct {
			Disabled      bool   `default:"false" json:"disabled,omitempty" yaml:"disabled,omitempty" toml:"disabled,omitempty"`
			Debug         bool   `default:"false" json:"debug,omitempty" yaml:"debug,omitempty" toml:"debug,omitempty"`
			CertFile      string `json:"cert_file,omitempty" yaml:"cert_file,omitempty" toml:"cert_file,omitempty"`
			KeyFile       string `json:"key_file,omitempty" yaml:"key_file,omitempty" toml:"key_file,omitempty"`
			TrustedCAFile string `json:"trusted_ca_file,omitempty" yaml:"trusted_ca_file,omitempty" toml:"trusted_ca_file,omitempty"`
		} `json:"tls,omitempty" yaml:"tls,omitempty" toml:"tls,omitempty"`
	} `json:"etcd,omitempty" yaml:"etcd,omitempty" toml:"etcd,omitempty"`
}

func Init(filepath string) (*Config, error) {

	c := &Config{}

	useConfigor := false
	if useConfigor {
		configor.Load(&c, "config.yml")
		fmt.Printf("config: %#v", c)
	}

	cfg, err := ini.Load(filepath)
	if err != nil {
		return nil, err
	}

	appSec := cfg.Section("app")
	c.Port = appSec.Key("port").Value()
	c.Auth = appSec.Key("auth").MustBool()

	etcdSec := cfg.Section("etcd")
	c.EtcdRootKey = etcdSec.Key("root_key").Value()
	c.DirValue = etcdSec.Key("dir_value").Value()
	c.EtcdEndPoints = etcdSec.Key("addr").Strings(",")
	c.TrustedCAFile = etcdSec.Key("ca_file").Value()
	c.CertFile = etcdSec.Key("cert_file").Value()
	c.KeyFile = etcdSec.Key("key_file").Value()
	c.EtcdUsername = etcdSec.Key("username").Value()
	c.EtcdPassword = etcdSec.Key("password").Value()

	return c, nil
}
