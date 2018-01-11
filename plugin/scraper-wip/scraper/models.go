package scraper

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/qor/media/media_library"
	"github.com/qor/sorting"
	"github.com/qor/validations"
	// "github.com/ucirello/goherokuname"
	// nlp "github.com/chewxy/lingo"
	// spacy "github.com/peter3125/goparser"
	// kai "github.com/peter3125/k-ai"
	// ""
)

// https://news.google.com/news/rss/headlines/section/topic/NATION/?ned=us&hl=en
// https://news.google.com/news/headlines?ned=us&hl=en
// url := "https://www.yahoo.com/news/rss/mostviewed"

// WEB SCRAPER ///////////////////////////////////////////////////////////////

type Result map[string]interface{} // Result represents a result

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

// Create a GORM-backend model
type Provider struct {
	gorm.Model      `json:"-" yaml:"-" toml:"-"`
	sorting.Sorting `json:"-" yaml:"-" toml:"-"`
	// ProviderID uint
	Name  string                   `etcd:"name" required:"true" json:"name" yaml:"name" toml:"name"` // gorm:"type:varchar(128);unique_index"
	Logo  media_library.MediaBox   `json:"-" yaml:"-" toml:"-"`
	Ranks []*ProviderWebRankConfig `json:"ranks,omitempty" yaml:"ranks,omitempty" toml:"ranks,omitempty"`
	// Endpoints []*Endpoint              `json:"endpoints,omitempty" yaml:"endpoints,omitempty" toml:"endpoints,omitempty"`
}

type ProviderWebRankConfig struct {
	gorm.Model      `json:"-" yaml:"-" toml:"-"`
	sorting.Sorting `json:"-" yaml:"-" toml:"-"`
	ProviderID      uint   `json:"-" yaml:"-" toml:"-"`
	Engine          string `json:"engine,omitempty" yaml:"engine,omitempty" toml:"engine,omitempty"`
	Score           string `json:"score,omitempty" yaml:"score,omitempty" toml:"score,omitempty"`
}

// Config represents...
type Config struct {
	gorm.Model      `json:"-" yaml:"-" toml:"-"`
	sorting.Sorting `json:"-" yaml:"-" toml:"-"`
	Disabled        bool        `default:"false" help:"Disable handler init" json:"disabled,omitempty" yaml:"disabled,omitempty" toml:"disabled,omitempty"`
	Mode            string      `default:"dev" help:"Env mode" json:"mode,omitempty" yaml:"mode,omitempty" toml:"mode,omitempty"`
	Cache           CacheConfig `gorm:"-" help:"Cache handler opts" json:"cache,omitempty" yaml:"cache,omitempty" toml:"cache,omitempty"`
	Env             EnvConfig   `gorm:"-" json:"env,omitempty" yaml:"env,omitempty" toml:"env,omitempty"`
	Etcd            EtcdConfig  `opts:"-" json:"etcd,omitempty" yaml:"etcd,omitempty" toml:"etcd,omitempty"`
	Port            int         `default:"3000" json:"port,omitempty" yaml:"port,omitempty" toml:"port,omitempty"`
	Dashboard       bool        `default:"false" help:"Initialize the Administration Interface" json:"dashboard,omitempty" yaml:"dashboard,omitempty" toml:"dashboard,omitempty"`
	Truncate        bool        `default:"true" help:"Truncate previous data" json:"truncate,omitempty" yaml:"truncate,omitempty" toml:"truncate,omitempty"`
	Migrate         bool        `default:"true" help:"Migrate to admin dashboard" json:"migrate,omitempty" yaml:"migrate,omitempty" toml:"migrate,omitempty"`
	Debug           bool        `default:"false" help:"Enable debug output" json:"debug,omitempty" yaml:"debug,omitempty" toml:"debug,omitempty"`
	Routes          []*Endpoint `gorm:"-" json:"routes,omitempty" yaml:"routes,omitempty" toml:"routes,omitempty"`
	Templates       Templates   `gorm:"-" json:"templates" yaml:"templates" toml:"templates"`
}

//the configuration file
type Templates map[string]map[string]*SelectorConfig

// type TemplateConfig SelectorConfig

type TemplateConfig2 struct {
	gorm.Model      `json:"-" yaml:"-" toml:"-"`
	sorting.Sorting `json:"-" yaml:"-" toml:"-"`
	EndpointID      uint                  `json:"-" yaml:"-" toml:"-"`
	Cache           bool                  `default:"true" etcd:"cache" json:"cache,omitempty" yaml:"cache,omitempty" toml:"cache,omitempty"`
	Collection      string                `json:"collection,omitempty" yaml:"collection,omitempty" toml:"collection,omitempty"`
	Description     string                `json:"description,omitempty" yaml:"description,omitempty" toml:"description,omitempty"`
	Required        bool                  `etcd:"required" default:"true" json:"required,omitempty" yaml:"required,omitempty" toml:"required,omitempty"`
	Items           string                `etcd:"items" json:"items,omitempty" yaml:"items,omitempty" toml:"items,omitempty"`
	Details         map[string]Extractors `gorm:"-" etcd:"details" json:"details,omitempty" yaml:"details,omitempty" toml:"details,omitempty"`
	StrictMode      bool                  `etcd:"strict_mode" default:"false" json:"strict_mode,omitempty" yaml:"strict_mode,omitempty" toml:"strict_mode,omitempty"`
	Debug           bool                  `etcd:"debug" default:"true" json:"debug,omitempty" yaml:"debug,omitempty" toml:"debug,omitempty"`
}

type CacheConfig struct {
	Disabled       bool          `default:"false" help:"Disable handler init" json:"disabled,omitempty" yaml:"disabled,omitempty" toml:"disabled,omitempty"`
	Engine         string        `json:"engine,omitempty" yaml:"engine,omitempty" toml:"engine,omitempty"`
	ListenAddr     string        `json:"listen_addr,omitempty" yaml:"listen_addr,omitempty" toml:"listen_addr,omitempty"`
	ExpirationTime time.Duration `json:"expiration_time,omitempty" yaml:"expiration_time,omitempty" toml:"expiration_time,omitempty"`
	PrefixPath     string        `json:"prefix_path,omitempty" yaml:"prefix_path,omitempty" toml:"prefix_path,omitempty"`
	Debug          bool          `default:"false" help:"Enable debug output for env vars processing" json:"debug,omitempty" yaml:"debug,omitempty" toml:"debug,omitempty"`
}

type EnvConfig struct {
	Disabled      bool                         `default:"false" help:"Disable handler init" json:"disabled,omitempty" yaml:"disabled,omitempty" toml:"disabled,omitempty"`
	Files         []string                     `json:"files,omitempty" yaml:"files,omitempty" toml:"files,omitempty"`
	VariablesList map[string]string            `json:"-" yaml:"-" toml:"-"`
	VariablesTree map[string]map[string]string `json:"-" yaml:"-" toml:"-"`
	Debug         bool                         `default:"false" help:"Enable debug output for env vars processing" json:"debug,omitempty" yaml:"debug,omitempty" toml:"debug,omitempty"`
}

// Endpoint represents a single remote endpoint. The performed query can be modified between each call by parameterising URL. See documentation.
type Endpoint struct {
	gorm.Model         `json:"-" yaml:"-" toml:"-"`
	sorting.Sorting    `json:"-" yaml:"-" toml:"-"`
	minFields          int                               `json:"-" yaml:"-" toml:"-"`
	count              string                            `gorm"-" json:"-" yaml:"-" toml:"-"`
	ready              bool                              `etcd:"-" json:"-" yaml:"-" toml:"-"`
	hash               string                            `etcd:"-" json:"-" yaml:"-" toml:"-"`
	Update             time.Time                         `json:"-" yaml:"-" toml:"-"`
	Disabled           bool                              `etcd:"disabled" json:"disabled,omitempty" yaml:"disabled,omitempty" toml:"disabled,omitempty"`
	Cache              bool                              `default:"true" etcd:"cache" json:"cache,omitempty" yaml:"cache,omitempty" toml:"cache,omitempty"`
	EtcdKey            string                            `etcd:"etcd_key" json:"etcd_key,omitempty" yaml:"etcd_key,omitempty" toml:"etcd_key,omitempty"`
	Connections        []Connection                      `json:"-" yaml:"-" toml:"-"`
	Source             string                            `etcd:"source" gorm:"-" json:"provider,omitempty" yaml:"provider,omitempty" toml:"provider,omitempty"`
	ProviderID         uint                              `json:"-" yaml:"-" toml:"-"`
	Provider           Provider                          `etcd:"provider" json:"provider_orm,omitempty" yaml:"provider_orm,omitempty" toml:"provider_orm,omitempty"`
	Comment            string                            `json:"comments,omitempty" yaml:"comments,omitempty" toml:"comments,omitempty"`
	Description        string                            `json:"description,omitempty" yaml:"description,omitempty" toml:"description,omitempty"`
	Groups             []*Group                          `etcd:"groups" json:"groups,omitempty" yaml:"groups,omitempty" toml:"groups,omitempty"`
	Route              string                            `etcd:"router" json:"route,omitempty" yaml:"route,omitempty" toml:"route,omitempty"`
	Method             string                            `gorm:"index" json:"method,omitempty" yaml:"method,omitempty" toml:"method,omitempty"`
	Template           string                            `gorm:"template" json:"template,omitempty" yaml:"template,omitempty" toml:"template,omitempty"`
	Domain             string                            `gorm:"-" json:"-" yaml:"-" toml:"-"`
	Host               string                            `gorm:"-" json:"-" yaml:"-" toml:"-"`
	Port               int                               `gorm:"-" json:"-" yaml:"-" toml:"-"`
	Pager              map[string]string                 `gorm:"-" json:"pager" yaml:"pager" toml:"pager"`
	Collection         bool                              `etcd:"collection" json:"collection,omitempty" yaml:"collection,omitempty" toml:"collection,omitempty"`
	Concurrency        int                               `default:"1" gorm:"concurrency" json:"concurrency,omitempty" yaml:"concurrency,omitempty" toml:"concurrency,omitempty"`
	BaseURL            string                            `etcd:"base_url" gorm:"base_url" json:"base_url,omitempty" yaml:"base_url,omitempty" toml:"base_url,omitempty"`
	PatternURL         string                            `etcd:"url" json:"url" yaml:"url" toml:"url"`
	Protocol           string                            `etcd:"protocol" gorm:"protocol" json:"protocol,omitempty" yaml:"protocol,omitempty" toml:"protocol,omitempty"`
	Transport          string                            `etcd:"transport" gorm:"transport" json:"transport,omitempty" yaml:"transport,omitempty" toml:"transport,omitempty"`
	Examples           map[string]map[string]string      `gorm:"-" json:"examples" yaml:"examples" toml:"examples"`
	Slug               string                            `etcd:"slug" json:"slug,omitempty" yaml:"slug,omitempty" toml:"slug,omitempty"`
	ExtractPaths       bool                              `etcd:"extract_paths" json:"extract_paths,omitempty" yaml:"extract_paths,omitempty" toml:"extract_paths,omitempty"`
	LeafPaths          []string                          `gorm:"-" json:"leaf_paths,omitempty" yaml:"leaf_paths,omitempty" toml:"leaf_paths,omitempty"`
	Body               string                            `gorm:"-" json:"body,omitempty" yaml:"body,omitempty" toml:"body,omitempty"`
	Selector           string                            `etcd:"selector" gorm:"index" default:"css" json:"selector,omitempty" yaml:"selector,omitempty" toml:"selector,omitempty"`
	HeadersIntercept   []string                          `etcd:"resp_headers_intercept" gorm:"-" json:"resp_headers_intercept,omitempty" yaml:"resp_headers_intercept,omitempty" toml:"resp_headers_intercept,omitempty"`
	Parameters         map[string]map[string]interface{} `etcd:"parameters" gorm:"-" json:"parameters,omitempty" yaml:"parameters,omitempty" toml:"parameters,omitempty"`
	HeadersJSON        map[string]string                 `etcd:"headers" gorm:"-" json:"headers,omitempty" yaml:"headers,omitempty" toml:"headers,omitempty"`
	BlocksJSON         map[string]*SelectorConfig        `etcd:"blocks" gorm:"-" json:"blocks,omitempty" yaml:"blocks,omitempty" toml:"blocks,omitempty"`
	Headers            []*HeaderConfig                   `json:"headers_orm,omitempty" yaml:"headers_orm,omitempty" toml:"headers_orm,omitempty"`
	Blocks             []*SelectorConfig                 `json:"blocks_orm,omitempty" yaml:"blocks_orm,omitempty" toml:"blocks_orm,omitempty"`
	EndpointProperties EndpointProperties                `etcd:"properties" sql:"type:text" json:"properties,omitempty" yaml:"properties,omitempty" toml:"properties,omitempty"`
	Extract            ExtractConfig                     `etcd:"extract" json:"extract,omitempty" yaml:"extract,omitempty" toml:"extract,omitempty"`
	Debug              bool                              `etcd:"debug" json:"debug,omitempty" yaml:"debug,omitempty" toml:"debug,omitempty"`
	StrictMode         bool                              `etcd:"strict_mode" json:"strict_mode,omitempty" yaml:"strict_mode,omitempty" toml:"strict_mode,omitempty"`
	Crawler            CrawlerConfig                     `etcd:"crawler" json:"crawler,omitempty" yaml:"crawler,omitempty" toml:"crawler,omitempty"`
}

type CrawlerConfig struct {
	// gorm.Model         `json:"-" yaml:"-" toml:"-"`
	// sorting.Sorting    `json:"-" yaml:"-" toml:"-"`
	// ExpirationDate *time.Time
	Disabled          bool                    `default:"false" gorm:"disabled" etcd:"disabled" json:"disabled,omitempty" yaml:"disabled,omitempty" toml:"disabled,omitempty"`
	Debug             bool                    `default:"false" etcd:"debug" json:"debug,omitempty" yaml:"debug,omitempty" toml:"debug,omitempty"`
	GetScreenshot     bool                    `default:"false" etcd:"screenshot" json:"screenshot,omitempty" yaml:"screenshot,omitempty" toml:"screenshot,omitempty"`
	IsContext         bool                    `default:"true" etcd:"is_context" json:"is_context,omitempty" yaml:"is_context,omitempty" toml:"is_context,omitempty"`
	IgnoreRobotsTxt   bool                    `default:"false" etcd:"ignore_robots_txt" json:"ignore_robots_txt,omitempty" yaml:"ignore_robots_txt,omitempty" toml:"ignore_robots_txt,omitempty"` // IgnoreRobotsTxt allows the Collector to ignore any restrictions set by
	AllowURLRevisit   bool                    `default:"false" etcd:"allow_url_revisit" json:"allow_url_revisit,omitempty" yaml:"allow_url_revisit,omitempty" toml:"allow_url_revisit,omitempty"` // AllowURLRevisit allows multiple downloads of the same URL
	AllowedDomains    []DomainConfig          `etcd:"allowed_domains" json:"allowed_domains,omitempty" yaml:"allowed_domains,omitempty" toml:"allowed_domains,omitempty"`                         // AllowedDomains is a domain whitelist.
	DisallowedDomains []DomainConfig          `etcd:"disallowed_domains" json:"disallowed_domains,omitempty" yaml:"disallowed_domains,omitempty" toml:"disallowed_domains,omitempty"`             // DisallowedDomains is a domain blacklist.
	URLFilters        []*regexp.Regexp        `etcd:"url_filters" json:"url_filters,omitempty" yaml:"url_filters,omitempty" toml:"url_filters,omitempty"`                                         // URLFilters is a list of regular expressions which restricts                                         // URLFilters is a list of regular expressions which restricts
	OnHTML            []CrawlerSelectorConfig `etcd:"selectors" json:"selectors,omitempty" yaml:"selectors,omitempty" toml:"selectors,omitempty"`
	Limits            LimitsConfig            `etcd:"limits" json:"limits,omitempty" yaml:"limits,omitempty" toml:"limits,omitempty"`
	MultiPart         MultiPartConfig         `default:"false" etcd:"is_multi_part" json:"is_multi_part,omitempty" yaml:"is_multi_part,omitempty" toml:"is_multi_part,omitempty"`
	CSV               CSVConfig               `etcd:"csv" json:"csv,omitempty" yaml:"csv,omitempty" toml:"csv,omitempty"`
	UserAgent         string                  `etcd:"user_agent" json:"user_agent,omitempty" yaml:"user_agent,omitempty" toml:"user_agent,omitempty"`                             // UserAgent is the User-Agent string used by HTTP requests
	CacheDir          string                  `default:"./shared/cache/colly" etcd:"cache_dir" json:"cache_dir,omitempty" yaml:"cache_dir,omitempty" toml:"cache_dir,omitempty"`  // CacheDir specifies a location where GET requests are cached as files.
	MaxBodySize       int                     `default:"10000" etcd:"max_body_size" json:"max_body_size,omitempty" yaml:"max_body_size,omitempty" toml:"max_body_size,omitempty"` // MaxBodySize is the limit of the retrieved response body in bytes.
	MaxDepth          int                     `default:"1" etcd:"max_depth" json:"max_depth,omitempty" yaml:"max_depth,omitempty" toml:"max_depth,omitempty"`                     // MaxDepth limits the recursion depth of visited URLs.
	// CacheLength     int             `etcd:"cache_length" json:"cache_length,omitempty" yaml:"cache_length,omitempty" toml:"cache_length,omitempty"`	// courses := make([]Course, 0, 200)
}

type MultiPartConfig struct {
	// gorm.Model         `json:"-" yaml:"-" toml:"-"`
	// sorting.Sorting    `json:"-" yaml:"-" toml:"-"`
	Disabled bool              `default:"false" gorm:"disabled" etcd:"disabled" json:"disabled,omitempty" yaml:"disabled,omitempty" toml:"disabled,omitempty"`
	Params   map[string]string `etcd:"params" json:"params,omitempty" yaml:"params,omitempty" toml:"params,omitempty"`
	Files    map[string]string `etcd:"files" json:"files,omitempty" yaml:"files,omitempty" toml:"files,omitempty"`
	Debug    bool              `default:"false" etcd:"debug" json:"debug,omitempty" yaml:"debug,omitempty" toml:"debug,omitempty"`
}

type CrawlerSelectorConfig struct {
	// gorm.Model         `json:"-" yaml:"-" toml:"-"`
	// sorting.Sorting    `json:"-" yaml:"-" toml:"-"`
	Disabled bool   `default:"false" gorm:"disabled" etcd:"disabled" json:"disabled,omitempty" yaml:"disabled,omitempty" toml:"disabled,omitempty"`
	Macro    string `etcd:"macro" json:"macro,omitempty" yaml:"macro,omitempty" toml:"macro,omitempty"`
	Selector string `etcd:"selector" json:"selector,omitempty" yaml:"selector,omitempty" toml:"selector,omitempty"`
	Debug    bool   `default:"false" etcd:"debug" json:"debug,omitempty" yaml:"debug,omitempty" toml:"debug,omitempty"`
}

type LimitsConfig struct {
	// gorm.Model         `json:"-" yaml:"-" toml:"-"`
	// sorting.Sorting    `json:"-" yaml:"-" toml:"-"`
	Disabled     bool          `default:"false" gorm:"disabled" etcd:"disabled" json:"disabled,omitempty" yaml:"disabled,omitempty" toml:"disabled,omitempty"`
	DomainGlob   string        `etcd:"domain_glob" json:"domain_glob,omitempty" yaml:"domain_glob,omitempty" toml:"domain_glob,omitempty"`             // DomainRegexp is a glob pattern to match against domains
	Parallelism  int           `default:"5" etcd:"parallelism" json:"parallelism,omitempty" yaml:"parallelism,omitempty" toml:"parallelism,omitempty"` // Parallelism is the number of the maximum allowed concurrent requests of the matching domains
	DomainRegexp string        `etcd:"domain_regexp" json:"domain_regexp,omitempty" yaml:"domain_regexp,omitempty" toml:"domain_regexp,omitempty"`     // DomainRegexp is a regular expression to match against domains
	Delay        time.Duration `etcd:"delay" json:"delay,omitempty" yaml:"delay,omitempty" toml:"delay,omitempty"`                                     // Delay is the duration to wait before creating a new request to the matching domains
	Debug        bool          `default:"false" etcd:"debug" json:"debug,omitempty" yaml:"debug,omitempty" toml:"debug,omitempty"`
}

type DomainConfig struct {
	// gorm.Model         `json:"-" yaml:"-" toml:"-"`
	// sorting.Sorting    `json:"-" yaml:"-" toml:"-"`
	Disabled bool `default:"false" gorm:"disabled" etcd:"disabled" json:"disabled,omitempty" yaml:"disabled,omitempty" toml:"disabled,omitempty"`
	MaxDepth int  `default:"1" etcd:"max_depth" json:"max_depth,omitempty" yaml:"max_depth,omitempty" toml:"max_depth,omitempty"`
	Debug    bool `default:"false" etcd:"debug" json:"debug,omitempty" yaml:"debug,omitempty" toml:"debug,omitempty"`
}

type CSVConfig struct {
	// gorm.Model         `json:"-" yaml:"-" toml:"-"`
	// sorting.Sorting    `json:"-" yaml:"-" toml:"-"`
	Disabled   bool   `gorm:"disabled" etcd:"disabled" json:"disabled,omitempty" yaml:"disabled,omitempty" toml:"disabled,omitempty"`
	PrefixPath string `etcd:"prefix_path" json:"prefix_path,omitempty" yaml:"prefix_path,omitempty" toml:"prefix_path,omitempty"`
	Separator  string `etcd:"separator" json:"separator,omitempty" yaml:"separator,omitempty" toml:"separator,omitempty"`
	Headers    string `etcd:"headers" json:"headers,omitempty" yaml:"headers,omitempty" toml:"headers,omitempty"`
	Debug      bool   `etcd:"debug" json:"debug,omitempty" yaml:"debug,omitempty" toml:"debug,omitempty"`
}

// Mail is the container of a single e-mail
type Mail struct {
	// gorm.Model         `json:"-" yaml:"-" toml:"-"`
	// sorting.Sorting    `json:"-" yaml:"-" toml:"-"`
	Disabled bool   `gorm:"disabled" etcd:"disabled" json:"disabled,omitempty" yaml:"disabled,omitempty" toml:"disabled,omitempty"`
	Title    string `json:"title,omitempty" yaml:"title,omitempty" toml:"title,omitempty"`
	Link     string `json:"link,omitempty" yaml:"link,omitempty" toml:"link,omitempty"`
	Author   string `json:"author,omitempty" yaml:"author,omitempty" toml:"author,omitempty"`
	Date     string `json:"date,omitempty" yaml:"date,omitempty" toml:"date,omitempty"`
	Message  string `json:"message,omitempty" yaml:"message,omitempty" toml:"message,omitempty"`
	Debug    bool   `etcd:"debug" json:"debug,omitempty" yaml:"debug,omitempty" toml:"debug,omitempty"`
}

// Endpoint represents a single remote endpoint. The performed query can be modified between each call by parameterising URL. See documentation.
type Template struct {
	gorm.Model      `json:"-" yaml:"-" toml:"-"`
	sorting.Sorting `json:"-" yaml:"-" toml:"-"`
	UUID            string                       `gorm:"uuid" etcd:"uuid" json:"uuid,omitempty" yaml:"uuid,omitempty" toml:"uuid,omitempty"`
	Disabled        bool                         `gorm:"disabled" etcd:"disabled" json:"disabled,omitempty" yaml:"disabled,omitempty" toml:"disabled,omitempty"`
	Cache           bool                         `gorm:"cache" default:"true" etcd:"cache" json:"cache,omitempty" yaml:"cache,omitempty" toml:"cache,omitempty"`
	Method          string                       `gorm:"index" json:"method,omitempty" yaml:"method,omitempty" toml:"method,omitempty"`
	Pager           map[string]string            `gorm:"-" json:"pager" yaml:"pager" toml:"pager"`
	Collection      bool                         `etcd:"collection" json:"collection,omitempty" yaml:"collection,omitempty" toml:"collection,omitempty"`
	Concurrency     int                          `gorm:"concurrency" default:"1" json:"concurrency,omitempty" yaml:"concurrency,omitempty" toml:"concurrency,omitempty"`
	Selector        string                       `gorm:"index" etcd:"selector" default:"css" json:"selector,omitempty" yaml:"selector,omitempty" toml:"selector,omitempty"`
	Catch           Intercept                    `gorm:"-" etcd:"resp_headers_intercept" json:"resp_headers_intercept,omitempty" yaml:"resp_headers_intercept,omitempty" toml:"resp_headers_intercept,omitempty"`
	Blocks          map[string]map[string]string `gorm:"-" etcd:"blocks" json:"blocks,omitempty" yaml:"blocks,omitempty" toml:"blocks,omitempty"`
	Debug           bool                         `gorm:"debug" etcd:"debug" json:"debug,omitempty" yaml:"debug,omitempty" toml:"debug,omitempty"`
	StrictMode      bool                         `gorm:"strict_mode" etcd:"strict_mode" json:"strict_mode,omitempty" yaml:"strict_mode,omitempty" toml:"strict_mode,omitempty"`
}

type Intercept struct {
	gorm.Model      `json:"-" yaml:"-" toml:"-"`
	sorting.Sorting `json:"-" yaml:"-" toml:"-"`
	hash            string     `gorm:"hash" etcd:"hash" json:"hash,omitempty" yaml:"hash,omitempty" toml:"hash,omitempty"`
	Disabled        bool       `gorm:"disabled" etcd:"disabled" json:"disabled,omitempty" yaml:"disabled,omitempty" toml:"disabled,omitempty"`
	Headers         []*Header  `json:"headers,omitempty" yaml:"headers,omitempty" toml:"headers,omitempty"`
	Body            []*Pattern `json:"patterns,omitempty" yaml:"patterns,omitempty" toml:"patterns,omitempty"`
	Debug           bool       `gorm:"debug" etcd:"debug" json:"debug,omitempty" yaml:"debug,omitempty" toml:"debug,omitempty"`
	StrictMode      bool       `gorm:"strict_mode" etcd:"strict_mode" json:"strict_mode,omitempty" yaml:"strict_mode,omitempty" toml:"strict_mode,omitempty"`
}

type Pattern struct {
	gorm.Model      `json:"-" yaml:"-" toml:"-"`
	sorting.Sorting `json:"-" yaml:"-" toml:"-"`
	hash            string `gorm:"hash" etcd:"hash" json:"hash,omitempty" yaml:"hash,omitempty" toml:"hash,omitempty"`
	Disabled        bool   `gorm:"disabled" etcd:"disabled" json:"disabled,omitempty" yaml:"disabled,omitempty" toml:"disabled,omitempty"`
	Pattern         string `gorm:"pattern" etcd:"pattern" json:"pattern,omitempty" yaml:"pattern,omitempty" toml:"pattern,omitempty"`
	Debug           bool   `gorm:"debug" etcd:"debug" json:"debug,omitempty" yaml:"debug,omitempty" toml:"debug,omitempty"`
	StrictMode      bool   `gorm:"strict_mode" etcd:"strict_mode" json:"strict_mode,omitempty" yaml:"strict_mode,omitempty" toml:"strict_mode,omitempty"`
}

type HTMLMeta struct {
	Title         string
	Description   string
	OGTitle       string
	OGDescription string
	OGImage       string
	OGAuthor      string
	OGPublisher   string
	OGSiteName    string
}

type MetaData struct {
	Title       string  `meta:"og:title"`
	Description string  `meta:"og:description,description"`
	Type        string  `meta:"og:type"`
	URL         url.URL `meta:"og:url"`
	VideoWidth  int64   `meta:"og:video:width"`
	VideoHeight int64   `meta:"og:video:height"`
}

type Header struct {
	gorm.Model      `json:"-" yaml:"-" toml:"-"`
	sorting.Sorting `json:"-" yaml:"-" toml:"-"`
	hash            string `gorm:"hash" etcd:"hash" json:"hash,omitempty" yaml:"hash,omitempty" toml:"hash,omitempty"`
	Disabled        bool   `gorm:"disabled" etcd:"disabled" json:"disabled,omitempty" yaml:"disabled,omitempty" toml:"disabled,omitempty"`
	Key             string `gorm:"key" etcd:"key" json:"key,omitempty" yaml:"key,omitempty" toml:"key,omitempty"`
	Value           string `gorm:"value" etcd:"value" json:"value,omitempty" yaml:"value,omitempty" toml:"value,omitempty"`
	Debug           bool   `gorm:"debug" etcd:"debug" json:"debug,omitempty" yaml:"debug,omitempty" toml:"debug,omitempty"`
	StrictMode      bool   `gorm:"strict_mode" etcd:"strict_mode" json:"strict_mode,omitempty" yaml:"strict_mode,omitempty" toml:"strict_mode,omitempty"`
}

type Screenshot struct {
	gorm.Model   `json:"-" yaml:"-" toml:"-"`
	Title        string                            `etcd:"title" json:"title,omitempty" yaml:"title,omitempty" toml:"title,omitempty"`
	EndpointID   uint                              `json:"-" yaml:"-" toml:"-"`
	SelectedType string                            `etcd:"selected_type" json:"selected_type,omitempty" yaml:"selected_type,omitempty" toml:"selected_type,omitempty"`
	File         media_library.MediaLibraryStorage `sql:"size:4294967295;" media_library:"url:/system/{{class}}/{{primary_key}}/{{column}}.{{extension}}" json:"-" yaml:"-" toml:"-"`
	// Category     Category
	// CategoryID   uint
}

func (screenshot Screenshot) Validate(db *gorm.DB) {
	if strings.TrimSpace(screenshot.Title) == "" {
		db.AddError(validations.NewError(screenshot, "Title", "Title can not be empty"))
	}
}

func (screenshot *Screenshot) SetSelectedType(typ string) {
	screenshot.SelectedType = typ
}

func (screenshot *Screenshot) GetSelectedType() string {
	return screenshot.SelectedType
}

func (screenshot *Screenshot) ScanMediaOptions(mediaOption media_library.MediaOption) error {
	if bytes, err := json.Marshal(mediaOption); err == nil {
		return screenshot.File.Scan(bytes)
	} else {
		return err
	}
}

func (screenshot *Screenshot) GetMediaOption() (mediaOption media_library.MediaOption) {
	mediaOption.Video = screenshot.File.Video
	mediaOption.FileName = screenshot.File.FileName
	mediaOption.URL = screenshot.File.URL()
	mediaOption.OriginalURL = screenshot.File.URL("original")
	mediaOption.CropOptions = screenshot.File.CropOptions
	mediaOption.Sizes = screenshot.File.GetSizes()
	mediaOption.Description = screenshot.File.Description
	return
}

/*
type ScreenShotVariationImageStorage struct{ oss.OSS }

func (colorVariation ScreenShot) MainImageURL() string {
	if len(colorVariation.Images.Files) > 0 {
		return colorVariation.Images.URL()
	}
	return "/images/default_product.png"
}

func (ScreenShotVariationImageStorage) GetSizes() map[string]*media.Size {
	return map[string]*media.Size{
		"small":  {Width: 320, Height: 320},
		"middle": {Width: 640, Height: 640},
		"big":    {Width: 1280, Height: 1280},
	}
}
*/

type Queries struct {
	gorm.Model          `json:"-" yaml:"-" toml:"-"`
	sorting.SortingDESC `json:"-" yaml:"-" toml:"-"`
	Keywords            []Query `etcd:"keywords" json:"keywords,omitempty" yaml:"keywords,omitempty" toml:"keywords,omitempty"`
}

type Query struct {
	gorm.Model `json:"-" yaml:"-" toml:"-"`
	InputQuery string `etcd:"input_query" json:"input_query,omitempty" yaml:"input_query,omitempty" toml:"input_query,omitempty"`
	Slug       string `etcd:"slug" json:"slug,omitempty" yaml:"slug,omitempty" toml:"slug,omitempty"`
	MD5        string `etcd:"md5" json:"md5,omitempty" yaml:"md5,omitempty" toml:"md5,omitempty"`
	SHA1       string `etcd:"sha1" json:"sha1,omitempty" yaml:"sha1,omitempty" toml:"sha1,omitempty"`
	UUID       string `etcd:"uuid" json:"uuid,omitempty" yaml:"uuid,omitempty" toml:"uuid,omitempty"`
	Blocked    bool   `etcd:"blocked" json:"blocked,omitempty" yaml:"blocked,omitempty" toml:"blocked,omitempty"`
}

type EndpointProperties []EndpointProperty // `etcd:"properties" json:"properties" yaml:"properties" toml:"properties"`

type EndpointProperty struct {
	Name  string `etcd:"name" json:"name" yaml:"name" toml:"name"`
	Value string `etcd:"value" json:"value" yaml:"value" toml:"value"`
}

func (endpointProperties *EndpointProperties) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, endpointProperties)
	case string:
		if v != "" {
			return endpointProperties.Scan([]byte(v))
		}
	default:
		return errors.New("not supported")
	}
	return nil
}

func (endpointProperties EndpointProperties) Value() (driver.Value, error) {
	if len(endpointProperties) == 0 {
		return nil, nil
	}
	return json.Marshal(endpointProperties)
}

// SelectorConfig represents a content selection rule for a single URL Pattern.
type SelectorConfig struct {
	gorm.Model      `json:"-" yaml:"-" toml:"-"`
	sorting.Sorting `json:"-" yaml:"-" toml:"-"`
	EndpointID      uint                  `json:"-" yaml:"-" toml:"-"`
	Cache           bool                  `default:"true" etcd:"cache" json:"cache,omitempty" yaml:"cache,omitempty" toml:"cache,omitempty"`
	EtcdKey         string                `etcd:"etcd_key" json:"etcd_key,omitempty" yaml:"etcd_key,omitempty" toml:"etcd_key,omitempty"`
	Collection      string                `json:"collection,omitempty" yaml:"collection,omitempty" toml:"collection,omitempty"`
	Description     string                `json:"description,omitempty" yaml:"description,omitempty" toml:"description,omitempty"`
	Required        bool                  `etcd:"required" default:"true" json:"required,omitempty" yaml:"required,omitempty" toml:"required,omitempty"`
	Items           string                `etcd:"items" json:"items,omitempty" yaml:"items,omitempty" toml:"items,omitempty"`
	Details         map[string]Extractors `gorm:"-" etcd:"details" json:"details,omitempty" yaml:"details,omitempty" toml:"details,omitempty"`
	Paths           map[string]string     `gorm:"-" etcd:"paths" json:"paths,omitempty" yaml:"paths,omitempty" toml:"paths,omitempty"`
	Matchers        []*MatcherConfig      `json:"matchers,omitempty" yaml:"matchers,omitempty" toml:"matchers,omitempty"`
	StrictMode      bool                  `etcd:"strict_mode" default:"false" json:"strict_mode,omitempty" yaml:"strict_mode,omitempty" toml:"strict_mode,omitempty"`
	Debug           bool                  `etcd:"debug" default:"true" json:"debug,omitempty" yaml:"debug,omitempty" toml:"debug,omitempty"`
}

// Extractor represents a pair of css selector and extracted node.
type Extractor struct {
	val string      `etcd:"value" json:"value" yaml:"value" toml:"value"`
	fn  extractorFn `gorm:"-" json:"-" yaml:"-" toml:"-"`
}

// Extractor represents a pair of css selector and extracted node.
type MatcherConfig struct {
	gorm.Model       `json:"-" yaml:"-" toml:"-"`
	sorting.Sorting  `json:"-" yaml:"-" toml:"-"`
	SelectorConfigID uint      `json:"-" yaml:"-" toml:"-"`
	Target           string    `etcd:"target" json:"target,omitempty" yaml:"target,omitempty" toml:"target,omitempty"`
	Selects          []Matcher `etcd:"selects" json:"selects,omitempty" yaml:"selects,omitempty" toml:"selects,omitempty"`
	EtcdKey          string    `etcd:"etcd_key" json:"etcd_key,omitempty" yaml:"etcd_key,omitempty" toml:"etcd_key,omitempty"`
}

//type Matchers {[]Matcher
type Matcher struct {
	gorm.Model      `json:"-" yaml:"-" toml:"-"`
	MatcherConfigID uint   `json:"-" yaml:"-" toml:"-"`
	Expression      string `etcd:"expr" json:"expr,omitempty" yaml:"expr,omitempty" toml:"expr,omitempty"`
}

type SelectorType struct {
	gorm.Model `json:"-" yaml:"-" toml:"-"`
	Name       string `etcd:"name" json:"name,omitempty" yaml:"name,omitempty" toml:"name,omitempty"`
	Engine     string `etcd:"engine" json:"engine,omitempty" yaml:"engine,omitempty" toml:"engine,omitempty"`
}

type TargetConfig struct {
	gorm.Model `json:"-" yaml:"-" toml:"-"`
	// EndpointID uint `json:"-" yaml:"-" toml:"-"`
	Name string `etcd:"name" json:"name,omitempty" yaml:"name,omitempty" toml:"name,omitempty"`
}

type HeaderConfig struct {
	gorm.Model `json:"-" yaml:"-" toml:"-"`
	EndpointID uint   `json:"-" yaml:"-" toml:"-"`
	Key        string `etcd:"key" json:"key,omitempty" yaml:"key,omitempty" toml:"key,omitempty"`
	Value      string `etcd:"value" json:"value,omitempty" yaml:"value,omitempty" toml:"value,omitempty"`
}

type BlocksConfig struct {
	gorm.Model `json:"-" yaml:"-" toml:"-"`
	Key        string         `etcd:"key" json:"key,omitempty" yaml:"key,omitempty" toml:"key,omitempty"`
	Value      SelectorConfig `etcd:"value" json:"value,omitempty" yaml:"value,omitempty" toml:"value,omitempty"`
}

type ExtractorsConfig struct {
	gorm.Model `json:"-" yaml:"-" toml:"-"`
	Key        string     `etcd:"key" json:"key,omitempty" yaml:"key,omitempty" toml:"key,omitempty"`
	Value      Extractors `etcd:"value" json:"value,omitempty" yaml:"value,omitempty" toml:"value,omitempty"`
}

// ExtractConfig represents a single sub-extraction rules url content configuration.
type ExtractConfig struct {
	gorm.Model `json:"-" yaml:"-" toml:"-"`
	Debug      bool `default:"true" json:"debug,omitempty" yaml:"debug,omitempty" toml:"debug,omitempty"`
	Links      bool `default:"true" json:"links,omitempty" yaml:"links,omitempty" toml:"links,omitempty"`
	Meta       bool `default:"true" json:"meta,omitempty" yaml:"meta,omitempty" toml:"meta,omitempty"`
	Opengraph  bool `default:"true" json:"opengraph,omitempty" yaml:"opengraph,omitempty" toml:"opengraph,omitempty"`
}

// OPENAPI SCRAPER ///////////////////////////////////////////////////////////////
type OpenAPIConfig struct {
	gorm.Model `json:"-" yaml:"-" toml:"-"`
	Name       string                `etcd:"name" json:"name,omitempty" yaml:"name,omitempty" toml:"name,omitempty"`
	Provider   Provider              `etcd:"provider" json:"provider,omitempty" yaml:"provider,omitempty" toml:"provider,omitempty"`
	Specs      []*OpenAPISpecsConfig `etcd:"specs" json:"specs,omitempty" yaml:"specs,omitempty" toml:"specs,omitempty"`
}

type OpenAPISpecsConfig struct {
	gorm.Model `json:"-" yaml:"-" toml:"-"`
	Slug       string `etcd:"slug" json:"slug,omitempty" yaml:"slug,omitempty" toml:"slug,omitempty"`
	Version    string `etcd:"version" json:"version,omitempty" yaml:"version,omitempty" toml:"version,omitempty"`
}

// REQUESTS API WEBMOCKS ///////////////////////////////////////////////////////////////
type Connection struct {
	gorm.Model `json:"-" yaml:"-" toml:"-"`
	// ID         uint     `gorm:"primary_key;AUTO_INCREMENT" json:"-" yaml:"-" toml:"-"`
	EndpointID uint     `json:"-" yaml:"-" toml:"-"`
	URL        string   `json:"url" yaml:"url" toml:"url"`
	Request    Request  `json:"request" yaml:"request" toml:"request"`
	Response   Response `json:"response" yaml:"response" toml:"response"`
	Provider   Provider `json:"provider" yaml:"provider" toml:"provider"`
	RecordedAt string   `json:"recorded_at" yaml:"recorded_at" toml:"recorded_at"`
}

type Request struct {
	gorm.Model `json:"-" yaml:"-" toml:"-"`
	// ID           uint   `gorm:"primary_key;AUTO_INCREMENT" json:"-" yaml:"-" toml:"-"`
	ConnectionID uint   `json:"-" yaml:"-" toml:"-"`
	Header       string `json:"header" yaml:"header" toml:"header"`
	Body         string `json:"body" yaml:"body" toml:"body"`
	Method       string `json:"method" yaml:"method" toml:"method"`
	URL          string `json:"url" yaml:"url" toml:"url"`
}

type Response struct {
	gorm.Model `json:"-" yaml:"-" toml:"-"`
	// ID           uint   `gorm:"primary_key;AUTO_INCREMENT" json:"-" yaml:"-" toml:"-"`
	ConnectionID uint   `json:"-" yaml:"-" toml:"-"`
	Status       string `json:"status" yaml:"status" toml:"status"`
	Header       string `json:"header" yaml:"header" toml:"header"`
	Body         string `json:"body" yaml:"body" toml:"body"`
}
