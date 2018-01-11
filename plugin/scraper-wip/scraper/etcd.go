package scraper

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	etcd "github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/etcdserver/api/v3rpc/rpctypes"
	"github.com/coreos/etcd/pkg/transport"
	"github.com/iancoleman/strcase"
	"github.com/k0kubun/pp"
	"github.com/oleiade/reflections"
	"github.com/roscopecoltran/e3ch"
	"golang.org/x/net/context"
	// "github.com/fatih/structs"
	// "github.com/iancoleman/orderedmap"
	// "github.com/coreos/etcd/clientv3/mirror"
	// etcderr "github.com/coreos/etcd/error"
	// "github.com/coreos/etcd/mvcc/mvccpb"
	// "github.com/jinuljt/getcds"
	// "github.com/damoye/etcd-config"
	// "github.com/roscopecoltran/e3w/routers"
)

/*
	Refs:
	- https://github.com/vmattos/apps-registrator/blob/master/etcd/etcd.go
	- https://github.com/Financial-Times/vulcan-config-builder/blob/master/main.go
	- https://github.com/xiang90/edb/blob/master/sql.go
	- https://github.com/rafaeljusto/etcetera/blob/master/etcetera.go
	- https://github.com/mickep76/etcdmap
	- https://github.com/Financial-Times/vulcan-config-builder/blob/master/main.go
	- https://github.com/ruprict/vulcand-atd-transformer/blob/master/templates.go
	- https://github.com/vmattos/apps-registrator/blob/master/etcd/etcd.go
	- https://github.com/vulcand/vulcand/blob/master/engine/etcdv3ng/etcd.go
	- https://github.com/vulcand/vulcand/blob/master/engine/etcdv2ng/etcd.go
	- https://github.com/xh4n3/ohetcd/blob/master/main.go
	- https://github.com/BlueDragonX/sentinel
	- https://github.com/mickep76/etcdtool/tree/master
	- https://github.com/jinuljt/getcds/blob/master/etcd.go
	- https://github.com/coreos/etcd/blob/master/clientv3/example_watch_test.go
	- https://github.com/coreos/etcd/blob/master/clientv3/example_test.go
	- https://github.com/kelseyhightower/confd/blob/master/backends/etcdv3/client.go
	- https://github.com/kelseyhightower/confd/blob/master/backends/client.go
*/

const (
	ETCD_CLIENT_TIMEOUT = 3 * time.Second
)

var (
	numDirs, numKeys        int
	routeIdRegex            = regexp.MustCompile("/routes/([^/]+)(?:/route)?$")
	headerRegex             = regexp.MustCompile("/routes/([^/]+)/headers/([^/]+)$")
	blockRegex              = regexp.MustCompile("/routes/([^/]+)/blocks/([^/]+)$")
	addressRegex            = regexp.MustCompile(`^[\.\-:\/\w]*:[0-9]{2,5}$`)
	serverRegex             = regexp.MustCompile("/backends/([^/]+)/servers/([^/]+)$")
	frontendIdRegex         = regexp.MustCompile("/frontends/([^/]+)(?:/frontend)?$")
	backendIdRegex          = regexp.MustCompile("/backends/([^/]+)(?:/backend)?$")
	hostnameRegex           = regexp.MustCompile("/hosts/([^/]+)(?:/host)?$")
	listenerIdRegex         = regexp.MustCompile("/listeners/([^/]+)")
	middlewareRegex         = regexp.MustCompile("/frontends/([^/]+)/middlewares/([^/]+)$")
	endpointRegex           = regexp.MustCompile("([A-Za-z0-9/]+)/endpoint/([A-Za-z0-9/]+)/config")
	configRegex             = regexp.MustCompile("([A-Za-z0-9/]+)/endpoint/([A-Za-z0-9/]+)/config/([^/]+)")
	scraperRegex            = regexp.MustCompile("([A-Za-z0-9/]+)/endpoint/([A-Za-z0-9/]+)/config/scraper")
	scraperParamRegex       = regexp.MustCompile("([A-Za-z0-9/]+)/endpoint/([A-Za-z0-9/]+)/config/scraper/([^/]+)")
	scraperGroupRegex       = regexp.MustCompile("([A-Za-z0-9/]+)/endpoint/([A-Za-z0-9/]+)/config/scraper/([^/]+)/([^/]+)")
	scraperExtractRegex     = regexp.MustCompile("([A-Za-z0-9/]+)/endpoint/([A-Za-z0-9/]+)/config/scraper/extract/([^/]+)")
	scraperBlocksRegex      = regexp.MustCompile("([A-Za-z0-9/]+)/endpoint/([A-Za-z0-9/]+)/config/scraper/blocks/([^/]+)/([^/]+)")
	scraperBlockDetailRegex = regexp.MustCompile("([A-Za-z0-9/]+)/endpoint/([A-Za-z0-9/]+)/config/scraper/blocks/([^/]+)/([^/]+)")
)

type EtcdConfig struct {
	Disabled            bool                    `default:"false" help:"Disable etcd client" json:"disabled,omitempty" yaml:"disabled,omitempty" toml:"disabled,omitempty"`
	Client              *etcd.Client            `gorm:"-" json:"-" yaml:"-" toml:"-"`
	Once                *sync.Once              `gorm:"-" json:"-" yaml:"-" toml:"-"`
	E3ch                *client.EtcdHRCHYClient `gorm:"-" json:"-" yaml:"-" toml:"-"`
	Context             context.Context         `gorm:"-" json:"-" yaml:"-" toml:"-"`
	CancelFunc          context.CancelFunc      `gorm:"-" json:"-" yaml:"-" toml:"-"`
	kv                  map[string]string       `gorm:"-" json:"-" yaml:"-" toml:"-"`
	mutex               sync.RWMutex            `gorm:"-" json:"-" yaml:"-" toml:"-"`
	rch                 etcd.WatchChan          `gorm:"-" json:"-" yaml:"-" toml:"-"`
	prefix              string                  `gorm:"-" json:"-" yaml:"-" toml:"-"`
	Handler             *Handler                `gorm:"-" json:"-" yaml:"-" toml:"-"`
	info                map[string]info         `gorm:"-" json:"-" yaml:"-" toml:"-"` // info creates a correlation between a path to a info structure that stores some extra information and make the API usage easier
	SyncIntervalSeconds int64                   `json:"sync_interval_seconds,omitempty" yaml:"sync_interval_seconds,omitempty" toml:"sync_interval_seconds,omitempty"`
	Consistency         string                  `json:"consistency,omitempty" yaml:"consistency,omitempty" toml:"consistency,omitempty"`
	RequireQuorum       bool                    `json:"require_quorum,omitempty" yaml:"require_quorum,omitempty" toml:"require_quorum,omitempty"`
	OnceError           error                   `gorm:"-" json:"-" yaml:"-" toml:"-"`
	InitCheck           bool                    `json:"init_check,omitempty" yaml:"init_check,omitempty" toml:"init_check,omitempty"`
	MaxDir              int                     `default:"10" json:"max_dir,omitempty" yaml:"max_dir,omitempty" toml:"max_dir,omitempty"`
	ApiVersion          int                     `default:"3" json:"api_version,omitempty" yaml:"api_version,omitempty" toml:"api_version,omitempty"`
	Peers               []string                `gorm:"-" json:"peers,omitempty" yaml:"peers,omitempty" toml:"peers,omitempty"`
	MaxTimeout          time.Duration           `json:"timeout,omitempty" yaml:"timeout,omitempty" toml:"timeout,omitempty"`
	DialTimeout         time.Duration           `json:"dial_timeout,omitempty" yaml:"dial_timeout,omitempty" toml:"dial_timeout,omitempty"`
	ReadTimeout         time.Duration           `json:"read_timeout,omitempty" yaml:"read_timeout,omitempty" toml:"read_timeout,omitempty"`
	WriteTimeout        time.Duration           `json:"write_timeout,omitempty" yaml:"write_timeout,omitempty" toml:"write_timeout,omitempty"`
	CommandTimeout      time.Duration           `json:"command_timeout,omitempty" yaml:"command_timeout,omitempty" toml:"command_timeout,omitempty"`
	Routes              []EtcdRouteConfig       `json:"routes" yaml:"routes" toml:"routes"`
	Username            string                  `json:"username,omitempty" yaml:"username,omitempty" toml:"username,omitempty"`
	Password            string                  `json:"password,omitempty" yaml:"password,omitempty" toml:"password,omitempty"`
	PasswordFilePath    string                  `json:"password_file,omitempty" yaml:"password_file,omitempty" toml:"password_file,omitempty"`
	IsSecured           bool                    `json:"tls,omitempty" yaml:"tls,omitempty" toml:"tls,omitempty"`
	CertFile            string                  `json:"cert_file,omitempty" yaml:"cert_file,omitempty" toml:"cert_file,omitempty"`
	KeyFile             string                  `json:"key_file,omitempty" yaml:"key_file,omitempty" toml:"key_file",omitempty`
	TrustedCAFile       string                  `json:"trusted_ca_file,omitempty" yaml:"trusted_ca_file,omitempty" toml:"trusted_ca_file,omitempty"`
	RootKey             string                  `default:"root" json:"root_key,omitempty" yaml:"root_key,omitempty" toml:"root_ke,omitemptyy"`
	DirValue            string                  `json:"dir_value,omitempty" yaml:"dir_value,omitempty" toml:"dir_value,omitempty"`
	SealKey             string                  `json:"seal_key,omitempty" yaml:"seal_key,omitempty" toml:"seal_key,omitempty"`
	TrustForwardHeader  bool                    `json:"trust_forward_header,omitempty" yaml:"trust_forward_header,omitempty" toml:"trust_forward_header,omitempty"`
	MemProfileRate      int                     `json:"mem_profile_rate,omitempty" yaml:"mem_profile_rate,omitempty" toml:"mem_profile_rate,omitempty"`
	Debug               bool                    `help:"Enable debug output" json:"debug,omitempty" yaml:"debug,omitempty" toml:"debug,omitempty"`
}

func (ectl *EtcdConfig) NewEtcdClient(conf etcd.Config) (*etcd.Client, error) {
	var err error
	ectl.Client, err = etcd.New(etcd.Config{
		Endpoints:   []string{"etcd1:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		fmt.Println("failed to init initial etcd v3 client, error: ", err)
		return nil, err
	} else {
		fmt.Println("[SUCCESS] etcd v3 client connected !")
	}
	return ectl.Client, nil
}

func (ectl *EtcdConfig) NewE3chClient() (*client.EtcdHRCHYClient, error) {
	if ectl.MaxDir == 0 {
		ectl.MaxDir = 10
	}
	if len(ectl.Peers) == 0 {
		ectl.Peers = []string{"http://localhost:2379"}
	}
	if ectl.RootKey == "" {
		ectl.RootKey = "sniperkit"
	}
	etcdConfig := etcd.Config{
		Endpoints: ectl.Peers,
		Username:  ectl.Username,
		Password:  ectl.Password,
	}
	if ectl.IsSecured {
		tlsInfo := transport.TLSInfo{
			CertFile:      ectl.CertFile,
			KeyFile:       ectl.KeyFile,
			TrustedCAFile: ectl.TrustedCAFile,
		}
		tlsConfig, err := tlsInfo.ClientConfig()
		if err != nil {
			return nil, err
		}
		etcdConfig.TLS = tlsConfig
	}
	if ectl.Client == nil {
		var cerr error
		ectl.Client, cerr = ectl.NewEtcdClient(etcdConfig)
		if cerr != nil {
			return nil, cerr
		}
	}
	client, err := client.New(ectl.Client, ectl.RootKey, ectl.DirValue)
	if err != nil {
		return nil, err
	}
	ectl.E3ch = client
	ectl.kv = map[string]string{}
	ectl.prefix = ectl.RootKey // ectl.RootKey + "/"

	// init and watch
	err = ectl.initAndWatch()
	if err != nil {
		if ectl.Debug {
			fmt.Println("initAndWatch, error: ", err)
		}
		ectl.Client.Close()
		return nil, err
	}

	// loop to update
	go func() {
		for {
			for wresp := range ectl.rch {
				for _, ev := range wresp.Events {
					endpointKey := string(ev.Kv.Key)
					var isScraperBlock, isScraperBlockDetails, isScraperExtract bool
					var endpointRoute, endpointField, endpointConfigType, scraperConfigParameter, scraperConfigGroup, scraperConfigBlock string
					if endpointIds := endpointRegex.FindStringSubmatch(endpointKey); len(endpointIds) == 2 {
						if ectl.Debug {
							fmt.Println("[COUNT] endpointIds: ", len(endpointIds))
							pp.Println("[PARTS] endpointIds: ", endpointIds)
						}
						endpointRoute = endpointIds[2]
					}
					if configIds := configRegex.FindStringSubmatch(endpointKey); len(configIds) == 3 {
						if ectl.Debug {
							fmt.Println("[COUNT] configIds: ", len(configIds))
							pp.Println("[PARTS] configIds: ", configIds)
						}
						endpointRoute = configIds[2]
						endpointConfigType = strcase.ToCamel(configIds[3])
					}
					if scraperParamIds := scraperParamRegex.FindStringSubmatch(endpointKey); len(scraperParamIds) >= 3 {
						if ectl.Debug {
							fmt.Println("[COUNT] scraperParamIds: ", len(scraperParamIds))
							pp.Println("[PARTS] scraperParamIds: ", scraperParamIds)
						}
						endpointRoute = scraperParamIds[2]
						scraperConfigParameter = scraperParamIds[3]
					}
					if scraperExtractIds := scraperExtractRegex.FindStringSubmatch(endpointKey); len(scraperExtractIds) >= 3 {
						if ectl.Debug {
							fmt.Println("[COUNT] scraperExtractIds: ", len(scraperExtractIds))
							pp.Println("[PARTS] scraperExtractIds: ", scraperExtractIds)
						}
						endpointRoute = scraperExtractIds[2]
						endpointField = strcase.ToCamel(scraperExtractIds[3])
						isScraperExtract = true
					}

					if configGroupParamIds := scraperGroupRegex.FindStringSubmatch(endpointKey); len(configGroupParamIds) >= 4 {
						if ectl.Debug {
							fmt.Println("[COUNT] configGroupParamIds: ", len(configGroupParamIds))
							pp.Println("[PARTS] configGroupParamIds: ", configGroupParamIds)
						}
						endpointRoute = configGroupParamIds[2]
						scraperConfigGroup = configGroupParamIds[3]
						scraperConfigParameter = configGroupParamIds[4]
						isScraperBlock = true
					}
					if configBlockParamIds := scraperBlocksRegex.FindStringSubmatch(string(ev.Kv.Key)); len(configBlockParamIds) >= 5 {
						if ectl.Debug {
							fmt.Println("[COUNT] configBlockParamIds: ", len(configBlockParamIds))
							pp.Println("[PARTS] configBlockParamIds: ", configBlockParamIds)
						}
						endpointRoute = configBlockParamIds[2]
						scraperConfigGroup = configBlockParamIds[3]
						scraperConfigBlock = configBlockParamIds[4]
						isScraperBlockDetails = true
					}
					if scraperConfigParameter != "" {
						endpointField = strcase.ToCamel(scraperConfigParameter)
						if ectl.Handler != nil {
							endpoint := ectl.Handler.Endpoint(endpointRoute)
							if ectl.Debug {
								fmt.Println("[VARS] endpointRoute: ", endpointRoute)
								fmt.Println("[VARS] endpointConfigType: ", endpointConfigType)
								fmt.Println("[VARS] endpointField: ", endpointField)
								fmt.Println("[VARS] endpoint.ready: ", endpoint.ready)
								fmt.Println("[VARS] scraperConfigParameter: ", scraperConfigParameter)
								fmt.Println("[VARS] scraperConfigGroup: ", scraperConfigGroup)
								fmt.Println("[VARS] scraperConfigBlock: ", scraperConfigBlock)
								fmt.Printf("[FLAG] isScraperBlock: %t\n", isScraperBlock)
								fmt.Printf("[FLAG] isScraperBlockDetails: %t\n", isScraperBlockDetails)
							}
							if endpoint.ready {
								if isScraperExtract {
									status, err := strconv.ParseBool(string(ev.Kv.Value))
									if err != nil {
										if ectl.Debug {
											fmt.Println("error: ", err)
										}
									}
									if err := reflections.SetField(&endpoint.Extract, endpointField, status); err != nil {
										if ectl.Debug {
											fmt.Println("SetField() error: ", err)
										}
									}
								} else if isScraperBlockDetails { // nested struct (extractBlocks,...)
									if ectl.Debug {
										fmt.Printf("UpdateConfig for route group: endpoint=%s > block=%s > key=%s, value=%s \n", endpointRoute, scraperConfigGroup, scraperConfigBlock, ev.Kv.Value)
									}
									/*
										if scraperConfigBlock == "items" {
											fmt.Println("---> scraperConfigGroup: ", scraperConfigGroup)
											fmt.Println("---> scraperConfigBlock: ", strcase.ToCamel(scraperConfigBlock))
											fmt.Println("---> endpoint.BlocksJSON[scraperConfigGroup].Items: ", endpoint.BlocksJSON[scraperConfigGroup].Items)
											isEndpointField, err := reflections.HasField(endpoint.BlocksJSON[scraperConfigGroup], "Items")
											if err != nil {
												if ectl.Debug {
													fmt.Println("error: ", err)
												}
											}
											if ectl.Debug {
												pp.Printf("[VARS] isEndpointField: %s\n", isEndpointField)
											}
											endpointFieldKind, err := reflections.GetFieldKind(endpoint.BlocksJSON[scraperConfigGroup], "Items")
											if err != nil {
												if ectl.Debug {
													fmt.Println("error: ", err)
												}
											}
											if ectl.Debug {
												pp.Printf("[VARS] endpointFieldKind: %T\n", endpointFieldKind)
											}
											// group := endpoint.BlocksJSON[scraperConfigGroup]
											if err := reflections.SetField(group, "Items", string(ev.Kv.Value)); err != nil {
												if ectl.Debug {
													fmt.Println("SetField() error: ", err)
												}
											}
										} else {
									*/
									if len(endpoint.BlocksJSON[scraperConfigGroup].Details[scraperConfigBlock]) > 0 {
										endpoint.BlocksJSON[scraperConfigGroup].Details[scraperConfigBlock][0].val = string(ev.Kv.Value)
									}
									//}
								} else {
									isEndpointField, err := reflections.HasField(endpoint, endpointField)
									if err != nil {
										if ectl.Debug {
											fmt.Println("error: ", err)
										}
									}
									if ectl.Debug {
										pp.Println("[VARS] isEndpointField: ", isEndpointField)
									}
									if isEndpointField {
										endpointFieldKind, err := reflections.GetFieldKind(endpoint, endpointField)
										if err != nil {
											if ectl.Debug {
												fmt.Println("error: ", err)
											}
										}
										switch endpointFieldKind {
										case reflect.String:
											if err := reflections.SetField(endpoint, endpointField, string(ev.Kv.Value)); err != nil {
												if ectl.Debug {
													fmt.Println("SetField() error: ", err)
												}
											}
										case reflect.Int:
											i, err := strconv.Atoi(string(ev.Kv.Value))
											if ectl.Debug {
												fmt.Println("SetField() error: ", err)
											}
											if err := reflections.SetField(endpoint, endpointField, i); err != nil {
												if ectl.Debug {
													fmt.Println("SetField() error: ", err)
												}
											}
										case reflect.Bool:
											b, err := strconv.ParseBool(string(ev.Kv.Value))
											if ectl.Debug {
												fmt.Println("SetField() error: ", err)
											}
											if err := reflections.SetField(endpoint, endpointField, b); err != nil {
												if ectl.Debug {
													fmt.Println("SetField() error: ", err)
												}
											}
										default:
											fmt.Println("Unkown type, SetField() ignored...")
										}
									}
								}
								if ectl.Debug {
									pp.Println(ectl.Handler.Endpoint(endpointRoute))
								}
							}
							if ectl.Debug {
								fmt.Printf("[UPDATE] set: field=%s, key=%s, value=%s \n", endpointField, ev.Kv.Key, ev.Kv.Value)
							}
						}
					}
				}

			}
			if ectl.Debug {
				fmt.Println("etcd-config watch channel closed")
			}
			for {
				err = ectl.initAndWatch()
				if err == nil {
					break
				}
				if ectl.Debug {
					fmt.Println("etcd-config get failed: ", err)
				}
				time.Sleep(time.Second)
			}
		}
	}()

	if ectl.InitCheck {
		report, err := ectl.CheckupE3ch()
		if ectl.Debug {
			fmt.Printf("[WARNING] failed to pass all the E3CH client init tests:\n- fatal_error:\n%#v\n- warnings:\n%s\n", err, strings.Join(report, "\n"))
		}
	}

	return client, client.FormatRootKey()
}

func (ectl *EtcdConfig) Search(input interface{}, key string) string {
	val := reflect.ValueOf(input).Elem()
	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		if key == strings.ToLower(typeField.Name) {
			return valueField.Interface().(string)
		}
	}
	return ""
}

// Get ...
func (ectl *EtcdConfig) Get(key string) string {
	ectl.mutex.RLock()
	defer ectl.mutex.RUnlock()
	return ectl.kv[key]
}

func (ectl *EtcdConfig) String() string {
	ectl.mutex.RLock()
	defer ectl.mutex.RUnlock()
	b, _ := json.MarshalIndent(ectl.kv, "", "  ")
	return string(b)
}

func (ectl *EtcdConfig) initAndWatch() error {
	ectl.rch = ectl.Client.Watch(context.TODO(), ectl.RootKey, etcd.WithPrefix())
	resp, err := ectl.Client.Get(context.TODO(), ectl.RootKey, etcd.WithPrefix())
	if err != nil {
		return err
	}
	for _, kv := range resp.Kvs {
		ectl.set(kv.Key, kv.Value)
	}
	return nil
}

/*
func (ectl *EtcdConfig) set(key, value []byte) {
	strKey := strings.TrimPrefix(string(key), ectl.RootKey)
	ectl.mutex.Lock()
	defer ectl.mutex.Unlock()
	if len(value) == 0 {
		delete(ectl.kv, string(strKey))
	} else {
		ectl.kv[string(strKey)] = string(value)
	}
}
*/

func (ectl *EtcdConfig) set(key, value []byte) {
	strKey := strings.TrimPrefix(string(key), ectl.RootKey)
	ectl.mutex.Lock()
	defer ectl.mutex.Unlock()
	if len(value) == 0 {
		delete(ectl.kv, string(strKey))
	} else {
		ectl.kv[string(strKey)] = string(value)
	}
}

func (ectl *EtcdConfig) CheckupE3ch() ([]string, error) {
	var warns []string
	var err error

	if ectl.E3ch == nil {
		return nil, errors.New("E3ch client not initialized")
	}

	err = ectl.E3ch.FormatRootKey() // set the rootKey as directory
	if err != nil {
		warns = append(warns, fmt.Sprintf("- failed to  set the rootKey as directory, error: %s", err))
		return warns, err
	}

	// Quick Test
	err = ectl.E3ch.CreateDir("/dir1")
	if err != nil {
		warns = append(warns, fmt.Sprintf("- failed to CreateDir '/dir1', error: %s", err))
	}

	err = ectl.E3ch.Create("/dir1/key1", "")
	if err != nil {
		warns = append(warns, fmt.Sprintf("- failed to Create '/dir1/key1', error: %s", err))
	}

	err = ectl.E3ch.Create("/dir1/key2", "")
	if err != nil {
		warns = append(warns, fmt.Sprintf("- failed to Create '/dir1/key2', error: %s", err))
	}

	err = ectl.E3ch.Create("/dir2/key2", "")
	if err != nil {
		warns = append(warns, fmt.Sprintf("- failed to Create '/dir2/key2', error: %s", err))
	}

	err = ectl.E3ch.Create("/dir", "")
	if err != nil {
		warns = append(warns, fmt.Sprintf("- failed to Create '/dir', error: %s", err))
	}

	err = ectl.E3ch.Put("/dir1/key1", "value1")
	if err != nil {
		warns = append(warns, fmt.Sprintf("- failed to Put: key='/dir1/key1', error: %s", err))
	}

	// return node value
	var node *client.Node
	node, err = ectl.E3ch.Get("/dir1/key1")
	if err != nil {
		warns = append(warns, fmt.Sprintf("- failed to Get: key='/dir1/key1', error: %s", err))
	}
	if ectl.Debug {
		fmt.Println("- value for node='/dir1/key1':")
		pp.Println(parseNode(node))
	}
	// return nodes in dir
	var nodes []*client.Node
	nodes, err = ectl.E3ch.List("/")
	if err != nil {
		warns = append(warns, fmt.Sprintf("- failed to List keys: key='/dir1', error: %s", err))
	}
	if ectl.Debug {
		fmt.Println("- nodes for list='/dir1':")
		for k, node := range nodes {
			fmt.Printf("#%d: \n", k)
			pp.Println(parseNode(node))
		}
	}
	err = ectl.E3ch.Delete("/dir")
	if err != nil {
		warns = append(warns, fmt.Sprintf("- failed to Delete: key='/dir', error: %s", err))
	}
	_, err = ectl.E3ch.List("/")
	if err != nil {
		warns = append(warns, fmt.Sprintf("- failed to List: key='/', error: %s", err))
	}
	return warns, nil // return nil as no major
}

// GetValues queries etcd for keys prefixed by prefix.
func (ectl *EtcdConfig) GetValues(keys []string) (map[string]string, error) {
	vars := make(map[string]string)
	for _, key := range keys {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(3)*time.Second)
		resp, err := ectl.Client.Get(ctx, key, etcd.WithPrefix(), etcd.WithSort(etcd.SortByKey, etcd.SortDescend))
		cancel()
		if err != nil {
			return vars, err
		}
		for _, ev := range resp.Kvs {
			vars[string(ev.Key)] = string(ev.Value)
		}
	}
	return vars, nil
}

func (ectl *EtcdConfig) WatchPrefix(prefix string, keys []string, waitIndex uint64, stopChan chan bool) (uint64, error) {
	// return something > 0 to trigger a key retrieval from the store
	if waitIndex == 0 {
		return 1, nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancelRoutine := make(chan bool)
	defer close(cancelRoutine)
	var err error

	go func() {
		select {
		case <-stopChan:
			cancel()
		case <-cancelRoutine:
			return
		}
	}()

	rch := ectl.Client.Watch(ctx, prefix, etcd.WithPrefix())
	for wresp := range rch {
		for _, ev := range wresp.Events {
			fmt.Println(string(ev.Kv.Key))
			// Only return if we have a key prefix we care about.
			// This is not an exact match on the key so there is a chance
			// we will still pickup on false positives. The net win here
			// is reducing the scope of keys that can trigger updates.
			for _, k := range keys {
				if strings.HasPrefix(string(ev.Kv.Key), k) {
					return uint64(ev.Kv.Version), err
				}
			}
		}
	}
	return 0, err
}

func (ectl *EtcdConfig) AddEndpoint(path string, endpointConfig Endpoint) error {
	return nil
}

func (ectl *EtcdConfig) RecursiveCreateDir(keyPath string) error {
	keyParts := strings.Split(keyPath, "/")
	if len(keyParts) > ectl.MaxDir {
		return errors.New(fmt.Sprintf("[ERROR] Input path='%s', Max directory (%d) per key exceeded: '%d'.\n", keyPath, ectl.MaxDir, len(keyParts)))
	}
	for i := 0; i <= len(keyParts); i++ {
		if strings.Join(keyParts[:i], "/") != "" {
			fmt.Printf("input: '%s', iter='%d' , parent_dir: '%s'\n", keyPath, i, strings.Join(keyParts[:i], "/"))
			ectl.E3ch.CreateDir(strings.Join(keyParts[:i], "/"))
		}
	}
	return nil
}

type Node struct {
	Key   string `json:"key" yaml:"key" toml:"key"`
	Value string `json:"value" yaml:"value" toml:"value"`
	IsDir bool   `json:"is_dir" yaml:"is_dir" toml:"is_dir"`
}

func parseNode(node *client.Node) *Node {
	return &Node{
		Key:   string(node.Key),
		Value: string(node.Value),
		IsDir: node.IsDir,
	}
}

type info struct {
	field   reflect.Value
	version uint64
}

// Route configuration struct.
type EtcdRouteConfig struct {
	Regexp string `json:"regexp" yaml:"regexp" toml:"regexp"`
	Schema string `json:"schema" yaml:"schema" toml:"schema"`
}

func newEtcdCtx() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), ETCD_CLIENT_TIMEOUT)
	return ctx
}

func CloneE3chClient(username, password string, client *client.EtcdHRCHYClient) (*client.EtcdHRCHYClient, error) {
	ectl, err := etcd.New(etcd.Config{
		Endpoints: client.EtcdClient().Endpoints(),
		Username:  username,
		Password:  password,
	})
	if err != nil {
		return nil, err
	}
	return client.Clone(ectl), nil
}

func notFound(e error) bool {
	return e == rpctypes.ErrEmptyKey
}

type EtcdService struct {
	Disabled          bool              `etcd:"disabled" json:"disabled,omitempty" yaml:"disabled,omitempty" toml:"disabled,omitempty"`
	Name              string            `etcd:"name" json:"name,omitempty" yaml:"name,omitempty" toml:"name,omitempty"`
	HasHealthCheck    bool              `etcd:"has_health_check" json:"has_health_check,omitempty" yaml:"has_health_check,omitempty" toml:"has_health_check,omitempty"`
	Addresses         map[string]string `etcd:"addresses" json:"addresses,omitempty" yaml:"addresses,omitempty" toml:"addresses,omitempty"`
	PathPrefixes      map[string]string `etcd:"path_prefixes" json:"path_prefixes,omitempty" yaml:"path_prefixes,omitempty" toml:"path_prefixes,omitempty"`
	PathHosts         map[string]string `etcd:"path_hosts" json:"path_hosts,omitempty" yaml:"path_hosts,omitempty" toml:"path_hosts,omitempty"`
	FailoverPredicate string            `etcd:"failover_predicate" json:"failover_predicate,omitempty" yaml:"failover_predicate,omitempty" toml:"failover_predicate,omitempty"`
	Debug             bool              `etcd:"debug" json:"debug,omitempty" yaml:"debug,omitempty" toml:"debug,omitempty"`
}

type EtcdRoute struct {
	Loaded   bool `json:"-" yaml:"-" toml:"-"`
	Disabled bool `etcd:"disabled" json:"disabled,omitempty" yaml:"disabled,omitempty" toml:"disabled,omitempty"`
	// server - endpoint handler
	Source  string `etcd:"source" json:"provider,omitempty" yaml:"provider,omitempty" toml:"provider,omitempty"`
	Route   string `etcd:"router" json:"route,omitempty" yaml:"route,omitempty" toml:"route,omitempty"`
	Method  string `etcd:"method" json:"method,omitempty" yaml:"method,omitempty" toml:"method,omitempty"`
	Comment string `etcd:"comment" json:"comment,omitempty" yaml:"comment,omitempty" toml:"comment,omitempty"`
	// remote content - extraction rules and blocks
	BaseURL          string                       `required:"true" etcd:"base_url" json:"base_url,omitempty" yaml:"base_url,omitempty" toml:"base_url,omitempty"`
	PatternURL       string                       `required:"true" etcd:"pattern_url" json:"pattern_url" yaml:"pattern_url" toml:"pattern_url"`
	TestURL          string                       `required:"true" etcd:"test_url" json:"test_url" yaml:"test_url" toml:"test_url"`
	Selector         string                       `etcd:"selector" default:"css" json:"selector,omitempty" yaml:"selector,omitempty" toml:"selector,omitempty"`
	HeadersIntercept map[string]string            `etcd:"resp_headers_intercept" json:"resp_headers_intercept,omitempty" yaml:"resp_headers_intercept,omitempty" toml:"resp_headers_intercept,omitempty"`
	Headers          map[string]string            `etcd:"headers" json:"headers,omitempty" yaml:"headers,omitempty" toml:"headers,omitempty"`
	Blocks           map[string]map[string]string `etcd:"blocks" json:"blocks,omitempty" yaml:"blocks,omitempty" toml:"blocks,omitempty"`
	Extract          map[string]bool              `etcd:"extract" json:"extract,omitempty" yaml:"extract,omitempty" toml:"extract,omitempty"`
	Groups           string                       `etcd:"groups" json:"groups,omitempty" yaml:"groups,omitempty" toml:"groups,omitempty"`
	StrictMode       bool                         `etcd:"strict_mode" json:"strict_mode,omitempty" yaml:"strict_mode,omitempty" toml:"strict_mode,omitempty"`
	Debug            bool                         `etcd:"debug" json:"debug,omitempty" yaml:"debug,omitempty" toml:"debug,omitempty"`
}

type EtcdProxy struct {
	Disabled  bool                    `etcd:"disabled" json:"disabled,omitempty" yaml:"disabled,omitempty" toml:"disabled,omitempty"`
	Frontends map[string]EtcdFrontend `etcd:"frontends" json:"frontends,omitempty" yaml:"frontends,omitempty" toml:"frontends,omitempty"`
	Backends  map[string]EtcdBackend  `etcd:"backends" json:"backends,omitempty" yaml:"backends,omitempty" toml:"backends,omitempty"`
	// Middlewares map[string]EtcdMiddleware `etcd:"middlewares" json:"middlewares,omitempty" yaml:"middlewares,omitempty" toml:"middlewares,omitempty"`
	// Plugins     map[string]EtcdPlugin     `etcd:"plugins" json:"plugins,omitempty" yaml:"plugins,omitempty" toml:"plugins,omitempty"`
	Debug bool `etcd:"debug" json:"debug,omitempty" yaml:"debug,omitempty" toml:"debug,omitempty"`
}

type EtcdFrontend struct {
	Disabled          bool        `etcd:"disabled" json:"disabled,omitempty" yaml:"disabled,omitempty" toml:"disabled,omitempty"`
	BackendID         string      `etcd:"backend_id" json:"backend_id,omitempty" yaml:"backend_id,omitempty" toml:"backend_id,omitempty"`
	Route             string      `etcd:"route" json:"route,omitempty" yaml:"route,omitempty" toml:"route,omitempty"`
	Type              string      `etcd:"type" json:"type,omitempty" yaml:"type,omitempty" toml:"type,omitempty"`
	Rewrite           EtcdRewrite `etcd:"rewrite" json:"rewrite,omitempty" yaml:"rewrite,omitempty" toml:"rewrite,omitempty"`
	FailoverPredicate string      `etcd:"failover_predicate" json:"failover_predicate,omitempty" yaml:"failover_predicate,omitempty" toml:"failover_predicate,omitempty"`
	Debug             bool        `etcd:"debug" json:"debug,omitempty" yaml:"debug,omitempty" toml:"debug,omitempty"`
}

type EtcdRewrite struct {
	Disabled   bool          `etcd:"disabled" json:"disabled,omitempty" yaml:"disabled,omitempty" toml:"disabled,omitempty"`
	ID         string        `etcd:"id" json:"id,omitempty" yaml:"id,omitempty" toml:"id,omitempty"`
	Type       string        `etcd:"type" json:"type,omitempty" yaml:"type,omitempty" toml:"type,omitempty"`
	Priority   int           `etcd:"priority" json:"priority,omitempty" yaml:"priority,omitempty" toml:"priority,omitempty"`
	Middleware EtcdRewriteMw `etcd:"middleware" json:"middleware,omitempty" yaml:"middleware,omitempty" toml:"middleware,omitempty"`
	Debug      bool          `etcd:"debug" json:"debug,omitempty" yaml:"debug,omitempty" toml:"debug,omitempty"`
}

type EtcdRewriteMw struct {
	Disabled    bool   `etcd:"disabled" json:"disabled,omitempty" yaml:"disabled,omitempty" toml:"disabled,omitempty"`
	Regexp      string `etcd:"regexp" json:"regexp,omitempty" yaml:"regexp,omitempty" toml:"regexp,omitempty"`
	Replacement string `etcd:"replacement" json:"replacement,omitempty" yaml:"replacement,omitempty" toml:"replacement,omitempty"`
	Debug       bool   `etcd:"debug" json:"debug,omitempty" yaml:"debug,omitempty" toml:"debug,omitempty"`
}

type EtcdBackend struct {
	Disabled bool                  `etcd:"disabled" json:"disabled,omitempty" yaml:"disabled,omitempty" toml:"disabled,omitempty"`
	Servers  map[string]EtcdServer `etcd:"servers" json:"servers,omitempty" yaml:"servers,omitempty" toml:"servers,omitempty"`
	Debug    bool                  `etcd:"debug" json:"debug,omitempty" yaml:"debug,omitempty" toml:"debug,omitempty"`
}

type EtcdServer struct {
	Disabled bool   `etcd:"disabled" json:"disabled,omitempty" yaml:"disabled,omitempty" toml:"disabled,omitempty"`
	URL      string `etcd:"url" json:"url,omitempty" yaml:"url,omitempty" toml:"url,omitempty"`
	Debug    bool   `etcd:"debug" json:"debug,omitempty" yaml:"debug,omitempty" toml:"debug,omitempty"`
}

type EtcdHandler struct {
	Disabled   bool                    `etcd:"disabled" json:"disabled,omitempty" yaml:"disabled,omitempty" toml:"disabled,omitempty"`
	Backends   map[string]EtcdBackend  `etcd:"backends" json:"backends,omitempty" yaml:"backends,omitempty" toml:"backends,omitempty"`
	Frontends  map[string]EtcdFrontend `etcd:"frontends" json:"frontends,omitempty" yaml:"frontends,omitempty" toml:"frontends,omitempty"`
	Routes     map[string]EtcdRoute    `etcd:"routes" json:"routes,omitempty" yaml:"routes,omitempty" toml:"routes,omitempty"`
	Services   map[string]EtcdService  `etcd:"services" json:"services,omitempty" yaml:"services,omitempty" toml:"services,omitempty"`
	StrictMode bool                    `etcd:"strict_mode" json:"strict_mode,omitempty" yaml:"strict_mode,omitempty" toml:"strict_mode,omitempty"`
	Debug      bool                    `etcd:"debug" json:"debug,omitempty" yaml:"debug,omitempty" toml:"debug,omitempty"`
}

func (rc *EtcdHandler) NewEndpoint(path string, conf Endpoint) (bool, error) {
	fmt.Printf("Register new endpoint: %s \n", path)
	return true, nil
}

/*
func convertErr(e error) error {
	if e == nil {
		return nil
	}
	switch e {
	case rpctypes.ErrEmptyKey:
		return &engine.NotFoundError{Message: e.Error()}

	case rpctypes.ErrDuplicateKey:
		return &engine.AlreadyExistsError{Message: e.Error()}
	}
	return e
}

func (n ng) path(keys ...string) string {
	return strings.Join(append([]string{n.etcdKey}, keys...), "/")
}

func (n *ng) setJSONVal(key string, v interface{}, ttl time.Duration) error {
	bytes, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return n.setVal(key, bytes, ttl)
}

func (n *ng) setVal(key string, val []byte, ttl time.Duration) error {
	ops := []etcd.OpOption{}
	if ttl > 0 {
		lgr, err := n.client.Grant(n.context, int64(ttl.Seconds()))
		if err != nil {
			return err
		}
		ops = append(ops, etcd.WithLease(lgr.ID))
	}

	_, err := n.client.Put(n.context, key, string(val), ops...)
	return convertErr(err)
}

func (n *ng) getJSONVal(key string, in interface{}) error {
	val, err := n.getVal(key)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), in)
}

func (n *ng) getVal(key string) (string, error) {
	response, err := n.client.Get(n.context, key)
	if err != nil {
		return "", convertErr(err)
	}

	if len(response.Kvs) != 1 {
		return "", &engine.NotFoundError{Message: "Key not found"}
	}

	return string(response.Kvs[0].Value), nil
}

func (n *ng) getKeysBySecondPrefix(keys ...string) ([]string, error) {
	var out []string
	targetPrefix := strings.Join(keys, "/")
	response, err := n.client.Get(n.context, targetPrefix, etcd.WithPrefix(), etcd.WithSort(etcd.SortByKey, etcd.SortAscend))
	if err != nil {
		if notFound(err) {
			return out, nil
		}
		return nil, err
	}

	//If /this/is/prefix then
	// allow /this/is/prefix/one/two
	// disallow /this/is/prefix/one/two/three
	// disallow /this/is/prefix/one
	for _, keyValue := range response.Kvs {
		if prefix(prefix(string(keyValue.Key))) == targetPrefix {
			out = append(out, string(keyValue.Key))
		}
	}
	return out, nil
}

func (n *ng) getVals(keys ...string) ([]Pair, error) {
	var out []Pair
	response, err := n.client.Get(n.context, strings.Join(keys, "/"), etcd.WithPrefix(), etcd.WithSort(etcd.SortByKey, etcd.SortAscend))
	if err != nil {
		if notFound(err) {
			return out, nil
		}
		return nil, err
	}

	for _, keyValue := range response.Kvs {
		out = append(out, Pair{string(keyValue.Key), string(keyValue.Value)})
	}
	return out, nil
}

func (n *ng) checkKeyExists(key string) error {
	_, err := n.client.Get(n.context, key)
	return convertErr(err)
}

*/
