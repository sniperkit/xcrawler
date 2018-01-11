package scraper

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	// "net/rpc"
	// "sync"

	"github.com/BurntSushi/toml"
	"github.com/ahmetb/go-linq"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/birkelund/boltdbcache"
	"github.com/cabify/go-couchdb"
	bolt "github.com/coreos/bbolt"
	"github.com/fatih/color"
	"github.com/foize/go.fifo"
	"github.com/go-resty/resty"
	"github.com/gregjones/httpcache"
	"github.com/gregjones/httpcache/diskcache"
	"github.com/gregjones/httpcache/leveldbcache"
	"github.com/hydrogen18/stoppableListener"
	"github.com/if1live/staticfilecache"
	"github.com/joho/godotenv"
	"github.com/k0kubun/pp"
	"github.com/klaidliadon/go-couch-cache"
	"github.com/klaidliadon/go-memcached"
	"github.com/klaidliadon/go-redis-cache"
	"github.com/peterbourgon/diskv"
	"github.com/roscopecoltran/configor"
	"github.com/roscopecoltran/mxj"
	"github.com/trustmaster/goflow"
	"github.com/victortrac/disks3cache"
	ctx "golang.org/x/net/context"
	"golang.org/x/oauth2"
	"gopkg.in/redis.v3"
	// cache "github.com/patrickmn/go-cache"
	// "github.com/gin-gonic/gin"
	// "github.com/go-fsnotify/fsnotify"
	// "gopkg.in/olahol/melody.v1"
	// "golang.org/x/net/context/ctxhttp"
	// "github.com/gregjones/httpcache/memcache"
	// "github.com/mikegleasonjr/forwardcache"
)

var (
	BoltDB         *bolt.DB
	transportCache *httpcache.Transport
	// httpCache      httpcache.Cache
	// redisClient    *redis.Client
	// couchdbClient  *couchdb.Client
	// memcacheClient *memcache.Client
	// ErrURLEmpty to warn users that they passed an empty string in
	ErrURLEmpty = fmt.Errorf("the url string is empty")
	// ErrDomainMissing domain was missing from the url
	ErrDomainMissing = fmt.Errorf("url domain e.g .com, .net was missing")
	// ErrUnresolvedOrTimedOut ...
	ErrUnresolvedOrTimedOut = fmt.Errorf("url could not be resolved or timeout")

	// EmailRegex provides a base email regex for scraping emails
	EmailRegex = regexp.MustCompile(`([a-z0-9!#$%&'*+\/=?^_{|}~-]+(?:\.[a-z0-9!#$%&'*+\/=?^_{|}~-]+)*(@|\sat\s)(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?(\.|\sdot\s))+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?)`)

	searchTermColor = color.New(color.FgGreen).SprintFunc()
	foundColor      = color.New(color.FgGreen).SprintFunc()
	notFoundColor   = color.New(color.FgRed).SprintFunc()
	newLineReplacer = strings.NewReplacer("\r\n", "", "\n", "", "\r", "")
)

type Car struct {
	year         int
	owner, model string
}

var owners, cars []string

func testLinq() {

	linq.From(cars).Where(func(c interface{}) bool {
		return c.(Car).year >= 2015
	}).Select(func(c interface{}) interface{} {
		return c.(Car).owner
	}).ToSlice(&owners)
	pp.Println("owners: ", owners)

	linq.From(cars).WhereT(func(c Car) bool {
		return c.year >= 2015
	}).SelectT(func(c Car) string {
		return c.owner
	}).ToSlice(&owners)
	pp.Println("owners: ", owners)
}

// A graph for our app
type App struct {
	flow.Graph
}

type Handler struct {
	Disabled bool              `default:"false" help:"Disable handler init" json:"disabled,omitempty" yaml:"disabled,omitempty" toml:"disabled,omitempty"`
	Env      EnvConfig         `opts:"-" json:"env,omitempty" yaml:"env,omitempty" toml:"env,omitempty"`
	Etcd     EtcdConfig        `opts:"-" json:"etcd,omitempty" yaml:"etcd,omitempty" toml:"etcd,omitempty"`
	Config   Config            `opts:"-" json:"config,omitempty" yaml:"config,omitempty" toml:"config,omitempty"`
	Headers  map[string]string `opts:"-" json:"headers,omitempty" yaml:"headers,omitempty" toml:"headers,omitempty"`
	App      *App              `opts:"-" json:"app,omitempty" yaml:"app,omitempty" toml:"app,omitempty"`
	Auth     string            `help:"Basic auth credentials <user>:<pass>" json:"auth,omitempty" yaml:"auth,omitempty" toml:"auth,omitempty"`
	Log      bool              `default:"false" opts:"-" json:"log,omitempty" yaml:"log,omitempty" toml:"log,omitempty"`
	Debug    bool              `opts:"debug" default:"false" help:"Enable debug output" json:"debug,omitempty" yaml:"debug,omitempty" toml:"debug,omitempty"`
	Verbose  bool              `opts:"verbose" default:"false" help:"Enable verbose output" json:"verbose,omitempty" yaml:"verbose,omitempty" toml:"verbose,omitempty"`
	Cache    struct {
		Control int `opts:"-" default:"120" json:"control,omitempty" yaml:"control,omitempty" toml:"control,omitempty"`
	} `opts:"-" json:"cache,omitempty" yaml:"cache,omitempty" toml:"cache,omitempty"`
	// Api
	stoppable *stoppableListener.StoppableListener
	events    *fifo.Queue
	api       *rest.Api
}

func OpenBucket(filename string, bucketName string, permissions os.FileMode) error {
	if filename == "" {
		filename = "data.db"
	}
	if bucketName == "" {
		bucketName = "common"
	}
	if permissions == 0000 {
		permissions = 0666
	}
	db, err := bolt.Open(filename, permissions, nil)
	if err != nil {
		log.Printf("[open] init db error: %v", err)
		return err
	}
	BoltDB = db

	tx, err := BoltDB.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err = tx.CreateBucketIfNotExists([]byte(bucketName)); err != nil {
		return err
	}

	return tx.Commit()
}

func (h *Handler) LoadConfigToml(configFile string) error {
	c := Config{}
	_, err := toml.DecodeFile(configFile, &c)
	if err != nil {
		log.Fatal("Error loading config file")
		return err
	}
	h.Config = c
	InitCache("./shared/cache/external")
	return nil
}

func (h *Handler) LoadConfigorFile(path string) error {

	/*
		config := flag.String("file", "config.yml", "configuration file")
		flag.StringVar(&Config.APPName, "name", "", "app name")
		flag.StringVar(&Config.DB.Name, "db-name", "", "database name")
		flag.StringVar(&Config.DB.User, "db-user", "root", "database user")
		flag.Parse()
		os.Setenv("CONFIGOR_ENV_PREFIX", "-")
	*/

	if path == "" {
		return errors.New("no config file provided")
	}
	c := Config{}

	fmt.Println("h.Debug", h.Debug)
	fmt.Println("h.Verbose", h.Verbose)

	configor.New(&configor.Config{
		// Environment: "production",
		ENVPrefix: "SNIPERKIT",
		Debug:     h.Debug,
		Verbose:   h.Verbose,
	}).Load(&c, path)

	if h.Debug {
		pp.Println(c.Templates)
		fmt.Printf("config filepath: %s\n", path)
	}
	// os.Exit(1)

	if h.Debug {
		fmt.Println("configor loading: ")
		pp.Println(c)
		fmt.Printf("config filepath: %s\n", path)
		pp.Println(c.Templates)
		for k, v := range c.Templates {
			pp.Println("k=", k)
			pp.Println("v=", v)
		}
	}

	// bug with yaml format with details block; due to mapping of keys and values
	// os.Exit(1)

	h.Config = c
	InitCache("./shared/cache/external")
	return nil
}

func (h *Handler) LoadConfigFile(path string) error {
	fmt.Printf("config filepath: %s\n", path)
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return h.LoadConfig(b)
}

func (h *Handler) GetConfigPaths(path string) []string {
	var paths []string
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return paths
	}
	mxj.JsonUseNumber = true
	mv, err := mxj.NewMapJson(b)
	if err != nil {
		fmt.Println("NewMapJson, error: ", err)
		return paths
	}
	if h.Debug {
		fmt.Println("NewMapJson, jdata:", string(b))
		fmt.Printf("NewMapJson, mv: \n %#v\n", mv)
	}
	mxj.LeafUseDotNotation()
	paths = mv.LeafPaths()
	return paths
}

var Endpoints struct {
	Disabled bool
	Routes   []string
}

func (h *Handler) LoadConfig(b []byte) error {
	c := Config{}

	if err := json.Unmarshal(b, &c); err != nil { //json unmarshal performs selector validation
		return err
	}

	// Create a new graph
	h.App = new(App)
	h.App.InitGraphState()

	// Add graph nodes
	// h.App.Add(new(Router), "router")

	// Connect the processes
	// h.App.Connect("router", "Show", "controller", "In")

	// Network ports
	// h.App.MapInPort("In", "router", "In")

	h.Etcd = c.Etcd
	if len(c.Env.Files) > 0 {
		envVars, err := godotenv.Read(c.Env.Files...)
		if err != nil {
			return err
		}
		c.Env.VariablesList = envVars
		envVarsTree := make(map[string]map[string]string)
		for k, v := range envVars {
			var varParentKey, varChildrenKey string
			varParts := strings.Split(k, "_")
			if len(varParts) > 1 {
				varParentKey = varParts[0]
				varChildrenKey = strings.Join(varParts[1:], "_")
			}
			if v != "" && varParentKey != "" && varChildrenKey != "" {
				envVarsTree[varParentKey] = make(map[string]string)
				envVarsTree[varParentKey][varChildrenKey] = v
			}
		}
		c.Env.VariablesTree = envVarsTree
	}

	if h.Log {
		for k, e := range c.Routes {
			// Ovveride value ?! which cases ?!
			// e.Debug = h.Debug
			if strings.HasPrefix(e.Route, "/") {
				e.Route = strings.TrimPrefix(e.Route, "/")
				c.Routes[k] = e
			}
			/*
				if strings.HasPrefix(k, "/") {
					delete(c, k)
					k = strings.TrimPrefix(k, "/")
					c[k] = e
				}
			*/
			if h.Debug {
				logf("Loaded endpoint: /%s", e.Route)
			}
			Endpoints.Routes = append(Endpoints.Routes, e.Route)
			if len(h.Headers) > 0 && h.Debug { // Copy the Header attributes (only if they are not yet set)
				fmt.Printf("h.Headers, len=%d:\n", len(h.Headers))
				pp.Println(h.Headers)
			}
			for k, v := range e.HeadersJSON {
				if len(e.HeadersJSON) > 0 && h.Debug {
					pp.Println("header key: ", k)
					pp.Println("header val: ", v)
				}
				for kl, vl := range c.Env.VariablesList {
					holderKey := fmt.Sprintf("{{%s}}", strings.Replace(kl, "\"", "", -1))
					v = strings.Replace(v, holderKey, vl, -1)
				}
				e.HeadersJSON[k] = strings.Trim(v, " ")
			}
			if e.Crawler.MaxDepth <= 0 {
				// if e.Crawler.MaxDepth <= 0 {
				e.Crawler.MaxDepth = 1
			}
			if e.HeadersJSON == nil {
				e.HeadersJSON = h.Headers
			} else {
				for k, v := range h.Headers {
					if _, ok := e.HeadersJSON[k]; !ok {
						e.HeadersJSON[k] = v
					}
				}
			}
			if len(e.HeadersJSON) > 0 && h.Debug {
				fmt.Printf("e.HeadersJSON, len=%d:\n", len(e.HeadersJSON))
				pp.Println(e.HeadersJSON)
			}
		}
		for k, t := range c.Templates {
			pp.Println("k=", k, "t=", t)
			c.Templates[k] = t
		}
		pp.Println(c.Templates)
		// os.Exit(1)
	}
	if h.Debug {
		logf("Enabled debug mode")
	}

	h.Config = c //replace config
	InitCache("./shared/cache/external")
	return nil
}

func NewCache(cachePath string) *httpcache.Transport {
	return newTransportWithDiskCache(cachePath, "diskv")
}

func InitCache(cachePath string) {
	transportCache = newTransportWithDiskCache(cachePath, "diskv")
}

func newTransportContext(cachePath string, cacheEngine string) ctx.Context {
	transportCache = newTransportWithDiskCache(cachePath, cacheEngine)
	c := &http.Client{Transport: transportCache}
	return context.WithValue(context.Background(), oauth2.HTTPClient, c)
}

func newTransportWithDiskCache(basePath string, engine string) *httpcache.Transport {
	fmt.Println("[newTransportWithDiskCache] basePath: ", basePath)
	switch engine {
	case "boltdbcache":
		cachePath, err := boltdbcache.New(filepath.Join(basePath, "cache"))
		fmt.Println("[boltdbcache] cachePath: ", cachePath)
		if err != nil {
			fmt.Println("error: ", err)
		}
		return httpcache.NewTransport(cachePath)
	case "staticfilecache":
		cachePath := staticfilecache.New(basePath)
		fmt.Println("[staticfilecache] cachePath: ", cachePath)
		return httpcache.NewTransport(cachePath)
	case "disks3cache":
		cachePath, err := ioutil.TempDir("", "myTempDir")
		if err != nil {
			fmt.Println("error: ", err)
		}
		fmt.Println("[disks3cache] cachePath: ", cachePath)
		var cacheSize uint64
		cacheSize = 512 // in megabytes
		s3CacheURL := "s3://s3-us-west-2.amazonaws.com/my-bucket"
		cache := disks3cache.New(cachePath, cacheSize, s3CacheURL)
		return httpcache.NewTransport(cache)
	case "leveldbcache":
		cachePath, err := leveldbcache.New(filepath.Join(basePath, "cache"))
		if err != nil {
			fmt.Println("error: ", err)
		}
		fmt.Println("[leveldbcache] cachePath: ", cachePath)
		return httpcache.NewTransport(cachePath)
	case "couchdb":
		trans := &http.Transport{
			MaxIdleConnsPerHost: 10,
			Proxy:               http.ProxyFromEnvironment,
			Dial: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 60 * time.Second,
			}).Dial,
			TLSHandshakeTimeout: 10 * time.Second,
		}
		client, err := couchdb.NewClient("http://localhost:5984", trans)
		if err != nil {
			fmt.Println("error: ", err)
		}
		cache := couchcache.New(client.DB("cache"))
		cache.Indexes()
		return httpcache.NewTransport(cache)
	case "memcache":
		/*
			if memcachedURL := os.Getenv("MEMCACHE_URL"); memcachedURL != "" {
				return httpcache.NewTransport(memcache.New(memcachedURL, time.Minute*10)).Client()
			} else {
				return httpcache.NewTransport(httpcache.NewMemoryCache()).Client()
			}
		*/
		cache := memcached.New("localhost:11211", time.Minute*10)
		cache.Indexes()
		return httpcache.NewTransport(cache)
	case "redis":
		client := redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379"})
		cache := rediscache.New(client, time.Second/3)
		cache.Indexes()
		return httpcache.NewTransport(cache)
	case "diskv":
		fmt.Println("[diskv] basePath: ", basePath)
		d := diskv.New(diskv.Options{
			BasePath:     basePath,
			CacheSizeMax: 500 * 1024 * 250, // 10MB
		})
		cache := diskcache.NewWithDiskv(d)
		return httpcache.NewTransport(cache)
	}
	return nil
}

func getClient() *http.Client {
	c := transportCache.Client()
	// c.Timeout = time.Duration(30 * time.Second)
	// TODO Client Transport of type *httpcache.Transport doesn't support CanelRequest; Timeout not supported
	return c
}

func getResty() *resty.Client {
	transport := http.Transport{
		MaxIdleConns:        30,
		MaxIdleConnsPerHost: 30,
	}

	return resty.New().SetTransport(&transport).SetRetryCount(3).SetTimeout(time.Duration(25 * time.Second)).SetRedirectPolicy(resty.FlexibleRedirectPolicy(15))
}

func FaviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./shared/www/favicon.ico")
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// oauth2 ?!

	if h.Auth != "" { // basic auth
		u, p, _ := r.BasicAuth()
		if h.Auth != u+":"+p {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Access Denied"))
			return
		}
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8") // always JSON!
	// w.Header().Set("Content-Encoding", "gzip")
	// w.Header().Set("Cache-Control", "max-age=120")
	if r.URL.Path == "" || r.URL.Path == "/" { // admin actions
		get := false
		if r.Method == "GET" {
			get = true
		} else if r.Method == "POST" {
			b, err := ioutil.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write(jsonerr(err))
				return
			}
			if err := h.LoadConfig(b); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write(jsonerr(err))
				return
			}
			get = true
		}
		if !get {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write(jsonerr(errors.New("Use GET or POST")))
			return
		}
		b, _ := json.MarshalIndent(h.Config, "", "  ")
		w.Write(b)
		return
	}

	id := r.URL.Path[1:] // endpoint id (excludes root slash)

	pp.Println("r.URL.Path: ", r.URL.Path)
	pp.Println("r.URL.Path[1:]: ", r.URL.Path[1:])
	pp.Println("id: ", id)

	endpoint := h.Endpoint(id) // load endpoint

	if endpoint == nil {
		w.WriteHeader(404)
		w.Write(jsonerr(fmt.Errorf("Endpoint /%s not found", id)))
		return
	}
	values := map[string]string{} // convert url.Values into map[string]string
	for k, v := range r.URL.Query() {
		values[k] = v[0]
	}
	var err error
	res := make(map[string][]Result, 0)
	if endpoint.Debug {
		pp.Printf("endpoint.Concurrency: %s \n", endpoint.Concurrency)
	}

	_, _, cacheFile, err := endpoint.getCacheKey(r, h.Debug)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(jsonerr(err))
		return
	}

	isCacheExpired := cacheExpired(cacheFile, 5*time.Second)
	if endpoint.Debug && endpoint.Cache {
		fmt.Printf("[HANDLER] endpoint.Cache: %t\nisCacheExpired: %t\ncacheFile: %s \n", endpoint.Cache, isCacheExpired, cacheFile)
	}
	if !isCacheExpired && endpoint.Cache {
		if endpoint.Debug {
			fmt.Printf("reading cache content: %s \n", cacheFile)
		}
		file, err := os.Open(cacheFile)
		if err != nil {
			if endpoint.Debug {
				fmt.Println("os.Open, error: ", err)
			}
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(jsonerr(err))
			return
		}
		defer file.Close()

		b, err := ioutil.ReadAll(file)
		if err != nil {
			if endpoint.Debug {
				fmt.Println("ioutil.ReadAll, error: ", err)
			}
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(jsonerr(err))
			return
		}
		w.Header().Set("Content-Encoding", "gzip")
		w.Write(gzipFast(&b))
		return
	} else {
		if endpoint.Debug {
			fmt.Printf("new content to fetch/cache: %s \n", cacheFile)
		}
		if endpoint.Concurrency >= 1 && len(endpoint.Pager["max"]) > 0 {
			ctx := context.Background()
			resChan := make(chan *ScraperResult, endpoint.Concurrency)
			go endpoint.ExecuteParallel(ctx, values, resChan)
			totalResults, totalErrors := 0, 0
			for endpointResult := range resChan {
				if endpointResult.Error == nil {
					for k, v := range endpointResult.List {
						if _, ok := res[k]; !ok {
							res[k] = make([]Result, 0)
						}
						for _, r := range v {
							res[k] = append(res[k], r)
						}
						totalResults = totalResults + len(v)
					}
					if endpoint.Debug {
						fmt.Printf("res length: %d \n", len(res))
					}
				} else {
					totalErrors++
				}
			}
			if endpoint.Debug {
				fmt.Printf("totalResults: %d/%d, totalErrors: %d \n", totalResults, len(res["result"]), totalErrors)
			}
		} else {
			pp.Println("values: ", values)
			res, err = endpoint.Execute(values)
			if err != nil {
				fmt.Println("error endpoint.Execute(...): ", err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(jsonerr(err))
				return
			}
			if endpoint.Debug {
				fmt.Printf("totalResults: %d \n", len(res["result"]))
			}
		}
		if endpoint.Debug && endpoint.Cache {
			fmt.Println("[OUTPUT] isCacheExpired: ", isCacheExpired, ", cacheFile: ", cacheFile)
		}

		if len(res) > 0 && endpoint.Cache {
			err = cacheResponse(cacheFile, res) // dump response
			if err != nil {
				return
			}
			if endpoint.Debug {
				fmt.Printf("new content cached to file: %s\n", cacheFile)
			}
		}

		enc := json.NewEncoder(w) // encode as JSON
		enc.SetEscapeHTML(false)
		enc.SetIndent("", "  ")
		if err := enc.Encode(res); err != nil {
			w.Write([]byte("JSON Error: " + err.Error()))
		}

		/*
			writer, err := gzip.NewWriterLevel(w, gzip.BestCompression)
			if err != nil {
				// Your error handling
				return
			}

			defer writer.Close()
			writer.Write(enc)
		*/

		/*
			var v interface{}
			if endpoint.List == "" && len(res) == 1 {
				v = res[0]
			} else {
				v = res
			}
			if err := enc.Encode(v); err != nil {
				w.Write([]byte("JSON Error: " + err.Error()))
			}
		*/
	}
	// fmt.Fprintf(w, "luc")
}

// Endpoint will return the Handler's Endpoint from its Config
func (h *Handler) Endpoint(path string) *Endpoint {
	var keyCfg int
	//if !strings.HasPrefix(path, "/") {
	//	path = fmt.Sprintf("/%s", path)
	//}
	fmt.Printf("path to match: %s\n", path)
	// pp.Println("h.Config.Routes: ", h.Config.Routes)
	for k, v := range h.Config.Routes {
		fmt.Printf("v.Route: %s, path: %s\n", v.Route, path)
		if v.Route == path {
			fmt.Println("v.Route == path, k=", k)
			keyCfg = k
			break
		}
	}
	if h.Config.Routes[keyCfg] != nil {
		if h.Config.Routes[keyCfg].Template != "" {
			for k, v := range h.Config.Templates {
				if k == h.Config.Routes[keyCfg].Template {
					h.Config.Routes[keyCfg].BlocksJSON = v
				}
			}
		}
		return h.Config.Routes[keyCfg]
	}
	fmt.Println("v.Route is nil ")
	return nil
}
