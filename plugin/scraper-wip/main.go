package main

import (
	// _ "expvar"
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-gonic/gin"
	"github.com/go-fsnotify/fsnotify"
	"github.com/googollee/go-socket.io"
	"github.com/jamisonhyatt/HttpParallelSync"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/jpillora/opts"
	"github.com/k0kubun/pp"
	"github.com/mholt/caddy"
	"github.com/mholt/caddy/caddyhttp/httpserver"
	"github.com/roscopecoltran/admin"
	"github.com/roscopecoltran/scraper/scraper"

	"github.com/wantedly/gorm-zap"
	"go.uber.org/zap"
	"gopkg.in/natefinch/lumberjack.v2"
	"gopkg.in/olahol/melody.v1"

	"github.com/geekypanda/httpcache"
	"gopkg.in/unrolled/secure.v1"
	//"github.com/soheilhy/cmux"
	// "github.com/stretchr/graceful"

	krakendcfg "github.com/devopsfaith/krakend/config"
	krakendlog "github.com/devopsfaith/krakend/logging"
	"github.com/devopsfaith/krakend/proxy"
	"github.com/devopsfaith/krakend/router/mux"
	// ref. https://github.com/devopsfaith/krakend/blob/master/examples/httpcache/main.go
	// "github.com/gin-gonic/contrib/cache"
	// "github.com/gin-gonic/gin"
	// "github.com/gregjones/httpcache"
	// ref. https://github.com/hacdias/filemanager/blob/master/cmd/filemanager/main.go
	// "github.com/asdine/storm"
	// "github.com/hacdias/filemanager"
	// "github.com/hacdias/filemanager/bolt"
	// h "github.com/hacdias/filemanager/http"
	// "github.com/hacdias/filemanager/staticgen"
	// "github.com/hacdias/fileutils"
	// "github.com/tleyden/open-ocr"
	// fsnotify "gopkg.in/fsnotify.v1"
	// "github.com/valyala/fasthttp"
	// "github.com/geekypanda/httpcache"
	// "github.com/meission/router"
	// "github.com/go-zoo/bone"
	// "github.com/birkelund/boltdbcache"
	// "golang.org/x/oauth2"
	// "github.com/mickep76/flatten"
	// "github.com/gin-contrib/cache"
	// "github.com/aviddiviner/gin-limit"
	// "github.com/gin-gonic/contrib/cache"
	// "github.com/gin-gonic/contrib/secure"
	// "github.com/gin-gonic/contrib/static"
	// "github.com/ashwanthkumar/slack-go-webhook"
	// "github.com/carlescere/scheduler"
	// "github.com/jungju/qor_admin_auth"
	// "github.com/qor/publish2"
	// "github.com/qor/validations"
	// "golang.org/x/crypto/bcrypt"
	// "github.com/roscopecoltran/scraper/db/redis"
	// "github.com/roscopecoltran/scraper/api"
	// https://github.com/dyllanwli/GoLang_project/blob/master/blockchain/main.go
	// https://github.com/aliostad/deep-learning-lang-detection/blob/1180fba0d2a7f6b470cb3c9a363b560787f5e7c5/data/test/go/ec5f82a852d053a084edbc39ac4b56f9381b7cf9test.go
)

/*
	Refs:
	- https://github.com/agencyrevolution/go-microservices-example/blob/master/utils/vulcand.go
	- https://github.com/helinwang/gotensor/blob/master/cmd/server/main.go#L31-L37
	- fastText
		1. Use gensim, a python topic modeling library.
		2. Convert the fastText model file to gensim model file using this python code.
		from gensim.models.word2vec import Word2Vec
		from gensim.models.wrapper import FastText
		model = FastText.load_fasttext_format('wiki.en')
		model = Word2Vec.save("wiki.en-gensim")
		3. Load your gensim model.
		from gensim.models.word2vec import Word2Vec
		model = Word2Vec.load("wiki-en-gensim")
		- Using gensim model, you can reduce model size 35x and achieve fast model loading. (978MB -> 27.8MB)

*/

var VERSION = "0.0.0"

type config struct {
	*scraper.Handler `type:"embedded"`

	ConfigFile string `type:"arg" help:"Path to JSON configuration file" json:"config_file" yaml:"config_file" toml:"config_file"`
	Host       string `default:"0.0.0.0" help:"Listening interface" json:"host" yaml:"host" toml:"host"`
	Port       int    `default:"8092" help:"Listening port" json:"port" yaml:"port" toml:"port"`
	NoLog      bool   `default:"false" help:"Disable access logs" json:"logs" yaml:"logs" toml:"logs"`

	EtcdHost string `default:"etcd-1,etcd-2" help:"Listening interface" json:"etcd_host" yaml:"etcd_host" toml:"etcd_host"`
	EtcdPort int    `default:"2379" help:"Listening port" json:"etcd_port" yaml:"etcd_port" toml:"etcd_port"`

	RedisAddr string `default:"127.0.0.1:6379" help:"Redis Addr" json:"redis_addr" yaml:"redis_addr" toml:"redis_addr"`
	RedisHost string `default:"127.0.0.1" help:"Redis host" json:"redis_host" yaml:"redis_host" toml:"redis_host"`
	RedisPort string `default:"6379" help:"Redis port" json:"redis_port" yaml:"redis_port" toml:"redis_port"`

	// redis.UseRedis(rhost)
}

var (
	// Serialize all modifications through these
	// commands chan interface{}
	// errors   chan error

	// Clients
	// clients     []sockjs.Session
	// clientsLock sync.RWMutex
	// Index     bleve.Index

	// handler http.Handler
	// indexLock sync.RWMutex
	// Signalled to exit everything
	// finish chan struct{}

	// Used to control when things are done
	// wg sync.WaitGroup

	AdminUI *admin.Admin
	DB      *gorm.DB

	Tables = []interface{}{
		&scraper.Connection{},
		&scraper.Request{},
		&scraper.Response{},
		&scraper.Screenshot{},
		&scraper.Matcher{},
		&scraper.Queries{},
		&scraper.ProviderWebRankConfig{},
		&scraper.MatcherConfig{},
		&scraper.TargetConfig{},
		&scraper.Provider{},
		&scraper.Group{},
		&scraper.Topic{},
		&scraper.Endpoint{},
		&scraper.SelectorType{},
		&scraper.ExtractorsConfig{},
		&scraper.BlocksConfig{},
		&scraper.HeaderConfig{},
		&scraper.SelectorConfig{},
		&scraper.ExtractConfig{},
		&scraper.OpenAPIConfig{},
		&scraper.OpenAPISpecsConfig{},
	}

	logger        *zap.Logger
	krakendlogger krakendlog.Logger
	/*
	   log.SetOutput(&lumberjack.Logger{
	       Filename:   "./shared/logs/http_parallel_sync/scraper.log",
	       MaxSize:    500, // megabytes
	       MaxBackups: 3,
	       MaxAge:     28, //days
	   })
	*/

	errInit error
)

var cacheDuration = 3600 * time.Second
var DefaultCaddyServeMux = http.NewServeMux()

/*
func setup(c *caddy.Controller) error {
	//func setup(c *setup.Controller) error {
	// (middleware.Middleware, error) {
	//return func(next middleware.Handler) middleware.Handler {
	//	return &handler{}
	//}, nil
	httpserver.GetConfig(c.Key).AddMiddleware(func(next httpserver.Handler) httpserver.Handler {
		return MuxHandler{Next: next}
	})
	return nil
}
*/

func setup(c *caddy.Controller) error {
	cnf := httpserver.GetConfig(c)
	for c.Next() {
		if !c.NextArg() { // expect at least one value
			return c.ArgErr() // otherwise it's an error
		}
		value := c.Val() // use the value
		fmt.Println(value)
	}
	mid := func(next httpserver.Handler) httpserver.Handler {
		return &MuxHandler{
			Next: next,
		}
	}

	cnf.AddMiddleware(mid)
	return nil
}

/*

func Setup(c *setup.Controller) (middleware.Middleware, error) {
	return func(next middleware.Handler) middleware.Handler {
		return &handler{}
	}, nil
}

type handler struct{}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) (int, error) {
	w.Write([]byte("Hello, I'm a caddy middleware"))
	return 0, nil
}
*/

func init() {
	caddy.RegisterPlugin("mux", caddy.Plugin{
		ServerType: "http",
		Action:     setup,
	})
}

type MuxHandler struct {
	Next httpserver.Handler
}

type TelegramHandler struct {
	Next httpserver.Handler
}

// Get the default Caddy ServeMux
func ServeMux() *http.ServeMux {
	return DefaultCaddyServeMux
}

// Register the handler for the given pattern in the default Caddy ServeMux
func Handle(pattern string, handler http.Handler) {
	DefaultCaddyServeMux.Handle(pattern, handler)
}

// Registers the handler function for the given pattern in the default Caddy ServeMux
func HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	DefaultCaddyServeMux.HandleFunc(pattern, handler)
}

func (m MuxHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) (int, error) {
	_, pattern := DefaultCaddyServeMux.Handler(r)
	if len(pattern) > 0 {
		DefaultCaddyServeMux.ServeHTTP(w, r)
		return 0, nil
	} else {
		// no matching filter
		return m.Next.ServeHTTP(w, r)
	}
}

func httpParallelSyncMain(cluster string) {

	if cluster == "" {
		cluster = "test"
	}

	log.SetOutput(&lumberjack.Logger{
		Filename:   fmt.Sprintf("./shared/logs/http_parallel_sync/%s.log", cluster),
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28, //days
	})

	client := newCaddyClient()

	err := HttpParallelSync.Sync(client, cluster, 2)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("complete")
}

func newCaddyClient() *HttpParallelSync.CaddyClient {
	caddy := HttpParallelSync.CaddyClient{
		Host: "localhost",
		Port: 2015,
		Ssl:  false,
	}

	caddy.HttpClient = &http.Client{
		Timeout: time.Second * 25,
	}
	var protocol string
	if caddy.Ssl {
		protocol = "https"
	} else {
		protocol = "http"
	}
	caddy.BaseURI = fmt.Sprintf("%s://%s:%v", protocol, caddy.Host, caddy.Port)
	return &caddy
}

// import "github.com/devopsfaith/krakend/config/viper"
func newKrakendMux(serviceConfig krakendcfg.ServiceConfig, logLevel string) {

	/*
		port := flag.Int("p", 0, "Port of the service")
		logLevel := flag.String("l", "ERROR", "Logging level")
		debug := flag.Bool("d", false, "Enable the debug")
		configFile := flag.String("c", "/etc/krakend/configuration.json", "Path to the configuration filename")
		flag.Parse()

		parser := viper.New()
		serviceConfig, err := parser.Parse(*configFile)
		if err != nil {
			log.Fatal("ERROR:", err.Error())
		}
		serviceConfig.Debug = serviceConfig.Debug || *debug
		if *port != 0 {
			serviceConfig.Port = *port
		}
	*/

	if logLevel == "" {
		logLevel = "ERROR"
	}

	klogger, err := krakendlog.NewLogger(logLevel, os.Stdout, "[SCRAPER]")
	if err != nil {
		log.Fatal("ERROR: ", err.Error())
	}
	krakendlogger = klogger

	secureMiddleware := secure.New(secure.Options{
		AllowedHosts:          []string{"127.0.0.1:8080", "example.com", "ssl.example.com"},
		SSLRedirect:           false,
		SSLHost:               "ssl.example.com",
		SSLProxyHeaders:       map[string]string{"X-Forwarded-Proto": "https"},
		STSSeconds:            315360000,
		STSIncludeSubdomains:  true,
		STSPreload:            true,
		FrameDeny:             true,
		ContentTypeNosniff:    true,
		BrowserXssFilter:      true,
		ContentSecurityPolicy: "default-src 'self'",
	})

	// routerFactory := mux.DefaultFactory(proxy.DefaultFactory(logger), logger)

	routerFactory := mux.NewFactory(mux.Config{
		Engine:       mux.DefaultEngine(),
		ProxyFactory: proxy.DefaultFactory(klogger),
		Middlewares:  []mux.HandlerMiddleware{secureMiddleware},
		Logger:       krakendlogger,
		HandlerFactory: func(cfg *krakendcfg.EndpointConfig, p proxy.Proxy) http.HandlerFunc {
			return httpcache.CacheFunc(mux.EndpointHandler(cfg, p), time.Minute)
		},
	})

	routerFactory.New().Run(serviceConfig)
}

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	useGinWrap := false

	logger, errInit = zap.NewProduction()

	h := &scraper.Handler{Log: true}
	c := config{
		Handler: h,
		Host:    "0.0.0.0",
		Port:    3000,
	}

	opts.New(&c).
		Repo("github.com/roscopecoltran/scraper").
		Version(VERSION).
		Parse()

	h.Log = !c.NoLog
	go func() {
		for {
			sig := make(chan os.Signal, 1)
			signal.Notify(sig, syscall.SIGHUP)
			<-sig
			//if err := h.LoadConfigorFile(c.ConfigFile); err != nil {
			if err := h.LoadConfigFile(c.ConfigFile); err != nil {
				log.Printf("[scraper] Failed to load configuration: %s", err)
			} else {
				log.Printf("[scraper] Successfully loaded new configuration")
			}
		}
	}()

	// alternative: https://github.com/Entalpi/H-News-Backend/blob/master/cmd/main.go
	/*
		// When closed make sure to call Close on all the underlying bolt.DB instances.
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGTERM, os.Interrupt)
		go func() {
			<-ch
			newestScraper.DatabaseService.Close()
			services.Commentsdb.Close()
			os.Exit(1)
		}()

		select {} // Block forever and ever
	*/

	m := melody.New()
	w, _ := fsnotify.NewWatcher()
	go func() {
		for {
			ev := <-w.Events
			if ev.Op == fsnotify.Write {
				content, _ := ioutil.ReadFile(ev.Name)
				fmt.Println("ev.Name:", ev.Name)
				fmt.Printf("content: %s\n", content)
				m.Broadcast(content)

				if err := h.LoadConfigFile(c.ConfigFile); err != nil {
					log.Printf("[scraper] Failed to load configuration: %s", err)
				} else {
					log.Printf("[scraper] Successfully loaded new configuration")
				}

			}
		}
	}()

	//if err := h.LoadConfigorFile(c.ConfigFile); err != nil {
	if err := h.LoadConfigFile(c.ConfigFile); err != nil {
		log.Fatal(err)
	}

	m.HandleConnect(func(s *melody.Session) {
		content, _ := ioutil.ReadFile(c.ConfigFile)
		s.Write(content)
		fmt.Println("c.ConfigFile:", c.ConfigFile)
		fmt.Printf("content: %s\n", content)
	})
	w.Add(c.ConfigFile)

	// cache
	// https://github.com/garycarr/httpcache/blob/f039dd6ff44cf40d52e8e86ef10bff41e592fd48/README.md
	fmt.Printf("Scraper.NumCPU: %d\n", runtime.NumCPU())
	fmt.Printf("Scraper.useGinWrap: %t\n", useGinWrap)
	fmt.Printf("Scraper.Etcd.Disabled? %t \n", h.Etcd.Disabled)
	fmt.Printf("Scraper.Etcd.InitCheck? %t \n", h.Etcd.InitCheck)
	fmt.Printf("Scraper.Etcd.Debug? %t \n", h.Etcd.Debug)
	e3ch, err := c.Etcd.NewE3chClient()
	if err != nil {
		fmt.Println("Could not connect to the ETCD cluster, error: ", err)
	}

	if e3ch != nil {
		h.Etcd.Handler = h
		h.Etcd.E3ch = e3ch
	}

	// if useBoneMux
	// mux := bone.New()
	mux := http.NewServeMux() // Register route

	/*
	   // mux.Get, Post, etc ... takes http.Handler
	   mux.Get("/home/:id", http.HandlerFunc(HomeHandler))
	   mux.Get("/profil/:id/:var", http.HandlerFunc(ProfilHandler))
	   mux.Post("/data", http.HandlerFunc(DataHandler))

	   // Support REGEX Route params
	   mux.Get("/index/#id^[0-9]$", http.HandlerFunc(IndexHandler))

	   // Handle take http.Handler
	   mux.Handle("/", http.HandlerFunc(RootHandler))
	*/

	if h.Config.Debug {
		fmt.Printf(" - IsLogger? %t \n", h.Log)
		fmt.Println(" - Env params: ")
		pp.Println(h.Config.Env.VariablesTree)
	}

	if err := scraper.OpenBucket("bucket.db", "scraper", 0666); err != nil {
		log.Fatal(err)
	}

	if h.Config.Migrate {
		if h.Config.Debug {
			fmt.Printf(" - IsTruncateTables? %t \n", h.Config.Truncate)
			fmt.Printf(" - IsMigrateEndpoints? %t \n", h.Config.Migrate)
		}
		DB, errInit = gorm.Open("sqlite3", "admin.db")
		if errInit != nil {
			panic("failed to connect database")
		}
		defer DB.Close()

		if h.Config.Debug {
			DB.LogMode(true)
			if errInit == nil {
				DB.SetLogger(gormzap.New(logger))
			}
		}

		scraper.MigrateTables(DB, h.Config.Truncate, Tables...) // Create RDB datastore
	}

	if h.Config.Dashboard {
		initDashboard()
		AdminUI.MountTo("/admin", mux) // amount to /admin, so visit `/admin` to view the admin interface
	}

	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}

	// https://github.com/googollee/go-socket.io/blob/master/example/main.go
	server.On("connection", func(so socketio.Socket) {
		log.Println("on connection")
		so.Join("chat")
		so.On("chat message", func(msg string) {
			fmt.Println(so, msg)
			log.Println("emit:", so.Emit("chat message", msg))
			so.BroadcastTo("chat", "chat message", msg)
		})
		so.On("disconnection", func() {
			log.Println("on disconnect")
		})
	})
	server.On("error", func(so socketio.Socket, err error) {
		log.Println("error:", err)
	})

	// Experimental
	// redis.UseRedis(c.RedisHost)
	// scraper.ConvertToJsonSchema()
	// scraper.SeedAlexaTop1M()
	// h = scraper.NewRequestCacher(mux, "./shared/cache/scraper")
	mux.Handle("/", h)
	//	mux.HandleFunc("/ws", m.HandleRequest())
	mux.Handle("/socket.io/", server)

	mux.HandleFunc("/favicon.ico", scraper.FaviconHandler)
	mux.HandleFunc("/test", handler)

	api_handler, err := setupApi()
	if err == nil {
		mux.Handle("/api/v2", api_handler)
	}

	// GetFunc, PostFunc etc ... takes http.HandlerFunc
	// mux.GetFunc("/test", Handler)

	mux.HandleFunc("/api/v1", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	if h.Config.Migrate {
		scraper.MigrateEndpoints(DB, h.Config, e3ch)
	}

	if useGinWrap { // With GIN

		gin.SetMode(gin.ReleaseMode)
		r := gin.Default()
		store := persistence.NewInMemoryStore(60 * time.Second)
		if h.Config.Debug {
			fmt.Println("store: ")
			pp.Println(store)
		}

		r.Any("/*w", gin.WrapH(mux))
		if err := r.Run(fmt.Sprintf("%s:%d", c.Host, c.Port)); err != nil {
			log.Fatalf("Can not run server, error: %s", err)
		}

	} else {

		log.Printf("Listening on: %s:%d", c.Host, c.Port)
		log.Fatal(http.ListenAndServe(c.Host+":"+strconv.Itoa(c.Port), mux))

	}

}

func setupApi() (http.Handler, error) {

	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		rest.Get("/countries", GetAllCountries),
		rest.Post("/countries", PostCountry),
		rest.Get("/countries/:code", GetCountry),
		rest.Delete("/countries/:code", DeleteCountry),
	)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	api.SetApp(router)
	return api.MakeHandler(), nil
}

type Country struct {
	Code string
	Name string
}

var store = map[string]*Country{}

var lock = sync.RWMutex{}

func GetCountry(w rest.ResponseWriter, r *rest.Request) {
	code := r.PathParam("code")

	lock.RLock()
	var country *Country
	if store[code] != nil {
		country = &Country{}
		*country = *store[code]
	}
	lock.RUnlock()

	if country == nil {
		rest.NotFound(w, r)
		return
	}
	w.WriteJson(country)
}

func GetAllCountries(w rest.ResponseWriter, r *rest.Request) {
	lock.RLock()
	countries := make([]Country, len(store))
	i := 0
	for _, country := range store {
		countries[i] = *country
		i++
	}
	lock.RUnlock()
	w.WriteJson(&countries)
}

func PostCountry(w rest.ResponseWriter, r *rest.Request) {
	country := Country{}
	err := r.DecodeJsonPayload(&country)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if country.Code == "" {
		rest.Error(w, "country code required", 400)
		return
	}
	if country.Name == "" {
		rest.Error(w, "country name required", 400)
		return
	}
	lock.Lock()
	store[country.Code] = &country
	lock.Unlock()
	w.WriteJson(&country)
}

func DeleteCountry(w rest.ResponseWriter, r *rest.Request) {
	code := r.PathParam("code")
	lock.Lock()
	delete(store, code)
	lock.Unlock()
	w.WriteHeader(http.StatusOK)
}

// "handler" is our handler function. It has to follow the function signature of a ResponseWriter and Request type
// as the arguments.
func handler(w http.ResponseWriter, r *http.Request) {
	// For this case, we will always pipe "Hello World" into the response writer
	fmt.Fprintf(w, "Hello World!")
}

/*
// Handler return a http.Handler that supports Vue Router app with history mode
func fileHandler(publicDir string) http.Handler {
	handler := http.FileServer(http.Dir(publicDir))

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		url := req.URL.String()

		// static files
		if strings.Contains(url, ".") || url == "/" {
			handler.ServeHTTP(w, req)
			return
		}

		// the all 404 gonna be served as root
		http.ServeFile(w, req, path.Join(publicDir, "/index.html"))
	})
}
*/

/*
	// ref. https://github.com/tleyden/open-ocr/blob/master/cli-worker/main.go
	// ref. https://github.com/tleyden/open-ocr/blob/master/cli-httpd/main.go
	// This assumes that there is a worker running
	// To test it:
	// curl -X POST -H "Content-Type: application/json" -d '{"img_url":"http://localhost:8081/img","engine":0}' http://localhost:8081/ocr
	rabbitConfig := ocrworker.DefaultConfigFlagsOverride(flagFunc)

	// any requests to root, just redirect to main page
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		text := `<h1>OpenOCR is running!<h1> Need <a href="http://www.openocr.net">docs</a>?`
		fmt.Fprintf(w, text)
	})

	http.Handle("/ocr", ocrworker.NewOcrHttpHandler(rabbitConfig))

	http.Handle("/ocr-file-upload", ocrworker.NewOcrHttpMultipartHandler(rabbitConfig))

	// add a handler to serve up an image from the filesystem.
	// ignore this, was just something for testing ..
	http.HandleFunc("/img", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../refactoring.png")
	})
*/

/*
func semaphoreTimeout() {
	sla := 100 * time.Millisecond
	sem := semaphore.New(1000)

	http.Handle("/do-with-timeout", http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		done := make(chan struct{})
		deadline := semaphore.WithTimeout(sla)

		go func() {
			release, err := sem.Acquire(deadline)
			if err != nil {
				return
			}
			defer release()
			defer close(done)

			// do some heavy work
		}()

		// wait what happens before
		select {
		case <-deadline:
			http.Error(rw, "operation timeout", http.StatusGatewayTimeout)
		case <-done:
			// send success response
		}
	}))
}

func semaphoreContextCancel() {
	deadliner := func(limit int, timeout time.Duration, handler http.HandlerFunc) http.HandlerFunc {
		throughput := semaphore.New(limit)
		return func(rw http.ResponseWriter, req *http.Request) {
			ctx := semaphore.WithContext(req.Context(), semaphore.WithTimeout(timeout))

			release, err := throughput.Acquire(ctx.Done())
			if err != nil {
				http.Error(rw, err.Error(), http.StatusGatewayTimeout)
				return
			}
			defer release()

			handler.ServeHTTP(rw, req.WithContext(ctx))
		}
	}

	http.HandleFunc("/do-with-deadline", deadliner(1000, time.Minute, func(rw http.ResponseWriter, req *http.Request) {
		// do some limited work
	}))
}

*/

/*
// import "github.com/lhside/chrome-go"
func chromeBridge() {
	// Read message from standard input.
	msg, err := chrome.Receive(os.Stdin)
	// Post message to standard output
	err := chrome.Post(msg, os.Stdout)
}

// import "github.com/sauyon/go-chromemessage/chromemsg"
func chromeBridge2() {
	msg := chromemsg.New()
	msg.Read()
	msg.Write()
}
*/
