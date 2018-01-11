package scraper

import (
	"context"
	"crypto/md5"
	//"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Obaied/RAKE.Go"
	"github.com/PuerkitoBio/goquery"
	"github.com/advancedlogic/GoOse"
	"github.com/alexflint/go-restructure"
	"github.com/antchfx/xquery/html"
	"github.com/archivers-space/warc"
	"github.com/benmanns/goworker"
	"github.com/cnf/structhash"
	"github.com/gebv/typed"
	"github.com/go-resty/resty"
	"github.com/gocolly/colly"
	"github.com/iancoleman/strcase"
	"github.com/jeevatkm/go-model"
	"github.com/jeffail/tunny"
	"github.com/k0kubun/pp"
	"github.com/kamildrazkiewicz/go-flow"
	"github.com/karlseguin/cmap"
	"github.com/leebenson/conform"
	"github.com/mgbaozi/gomerge"
	"github.com/mmcdole/gofeed"
	"github.com/oleiade/reflections"
	"github.com/parnurzeal/gorequest"
	"github.com/roscopecoltran/mxj"
	comap "github.com/streamrail/concurrent-map"
	"github.com/trustmaster/goflow"
	"github.com/tsak/concurrent-csv-writer"
	// "github.com/urandom/text-summary/summarize"
	"github.com/vbauerster/mpb"
	"github.com/vbauerster/mpb/decor"
	"golang.org/x/net/html"
	// tg "github.com/galeone/tfgo"
	// tf "github.com/tensorflow/tensorflow/tensorflow/go"
	// "github.com/ekzhu/lshensemble"
	// "github.com/gpahal/go-meta"
	// "github.com/jonlaing/htmlmeta"
	// m "github.com/keighl/metabolize"
	// "github.com/frozzare/go-ogp"
	// "github.com/as27/fieldcopy"
	// "github.com/foize/go.fifo"
	// "github.com/hydrogen18/stoppableListener"
	// cache "github.com/patrickmn/go-cache"
	// "github.com/patrickmn/go-cache"
	// "github.com/LibertyLocked/cachestore"
	// "github.com/fatih/color"
	// "golang.org/x/net/proxy" // https://github.com/prsolucoes/go-tor-crawler/blob/master/main.go
	// jsonql "github.com/ugurozgen/json-transformer"
	// cregex "github.com/mingrammer/commonregex"
	// rx "github.com/yargevad/regexpx"
	// "github.com/naivesound/expr-go"
	// "github.com/OAGr/rulebook"
	// "github.com/Maldris/mathparse"
	// "github.com/goanywhere/regex"
	// "github.com/xyproto/lookup"
	// "github.com/elgs/jsonql"
	// "github.com/h12w/dfa"
	// "github.com/nytlabs/gojee"
	// "github.com/cevaris/ordered_map"
	// "github.com/iancoleman/orderedmap"
	// "golang.org/x/net/context/ctxhttp"
	// "github.com/datatogether/pdf"
	// "github.com/datatogether/linked_data"
	// "github.com/datatogether/linked_data/dcat"
	// "github.com/datatogether/linked_data/pod"
	// "github.com/datatogether/linked_data/sciencebase"
	// "github.com/datatogether/linked_data/jsonld"
	// "github.com/datatogether/linked_data/xmp"
	// "github.com/ctessum/requestcache"
	// "github.com/otiai10/cachely"
	// "github.com/buger/jsonparser"
	// "github.com/go-aah/aah"
	// "github.com/creack/spider"
	// "github.com/whyrusleeping/json-filter"
	// "github.com/wolfeidau/unflatten"
	// "github.com/jzaikovs/t"
	// "github.com/linkosmos/urlutils"
	// "github.com/microcosm-cc/bluemonday"
	// "github.com/kennygrant/sanitize"
	// "github.com/slotix/slugifyurl"
	// "github.com/antchfx/xpath"
	// "github.com/advancedlogic/GoOse"
	// "github.com/ynqa/word-embedding/builder"
	// "github.com/ynqa/word-embedding/config"
	// "github.com/ynqa/word-embedding/validate"
)

/*
	Refs:
	- github.com/ahmetb/go-linq
	- interfacer -for \"github.com/roscopecoltran/scraper/scraper\".Config -as mock.Scraper
	- github.com/slotix/dataflowkit
	- github.com/slotix/pageres-go-wrapper
	- github.com/fern4lvarez/go-metainspector
	- github.com/gpahal/go-meta
	- https://github.com/scrapinghub/mdr
	- https://github.com/scrapinghub/aile/blob/master/demo2.py
	- https://github.com/datatogether/sentry
	- https://github.com/sourcegraph/webloop
	- https://github.com/107192468/sp/blob/master/src/readhtml/readhtml.go
	- https://github.com/nikolay-turpitko/structor
	- https://github.com/dreampuf/paw/tree/master/src/web
	- https://github.com/rakanalh/grawler/blob/master/processors/text.go
	- https://github.com/rakanalh/grawler/blob/master/extractor/xpath.go
	- https://github.com/rakanalh/grawler/blob/master/extractor/css.go
	- https://github.com/ErosZy/labour/blob/master/parser/pageItemXpathParser.go
	- https://github.com/ErosZy/labour
	- https://github.com/cugbliwei/crawler/blob/master/extractor/selector.go
	- https://github.com/xlvector/higgs/blob/master/extractor/selector.go
	- github.com/tchssk/link
	- https://github.com/peterhellberg/link
	- https://github.com/jpillora/scraper/commit/0b5e5ce320ffaaaf86fb3ba9cc49458df3406a86
	- https://github.com/KKRainbow/segmentation-server/blob/master/main.go
	- https://github.com/mhausenblas/github-api-fetcher/blob/master/main.go
	- https://github.com/hoop33/limo/blob/master/service/github.go#L39
	- https://github.com/creack/spider/blob/master/example_test.go
	- https://github.com/suwhs/go-goquery-utils/tree/master/pipes
	- https://github.com/andrewstuart/goq

	- https://golanglibs.com/search?page=4&q=expression
	- https://github.com/elves/elvish
	- https://github.com/google/zoekt
	- https://github.com/pointlander/peg
	- https://github.com/pointlander/peg/blob/master/grammars/c/c.peg
	- https://github.com/icochico/regexbench

	- https://github.com/hermanschaaf/go-string-concat-benchmarks
	- https://github.com/elgs/jsonql
	- https://github.com/alexflint/go-restructure/tree/master/samples
	- https://github.com/prsolucoes/go-tor-crawler/blob/master/main.go

	- https://github.com/mergermarket/news-aggro/blob/26d6834805449ed74684c3facac0813e339a8d8d/main.go#L21
	- github.com/mateuszdyminski/bloom-filter
	- https://github.com/Flaque/wikiracer
*/

const defaultUA = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/53.0.2785.143 Safari/537.36"

var (
	cacheDuration = 3600 * time.Second
	userAgents    = []string{
		"Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
		"GalaxyBot/1.0 (http://www.galaxy.com/galaxybot.html)",
		"Googlebot-Image/1.0",
	}
)

// https://github.com/mmadfox/scraper/blob/master/visits.go
type Visiter interface {
	Visit(string) bool
	ResetVisit(string) error
	Drop() error
	Close() error
}

type memoryVisits struct {
	m     comap.ConcurrentMap
	mutex *sync.Mutex
}

func (v *memoryVisits) Visit(u string) bool {
	return !v.m.SetIfAbsent(u, 1)
}

func (v *memoryVisits) ResetVisit(u string) error {
	v.m.Remove(u)
	return nil
}

func (v *memoryVisits) Drop() error {
	v.mutex.Lock()
	defer v.mutex.Unlock()

	v.m = comap.New()
	return nil
}

func (v *memoryVisits) Close() error {
	return nil
}

func NewMemoryVisits() Visiter {
	return &memoryVisits{
		m:     comap.New(),
		mutex: &sync.Mutex{},
	}
}

func RandomUserAgent() string {
	return userAgents[rand.Intn(len(userAgents))]
}

func myFunc(queue string, args ...interface{}) error {
	fmt.Printf("From %s, %v\n", queue, args)
	return nil
}

/*
func newMyFunc(uri string) func(queue string, args ...interface{}) error {
	foo := NewFoo(uri)
	return func(queue string, args ...interface{}) error {
		foo.Bar(args)
		return nil
	}
}
*/

func testGoose(link string) {
	g := goose.New()
	if link == "" {
		link = "http://edition.cnn.com/2012/07/08/opinion/banzi-ted-open-source/index.html"
	}
	article, _ := g.ExtractFromURL(link)
	println("link: ", link)
	println("title: ", article.Title)
	println("description: ", article.MetaDescription)
	println("keywords: ", article.MetaKeywords)
	println("content: ", article.CleanedText)
	println("url: ", article.FinalURL)
	println("top image: ", article.TopImage)
	println("============================================================")
}

func init() {
	settings := goworker.WorkerSettings{
		URI:            "redis://localhost:6379/",
		Connections:    100,
		Queues:         []string{"myqueue", "delimited", "queues"},
		UseNumber:      true,
		ExitOnComplete: false,
		Concurrency:    2,
		Namespace:      "resque:",
		Interval:       5.0,
	}
	goworker.SetSettings(settings)
	goworker.Register("MyClass", myFunc)
	// goworker.Register("MyClass", newMyFunc("http://www.example.com/"))
}

func testGoWorker() {

	goworker.Enqueue(&goworker.Job{
		Queue: "myqueue",
		Payload: goworker.Payload{
			Class: "MyClass",
			Args:  []interface{}{"hi", "there"},
		},
	})

	if err := goworker.Work(); err != nil {
		fmt.Println("Error:", err)
	}
}

func testGoFlow() {
	f1 := func(r map[string]interface{}) (interface{}, error) {
		fmt.Println("function1 started")
		time.Sleep(time.Millisecond * 1000)
		return 1, nil
	}

	f2 := func(r map[string]interface{}) (interface{}, error) {
		time.Sleep(time.Millisecond * 1000)
		fmt.Println("function2 started", r["f1"])
		return "some results", nil // errors.New("Some error")
	}

	f3 := func(r map[string]interface{}) (interface{}, error) {
		fmt.Println("function3 started", r["f1"])
		return nil, nil
	}

	f4 := func(r map[string]interface{}) (interface{}, error) {
		fmt.Println("function4 started", r)
		return nil, nil
	}

	res, err := goflow.New().
		Add("f1", nil, f1).
		Add("f2", []string{"f1"}, f2).
		Add("f3", []string{"f1"}, f3).
		Add("f4", []string{"f2", "f3"}, f4).
		Do()

	fmt.Println(res, err)
}

var quaternionRegexp = restructure.MustCompile(QuotedQuaternion{}, restructure.Options{})

type EmailAddress struct {
	_    struct{} `^`
	User string   `\w+`
	_    struct{} `@`
	Host string   `[^@]+`
	_    struct{} `$`
}

type Hostname struct {
	Domain string   `\w+`
	_      struct{} `\.`
	TLD    string   `\w+`
}

type EmailAddress2 struct {
	_    struct{} `^`
	User string   `[a-zA-Z0-9._%+-]+`
	_    struct{} `@`
	Host *Hostname
	_    struct{} `$`
}

// Matches "123", "1.23", "1.23e-4", "-12.3E+5", ".123"
type Float struct {
	Sign     *Sign     `?` // sign is optional
	Whole    string    `[0-9]*`
	Period   struct{}  `\.?`
	Frac     string    `[0-9]+`
	Exponent *Exponent `?` // exponent is optional
}

// Matches "e+4", "E6", "e-03"
type Exponent struct {
	_    struct{} `[eE]`
	Sign *Sign    `?` // sign is optional
	Num  string   `[0-9]+`
}

// Matches "+" or "-"
type Sign struct {
	Ch string `[+-]`
}

type RealPart struct {
	Sign string `regexp:"[+-]?"`
	Real string `regexp:"[0-9]+"`
}

type SignedInt struct {
	Sign string `regexp:"[+-]"`
	Real string `regexp:"[0-9]+"`
}

type IPart struct {
	Magnitude SignedInt
	_         struct{} `regexp:"i"`
}

type JPart struct {
	Magnitude SignedInt
	_         struct{} `regexp:"j"`
}

type KPart struct {
	Magnitude SignedInt
	_         struct{} `regexp:"k"`
}

// matches "1+2i+3j+4k", "-1+2k", "-1", etc
type Quaternion struct {
	Real *RealPart
	I    *IPart `regexp:"?"`
	J    *JPart `regexp:"?"`
	K    *KPart `regexp:"?"`
}

// matches the quoted strings `"-1+2i+3j+4k"`, `"3-4k"`, `"12+34i"`, etc
type QuotedQuaternion struct {
	_          struct{} `regexp:"^"`
	_          struct{} `regexp:"\""`
	Quaternion *Quaternion
	_          struct{} `regexp:"\""`
	_          struct{} `regexp:"$"`
}

func (c *QuotedQuaternion) UnmarshalJSON(b []byte) error {
	if !quaternionRegexp.Find(c, string(b)) {
		return fmt.Errorf("%s is not a quaternion number", string(b))
	}
	return nil
}

// this struct is handled by JSON
type Var struct {
	Name  string
	Value *QuotedQuaternion
}

func prettyPrint(x interface{}) string {
	buf, err := json.MarshalIndent(x, "", "  ")
	if err != nil {
		return err.Error()
	}
	return string(buf)
}

func testRestructureQuotedQuaternion() {
	src := `{"name": "foo", "value": "1+2i+3j+4k"}`
	var v Var
	err := json.Unmarshal([]byte(src), &v)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(prettyPrint(v))
}

func testRestructureSubstring() {
	content := "joe@example.com"
	expr := regexp.MustCompile(`^([a-zA-Z0-9._%+-]+)@((\w+)\.(\w+))$`)
	indices := expr.FindStringSubmatchIndex(content)
	if len(indices) > 0 {
		userBegin, userEnd := indices[2], indices[3]
		var user string
		if userBegin != -1 && userEnd != -1 {
			user = content[userBegin:userEnd]
		}

		domainBegin, domainEnd := indices[6], indices[7]
		var domain string
		if domainBegin != -1 && domainEnd != -1 {
			domain = content[domainBegin:domainEnd]
		}

		tldBegin, tldEnd := indices[8], indices[9]
		var tld string
		if tldBegin != -1 && tldEnd != -1 {
			tld = content[tldBegin:tldEnd]
		}

		fmt.Println(user)   // prints "joe"
		fmt.Println(domain) // prints "example"
		fmt.Println(tld)    // prints "com"
	}
}

func testRestructureNested() {
	var addr EmailAddress2
	success, _ := restructure.Find(&addr, "joe@example.com")
	if success {
		fmt.Println(addr.User)        // prints "joe"
		fmt.Println(addr.Host.Domain) // prints "example"
		fmt.Println(addr.Host.TLD)    // prints "com"
	}
}

func testRestructure() {
	var addr EmailAddress
	restructure.Find(&addr, "joe@example.com")
	fmt.Println(addr.User) // prints "joe"
	fmt.Println(addr.Host) // prints "example.com"
}

/*
func testJsonql(jsonString string) {
	parser, err := jsonql.NewStringQuery(jsonString)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(parser.Query("name='elgs'"))
	//[map[skills:[Golang Java C] name:elgs gender:m age:35]] <nil>

	fmt.Println(parser.Query("name='elgs' && gender='f'"))
	//[] <nil>

	fmt.Println(parser.Query("age<10 || (name='enny' && gender='f')"))
	//[map[gender:f age:36 skills:[IC Electric design Verification] name:enny] map[age:1 skills:[Eating Sleeping Crawling] name:sam gender:m]] <nil>

	fmt.Println(parser.Query("age<10"))
	//[map[name:sam gender:m age:1 skills:[Eating Sleeping Crawling]]] <nil>

	fmt.Println(parser.Query("1=0"))
	//[] <nil>

	fmt.Println(parser.Query("age=(2*3)^2"))
	//[map[skills:[IC Electric design Verification] name:enny gender:f age:36]] <nil>

	fmt.Println(parser.Query("name ~= 'e.*'"))
	//[map[age:35 skills:[Golang Java C] name:elgs gender:m] map[skills:[IC Electric design Verification] name:enny gender:f age:36]] <nil>

	fmt.Println(parser.Query("name='el'+'gs'"))
	fmt.Println(parser.Query("age=30+5.0"))
	fmt.Println(parser.Query("age=40.0-5"))
	fmt.Println(parser.Query("age=70-5*7"))
	fmt.Println(parser.Query("age=70.0/2.0"))
	fmt.Println(parser.Query("age=71%36"))
	// [map[name:elgs gender:m age:35 skills:[Golang Java C]]] <nil>
}
*/

func csvWriterTest() {
	// Create `dump.csv` in `./shared/data` directory
	csv, err := ccsv.NewCsvWriter("./shared/data/dump.csv")
	if err != nil {
		panic("Could not open `sample.csv` for writing")
	}

	// Flush pending writes and close file upon exit of main()
	defer csv.Close()

	count := 99

	done := make(chan bool)

	for i := count; i > 0; i-- {
		go func(i int) {
			csv.Write([]string{strconv.Itoa(i), "bottles", "of", "beer"})
			done <- true
		}(i)
	}

	for i := 0; i < count; i++ {
		<-done
	}
}

func warcReadAllTest() {
	f, err := os.Open("./shared/testdata/test.warc")
	if err != nil {
		fmt.Println("error: ", err.Error())
		return
	}
	defer f.Close()

	records, err := warc.NewReader(f).ReadAll()
	if err != nil {
		fmt.Println("error: ", err)
		return
	}

	if len(records) <= 0 {
		fmt.Printf("record length mismatch: %d isn't enough records", len(records))
		return
	}

	for _, r := range records {
		fmt.Println(r.Type().String())
	}
}

func testMBP() {
	var wg sync.WaitGroup
	p := mpb.New(mpb.WithWaitGroup(&wg))
	total := 100
	numBars := 3
	wg.Add(numBars)

	for i := 0; i < numBars; i++ {
		name := fmt.Sprintf("Bar#%d:", i)
		bar := p.AddBar(int64(total),
			mpb.PrependDecorators(
				decor.StaticName(name, 0, 0),
				// DSyncSpace is shortcut for DwidthSync|DextraSpace
				// means sync the width of respective decorator's column
				// and prepend one extra space.
				decor.Percentage(3, decor.DSyncSpace),
			),
			mpb.AppendDecorators(
				decor.ETA(2, 0),
			),
		)
		go func() {
			defer wg.Done()
			for i := 0; i < total; i++ {
				time.Sleep(time.Duration(rand.Intn(10)+1) * time.Second / 100)
				bar.Increment()
			}
		}()
	}
	// Wait for incr loop goroutines to finish,
	// and shutdown mpb's rendering goroutine
	p.Stop()
}

func typedTest(path string) {
	// directly from a map[string]interace{}
	// typed := typed.New(a_map)
	// from a json []byte
	// typed, err := typed.Json(data)
	// from a file containing JSON
	typ, _ := typed.JsonFile(path)
	pp.Print(typ)
}

func simpleGet() {
	resp, err := resty.R().Get("http://httpbin.org/get") // GET request
	if err != nil {
		fmt.Println("error: ", err)
	}
	// explore response object
	fmt.Printf("\nError: %v", err)
	fmt.Printf("\nResponse Status Code: %v", resp.StatusCode())
	fmt.Printf("\nResponse Status: %v", resp.Status())
	fmt.Printf("\nResponse Time: %v", resp.Time())
	fmt.Printf("\nResponse Received At: %v", resp.ReceivedAt())
	fmt.Printf("\nResponse Body: %v", resp) // or resp.String() or string(resp.Body())
}

func goModel(req http.Request) {
	// let's say you have just decoded/unmarshalled your request body to struct object.
	tempPeople, _ := ParseJson(req.Body)
	people := People{}
	// tag your Product fields with appropriate options like
	// -, omitempty, notraverse to get desired result.
	// Not to worry, go-model does deep copy :)
	errs := model.Copy(&people, tempPeople)
	fmt.Println("Errors:", errs)

	fmt.Printf("\nSource: %#v\n", tempPeople)
	fmt.Printf("\nDestination: %#v\n", people)
}

func cmapTest() {
	m := cmap.New()
	m.Set("power", 9000)
	value, _ := m.Get("power")
	pp.Print(value)
	m.Delete("power")
	m.Len()
}

type People struct {
	Name  string `json:"name"`
	Sex   string `json:"sex"`
	Age   int    `json:"age"`
	Times int    `json:"times"`
}

// body as string
func gomergeTest(body []byte) {

	var tom People
	tom = People{
		Name:  "tom",
		Sex:   "male",
		Age:   18,
		Times: 1,
	}

	var request_data map[string]interface{}
	if err := json.Unmarshal(body, &request_data); err != nil {
		panic(err)
	}
	if err := gomerge.Merge(&tom, request_data); err != nil {
		panic(err)
	}
	result, _ := json.Marshal(tom)
	fmt.Println(result)
}

// c is a cache for api call
// var c = cache.New(cache.NoExpiration, cache.NoExpiration)

// var WorkQueue = make(chan WorkRequest, 10) // Simultaneous requests!
// var WorkQueue = make(chan WorkRequest, 1)    // Buffered channel to not lose anything
// var StopDispatcher = make(chan chan bool, 1) // Stop all!
// var apiToken = make(chan struct{}, 40)
// ref. https://github.com/parnurzeal/gorequest#endbytes
func requestMultiPart(endpoint string, payload string, filePath string, filename string, fieldname string) (resp gorequest.Response, body string, errs []error) {
	payload = `{"query1":"test"}`
	f, err := filepath.Abs(filePath)
	if err != nil {
		errs = append(errs, err)
		return nil, "", errs
	}

	bytesOfFile, err := ioutil.ReadFile(f)
	if err != nil {
		errs = append(errs, err)
		return nil, "", errs
	}

	request := gorequest.New().Timeout(2 * time.Millisecond)
	resp, body, errs = request.Post(endpoint).
		Type("multipart").
		Send(payload).SendFile(bytesOfFile, filename, fieldname).
		Retry(3, 5*time.Second, http.StatusBadRequest, http.StatusInternalServerError).
		End()
		/*
			RedirectPolicy(func(req Request, via []*Request) error {
				if req.URL.Scheme != "https" {
					return http.ErrUseLastResponse
				}
			}).
			// Set("Accept", "application/json").
			// AppendHeader("Accept", "application/json").
			// Set("Accept-Language", "en-US,en;q=0.8").
			// Set("Cache-Control", "max-age=0").
			// Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/47.0.2526.106 Safari/537.36").
		*/
		/*
			if len(errs) > 0 {
				return nil, errs[0]
			}
		*/

	return resp, body, errs

}

// StarResult wraps a star and an error
type ScraperResult struct {
	List  map[string][]Result
	Error error
}

func (e *Endpoint) ExecuteParallel(ctx context.Context, params map[string]string, resChan chan<- *ScraperResult) { // Execute will execute an Endpoint with the given params

	currentPage, _ := strconv.Atoi(e.Pager["offset"])
	lastPage, _ := strconv.Atoi(e.Pager["max"])

	offsetHolder := e.Pager["offset_var"]
	params[offsetHolder] = e.Pager["offset"]

	limitHolder := e.Pager["limit_var"]
	params[limitHolder] = e.Pager["limit"]
	for k, v := range e.Parameters {
		if _, ok := params[k]; !ok {
			if e.Debug {
				fmt.Printf("[WARNING] Parameters missing: k=%s, v=%s \n", k, v)
			}
		}
	}
	if e.Debug {
		fmt.Println("params")
		pp.Println(params)
	}
	for currentPage <= lastPage {
		res, err := e.Execute(params)
		if err != nil {
			resChan <- &ScraperResult{
				Error: err,
				List:  nil,
			}
		} else {
			if len(res) == 0 {
				lastPage = currentPage
				break
			} else {
				resChan <- &ScraperResult{
					Error: err,
					List:  res,
				}
			}
		}
		//if e.Debug {
		fmt.Println("res count: ", len(res), ", currentPage: ", currentPage, ", lastPage: ", lastPage)
		//}
		currentPage++ // Go to the next page
		params[offsetHolder] = strconv.Itoa(currentPage)
	}
	fmt.Println("closing resChan...")
	close(resChan)
}

/*
func enhancedGet() {
	resp, err := resty.R().
		SetQueryParams(map[string]string{
			"page_no": "1",
			"limit":   "20",
			"sort":    "name",
			"order":   "asc",
			"random":  strconv.FormatInt(time.Now().Unix(), 10),
		}).
		SetHeader("Accept", "application/json").
		SetAuthToken("BC594900518B4F7EAC75BD37F019E08FBC594900518B4F7EAC75BD37F019E08F").
		Get("/search_result")

	// Sample of using Request.SetQueryString method
	resp, err := resty.R().
		SetQueryString("productId=232&template=fresh-sample&cat=resty&source=google&kw=buy a lot more").
		SetHeader("Accept", "application/json").
		SetAuthToken("BC594900518B4F7EAC75BD37F019E08FBC594900518B4F7EAC75BD37F019E08F").
		Get("/show_product")
}
*/

func (e *Endpoint) getCacheKey(req *http.Request, debug bool) (string, string, string, error) {
	reqBytes, err := httputil.DumpRequest(req, true)
	if err != nil {
		return "", "", "", errors.New("dump request")
	}
	if debug {
		pp.Println(string(reqBytes))
	}
	cacheKey := fmt.Sprintf("%s_%x-%s_%s", e.hash, md5.Sum(reqBytes), req.Method, req.URL.String())
	cacheSlug := slugifier.Slugify(cacheKey)
	if e.Debug {
		fmt.Println("cacheSlug: ", cacheSlug)
	}
	cacheFile := fmt.Sprintf("./shared/cache/internal/%s.json", cacheSlug)
	return cacheKey, cacheSlug, cacheFile, nil
}

func (e *Endpoint) getHash(crypto string) (string, error) { // Execute will execute an Endpoint with the given params
	/*
		hash, err := structhash.Hash(e, 1)
		if err != nil {
			return "", err
		}
		fmt.Println("hash: ", hash)
		fmt.Println(structhash.Version(hash))
		if crypto == "md5" {
			fmt.Printf("structhash.Md5: %x\n", structhash.Md5(e, 1))
			fmt.Printf(" md5.Sum: %x\n", md5.Sum(structhash.Dump(e, 1)))
		}
		if crypto == "sha1" {
			fmt.Printf("structhash.Sha1: %x\n", structhash.Sha1(e, 1))
			fmt.Printf("sha1.Sum: %x\n", sha1.Sum(structhash.Dump(e, 1)))
		}
	*/
	return fmt.Sprintf("%x", structhash.Sha1(e, 1)), nil
}

// Simple JSON response generator
type Responder struct {
	flow.Component
	// In <-chan *RequestPacket
}

type RequestPacket struct {
	Req  *http.Request
	Res  http.ResponseWriter
	Code int
	Data interface{}
	Done chan bool
}

// Immediately pops the request with error response
func (p *RequestPacket) Error(code int, msg string) {
	p.Res.WriteHeader(code)
	js, _ := json.Marshal(Error{Code: code, Msg: msg})
	p.Res.Write(js)
	p.Done <- true
}

type GetRequestPacket struct {
	*RequestPacket
	Since int64
}

type PostRequestPacket struct {
	*RequestPacket
	Author string
	Text   string
}

type Error struct {
	Code int
	Msg  string
}

// https://github.com/Jeffail/tunny
func testTunny() {
	numCPUs := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPUs + 1) // numCPUs hot threads + one for async tasks.

	pool, _ := tunny.CreatePool(numCPUs, func(object interface{}) interface{} {
		input, _ := object.([]byte)

		// Do something that takes a lot of work
		output := input

		return output
	}).Open()

	defer pool.Close()

	http.HandleFunc("/work", func(w http.ResponseWriter, r *http.Request) {
		input, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Internal error", http.StatusInternalServerError)
		}

		// Send work to our pool
		result, _ := pool.SendWork(input)

		w.Write(result.([]byte))
	})

	http.ListenAndServe(":8080", nil)
}

func (e *Endpoint) getRake(body string) []string { // , excludeDomains []string) []string {
	var keywords []string
	candidates := rake.RunRake(body)
	for _, candidate := range candidates {
		keywords = append(keywords, candidate.Key)
		// keywords = append(keywords, fmt.Sprintf("%s", candidate.Key))
		fmt.Printf("%s --> %f\n", candidate.Key, candidate.Value)
	}
	fmt.Printf("\nsize: %d\n", len(keywords))
	return keywords
}

func (e *Endpoint) getLinks(url string) []string { // , excludeDomains []string) []string {
	var external_links []string
	if url == "" {
		return external_links
	}
	if e.Crawler.Limits.Parallelism <= 0 {
		e.Crawler.Limits.Parallelism = 5
	}
	/*
		if !e.Crawler.CSV.Disabled {
			if e.Crawler.CSV.PrefixPath == "" {
				file, err := os.Create(e.Crawler.CSV.PrefixPath)
			}
			if err != nil {
				log.Fatalf("Cannot create file %q: %s\n", e.Crawler.CSV.PrefixPath, err)
				return external_links
			}
			defer file.Close()
			writer := csv.NewWriter(file)
			defer writer.Flush()

			headers := strings.Split(e.Crawler.CSV.Headers, ",")
			writer.Write(headers) // Write CSV header
		}
	*/

	c := colly.NewCollector() // Instantiate default collector
	/*
		threads := make(map[string][]Mail)
		threadCollector := colly.NewCollector()
		mailCollector := colly.NewCollector()
	*/

	c.CacheDir = "./shared/cache/external/"
	c.UserAgent = "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36"

	// Limit the maximum parallelism to 5, This is necessary if the goroutines are dynamically created to control the limit of simultaneous requests.
	// Parallelism can be controlled also by spawning fixed number of go routines.
	c.Limit(&colly.LimitRule{DomainGlob: e.Crawler.Limits.DomainGlob, Parallelism: e.Crawler.Limits.Parallelism}) // , DisallowedDomains: ""})

	c.MaxDepth = 1                            // e.Crawler.MaxDepth           // MaxDepth is 1, so only the links on the scraped page and links on those pages are visited
	c.DisallowedDomains = []string{e.BaseURL} // Disallow url's domain owner
	// c.AllowedDomains = []string{"store.xkcd.com"} // Allow requests only to store.xkcd.com

	c.OnRequest(func(r *colly.Request) { // Before making a request put the URL with the key of "url" into the context of the request
		// pp.Println("OnRequest: ", r)
		r.Ctx.Put("url", r.URL.String())
	})

	c.OnResponse(func(r *colly.Response) { // After making a request get "url" from the context of the request
		// pp.Println("OnResponse: ", r)
		pp.Println(r.Ctx.Get("url"))
	})

	c.OnError(func(r *colly.Response, err error) { // Set error handler
		pp.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) { // On every a element which has href attribute call callback
		link := e.Attr("href")
		pp.Println(link) // Print link
		if strings.HasPrefix(link, "http") {
			external_links = append(external_links, link)
			// object, err := meta.ParseHTML(resp.Body)
			pp.Println("link: ", link)
		}
		go e.Request.Visit(link) // Visit link found on page on a new thread
	})

	c.Visit(url)  // Start scraping on https://en.wikipedia.org
	c.Wait()      // Wait until threads are finished
	pp.Println(c) // Display collector's statistics
	return external_links
}

/*
// Processes a request packet and sends the response JSON
func (r *Responder) OnIn(p *RequestPacket) {
	js, err := json.Marshal(p.Data)
	if err != nil {
		p.Error(http.StatusInternalServerError, "Could not marshal JSON")
		return
	}
	p.Res.Write(js)
	p.Done <- true
}
*/

// callCount is an API call counter used for debug
// var callCount uint16
// var callCountMutex sync.Mutex

func (e *Endpoint) Execute(params map[string]string) (map[string][]Result, error) { // Execute will execute an Endpoint with the given params
	if e.Debug {
		fmt.Println("endpoint handler config: ")
		pp.Println(e)
	}
	url, err := template(true, fmt.Sprintf("%s%s", e.BaseURL, e.PatternURL), params) //render url using params
	if err != nil {
		return nil, err
	}
	if e.Debug {
		fmt.Println("url: ", url)
	}

	method := e.Method //default method
	if method == "" {
		method = "GET"
	}

	body := io.Reader(nil) //render body (if set)
	if e.Body != "" {
		s, err := template(true, e.Body, params)
		if err != nil {
			return nil, err
		}
		body = strings.NewReader(s)
	}

	req, err := http.NewRequest(method, url, body) //create HTTP request
	if err != nil {
		return nil, err
	}

	if e.HeadersJSON != nil {
		for k, v := range e.HeadersJSON {
			if e.Debug {
				logf("use header %s=%s", k, v)
			}
			req.Header.Set(k, v)
		}
	}

	isResty := false
	if isResty {
		// https://github.com/go-resty/resty#various-post-method-combinations
		restyResp, err := resty.R().Get(url)
		if e.Debug {
			// explore response object
			fmt.Printf("\nError: %v", err)
			fmt.Printf("\nResponse Status Code: %v", restyResp.StatusCode())
			fmt.Printf("\nResponse Status: %v", restyResp.Status())
			fmt.Printf("\nResponse Time: %v", restyResp.Time())
			fmt.Printf("\nResponse Received At: %v", restyResp.ReceivedAt())
			fmt.Printf("\nResponse Body: %v", restyResp) // or resp.String() or string(resp.Body())
		}
	}

	var cacheFile string
	if e.Cache {
		if _, _, cacheFile, err = e.getCacheKey(req, e.Debug); err != nil {
			return nil, err
		}
	}

	if e.Cache && cacheFile != "" {
		isCacheExpired := cacheExpired(cacheFile, cacheDuration)
		if e.Debug && e.Cache {
			fmt.Printf("[ENDPOINT] isCacheExpired: %t\ncacheFile: %s \n", isCacheExpired, cacheFile)
		}
		if !isCacheExpired && e.Cache {
			return cacheContent(cacheFile)
		}
	}

	pp.Println("request: ", req)

	resp, err := getClient().Do(req)
	if err != nil {
		pp.Println(err)
		return nil, err
	}
	defer resp.Body.Close()

	if e.Debug { //show received headers
		fmt.Println("Response Headers: ")
		pp.Println(resp.Header)
		fmt.Println("Response Headers to intercept: ")
		pp.Println(e.HeadersIntercept)
	}

	for k, v := range resp.Header {
		if contains(e.HeadersIntercept, k) {
			if e.Debug {
				logf(" [INTERCEPT] header key=%s, value=%s", k, v)
			}
		}
	}

	fmt.Println("resp.StatusCode: ", resp.StatusCode)

	if resp.StatusCode != 200 {
		if e.Debug {
			logf("%s %s => %s", method, url, resp.Status)
		}
		if resp.StatusCode == 403 {
			if e.Debug {
				logf("%s %s =>\n%s\n", method, url, resp.Body)
			}
			return nil, errors.New("Unauthorized request")
		}
	}

	aggregate := make(map[string][]Result, 0)

	if e.Debug {
		fmt.Println("e.Selector: ", e.Selector)
	}

	switch e.Selector {
	case "wiki":
		if e.Debug {
			fmt.Println("Using 'WIKI' extractor")
		}
	case "md":
		if e.Debug {
			fmt.Println("Using 'MARKDOWN' extractor")
		}
	case "csv":
	case "tsv":
		if e.Debug {
			fmt.Printf("Using '%s-DELIMITED' extractor \n", e.Selector)
		}
	// https://stackoverflow.com/questions/24879587/xml-newdecoderresp-body-decode-giving-eof-error-golang
	case "xml":
		mv, err := mxj.NewMapXmlReader(resp.Body)
		if err != nil {
			return nil, err
		}
		if e.Debug {
			pp.Print(mv)
		}
		if e.ExtractPaths {
			mxj.LeafUseDotNotation()
			if e.Debug {
				fmt.Println("mv.LeafPaths(): ")
				pp.Println(mv.LeafPaths())
			}
			e.LeafPaths = leafPathsPatterns(mv.LeafPaths())
			if e.Debug {
				for _, v := range e.LeafPaths {
					fmt.Println("path:", v) // , "value:", v.Value)
				}
			}
		}
		for b, s := range e.BlocksJSON {
			if s.Items != "" {
				r := e.extractMXJ(mv, s.Items, s.Details)
				if e.Debug {
					fmt.Println("extractMXJ: ")
					pp.Println(r)
				}
				if r != nil {
					aggregate[b] = r
				}
			}
		}
	case "json":
		var mv mxj.Map
		var err error
		mxj.JsonUseNumber = true
		if e.Collection {
			mv, err = mxj.NewMapJsonArrayReaderAll(resp.Body)
		} else {
			mv, err = mxj.NewMapJsonReaderAll(resp.Body)
		}
		if err != nil {
			if e.Debug {
				fmt.Println("NewMapJsonReaderAll: ", err)
			}
			return nil, err
		}
		if e.ExtractPaths {
			mxj.LeafUseDotNotation()
			e.LeafPaths = leafPathsPatterns(mv.LeafPaths())
			if e.Debug {
				fmt.Println("mv.LeafPaths(): ")
				pp.Println(mv.LeafPaths())
				for _, v := range e.LeafPaths {
					fmt.Println("path:", v)
				}
			}
		}
		for b, s := range e.BlocksJSON {
			if e.Debug {
				pp.Println("s.Items: ", s.Items)
				pp.Println("s.Details: ", s.Details)
			}
			if s.Items != "" {
				r := e.extractMXJ(mv, s.Items, s.Details)
				if e.Debug {
					pp.Println(r)
				}
				if r != nil {
					aggregate[b] = r
				}
			}
			if e.Debug {
				fmt.Println(" - block_key: ", b)
				pp.Println(s)
			}
		}
	case "rss":
		fp := gofeed.NewParser()
		xml := resp.Body
		feed, err := fp.Parse(xml)
		if err != nil {
			return nil, err
		}
		// pp.Println("feed.Items: ", feed.Items)
		for b, s := range e.BlocksJSON {
			var results []Result
			if e.Debug {
				fmt.Println("[RSS] items count: ", len(feed.Items))
			}
			for _, item := range feed.Items {
				pp.Println("feed.Item: ", item)
				pp.Println("s.Details: ", s.Details)
				if item != nil {
					res := e.extractRss(item, s.Details)
					if len(res) > 0 {
						results = append(results, res)
					}
				}
			}
			/*
				// Additional info
				r["author"] = feed.Author
				r["categories"] = feed. Categories
				r["custom"] = feed.Custom
				r["copyright"] = feed.Copyright
				r["description"] = feed.Description
				r["type"] = feed.FeedType
				r["language"] = feed.Language
				r["title"] = feed.Title
				r["published"] = feed.Published
				r["updated"] = feed.Updated
			*/
			if len(results) > 0 {
				aggregate[b] = results
			}
		}
	case "xpath":
		doc, err := htmlquery.Parse(resp.Body)
		if err != nil {
			return nil, err
		}
		for b, s := range e.BlocksJSON {
			if s.Items != "" {
				if e.Debug {
					pp.Print(s)
				}
				var results []Result
				htmlquery.FindEach(doc, s.Items, func(i int, node *html.Node) {
					r := e.extractXpath(node, s.Details)
					if len(r) == len(s.Details) {
						results = append(results, r)
					} else if len(r) > 0 {
						if s.StrictMode == false {
							results = append(results, r)
						}
					}
					conform.Strings(r)
					if r != nil {
						r["id"] = strconv.Itoa(i)
						results = append(results, r)
					}
					if e.Debug {
						fmt.Print(" ---[ result: \n")
						pp.Print(r)
						fmt.Print(" ]---- \n")
					}
				})
				if results != nil {
					aggregate[b] = results
				}
			}
		}
	case "css":
		doc, err := goquery.NewDocumentFromReader(resp.Body) //parse HTML
		if err != nil {
			return nil, err
		}
		sel := doc.Selection

		for b, s := range e.BlocksJSON {
			if e.Debug {
				pp.Print(b)
				pp.Print(s)
			}
			var results []Result
			if s.Items != "" {
				sels := sel.Find(s.Items)
				if e.Debug {
					logf("list: %s => #%d elements", s.Items, sels.Length())
				}
				sels.Each(func(i int, sel *goquery.Selection) {
					r := e.extractCss(sel, s.Details)
					if len(r) == len(s.Details) && e.StrictMode {
						results = append(results, r)
					} else if len(r) > 0 && !e.StrictMode {
						results = append(results, r)
					} else if e.Debug {
						logf("excluded #%d: has %d fields, expected %d", i, len(r), len(s.Details))
					}
				})
				/*
					g := goose.New()
					article := g.ExtractFromURL(results["url"])
					println("title", article.Title)
					println("description", article.MetaDescription)
					println("keywords", article.MetaKeywords)
					println("content", article.CleanedText)
					println("url", article.FinalURL)
					println("top image", article.TopImage)
				*/
			} else {
				results[0] = e.extractCss(sel, s.Details)
				// results = append(results, e.extract(sel))
			}

			if results != nil {
				aggregate[b] = results
			}
		}
	default:
		fmt.Println("unkown selector type")
	}

	conform.Strings(aggregate)

	if len(aggregate) > 0 && e.Cache && cacheFile != "" {
		err = cacheResponse(cacheFile, aggregate) // dump response
		if err != nil {
			return nil, err
		}
	}

	//if e.Debug {
	//	aggregate["raw"] = resp.Body
	//}

	return aggregate, nil
}

// https://github.com/xrd/docker-search/blob/master/client.go#L70-L89

func cacheExpired(cacheFile string, maxAge time.Duration) bool {
	fi, err := os.Stat(cacheFile)
	if err != nil {
		return true
	}
	expireTime := fi.ModTime().Add(maxAge)
	fmt.Println("maxAge: ", maxAge)
	fmt.Println("expireTime: ", expireTime)
	fmt.Println("expired: ", time.Now().After(expireTime))
	return time.Now().After(expireTime)
}

func cacheContentRaw(cacheFile string) ([]byte, error) {
	file, err := os.Open(cacheFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	b, err := ioutil.ReadAll(file)
	return b, nil
}

func cacheContent(cacheFile string) (map[string][]Result, error) {
	aggregate := make(map[string][]Result, 0)
	file, err := ioutil.ReadFile(cacheFile)
	if err != nil {
		fmt.Printf("File error: %v\n", err)
		return nil, err
	}
	json.Unmarshal(file, &aggregate)
	return aggregate, nil
}

func cacheResponse(cacheSlug string, aggregate map[string][]Result) error {
	dump, err := json.Marshal(aggregate)
	// dump, err := json.MarshalIndent(aggregate, "", "    ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(cacheSlug, dump, 0644)
	if err != nil {
		return err
	}
	return nil
}

func Post(url string, params map[string]interface{}) ([]byte, error) {
	resp, err := resty.R().
		SetBody(params).
		Post(url)

	if err != nil {
		return nil, err
	}
	return resp.Body(), nil
}

func Get(url string, params map[string]string) ([]byte, error) {
	resp, err := resty.R().
		SetQueryParams(params).
		SetHeader("Accept", "application/json").
		Get(url)
	if err != nil {
		return nil, err
	}
	return resp.Body(), nil
}

func leafPathsPatterns(input []string) []string {
	var output []string
	var re = regexp.MustCompile(`.([0-9]+)`)
	for _, value := range input {
		value = re.ReplaceAllString(value, `[*]`)
		if !contains(output, value) {
			output = append(output, value)
		}
	}
	return dedup(output)
}

func sanitizeText(input string) string {
	// input =
	input = strings.Replace(input, " ...", "...", -1)
	input = strings.Replace(input, "\n", " ", -1)
	input = strings.Replace(input, "\t", " ", -1)
	input = strings.Trim(input, " ")
	return input
}

func (e *Endpoint) extractCss(sel *goquery.Selection, fields map[string]Extractors) Result { //extract 1 result using this endpoints extractor map
	r := Result{}
	if e.Debug {
		pp.Println(fields)
	}
	// ParseExtractor()
	for field, ext := range fields {
		if v := ext.execute(sel); v != "" {
			if field == "desc" {
				r["rake"] = e.getRake(v)
			}
			if field == "url" || field == "link" || field == "page_url" {
				link := v
				if !strings.HasPrefix(v, "http") {
					link = fmt.Sprintf("%s%s", e.BaseURL, v)
				}
				testGoose(link)
				r["external_links"] = e.getLinks(link)
			}
			if field == "url" && !strings.HasPrefix(v, "http") {
				r[field] = sanitizeText(strings.Trim(fmt.Sprintf("%s%s", e.BaseURL, v), " "))
			} else {
				r[field] = sanitizeText(strings.Trim(v, " "))
			}
		} else if e.Debug {
			logf("missing field: %s", field)
		}
	}
	return r
}

func (e *Endpoint) extractRss(item *gofeed.Item, fields map[string]Extractors) Result { //extract 1 result using this endpoints extractor map
	var fieldsList []string
	if e.Debug {
		pp.Println(fields)
	}
	for k, v := range fields {
		if e.Debug {
			pp.Println("fieldsList: k=", k, ", v=", v[0].val)
		}
		fieldsList = append(fieldsList, strcase.ToCamel(v[0].val))
	}
	r := Result{}
	// ParseExtractor()
	for _, field := range fieldsList {
		/*
			res, err := ParseExtractor(field)
			if err != nil {
				if e.Debug {
					fmt.Println("error: ", err)
				}
			}
		*/
		has, _ := reflections.HasField(item, field)
		if has {
			value, err := reflections.GetField(item, field)
			if err != nil {
				if e.Debug {
					fmt.Println("error: ", err)
				}
				continue
			}
			key := strings.ToLower(field)
			if e.Debug {
				pp.Println("reflected value: ", value)
			}
			if value != nil {
				if key == "desc" || key == "description" || key == "summary" {
					r["rake"] = e.getRake(fmt.Sprintf("%s", value))
				}
				if key == "url" || key == "link" || key == "page_url" {
					link := fmt.Sprintf("%s", value)
					if !strings.HasPrefix(link, "http") {
						link = fmt.Sprintf("%s%s", e.BaseURL, link)
					}
					testGoose(link)
					r["external_links"] = e.getLinks(link)
				}
				r[key] = value // sanitizeText(value)
			}
		}
		// pp.Println(" !!!! ParseExtractor result:", res)
	}
	if e.Debug {
		pp.Println("fields:", fields)
		pp.Println("fieldsList:", fieldsList)
		fmt.Println("item attr length:", len(r))
	}

	if e.Debug {
		fmt.Println("Results for RSS feed:")
		pp.Println(r)
	}

	return r
}

// import "github.com/jzaikovs/t"
func (e *Endpoint) extractMXJ(mv mxj.Map, items string, fields map[string]Extractors) []Result { //extract 1 result using this endpoints extractor map
	var r []Result
	if e.Debug {
		pp.Println(fields)
	}
	list, err := mv.ValuesForPath(items)
	if err != nil {
		fmt.Println("Error: ", err)
		// return nil
	}
	if e.Debug {
		pp.Println(list)
	}
	// ParseExtractor()
	for i := 0; i < len(list); i++ {
		l := Result{}
		for attr, field := range fields {
			var keyPath string
			var node []interface{}
			if len(field) == 1 {
				keyPath = fmt.Sprintf("%#s[%#d].%#s", items, i, field[0].val)
				if e.Debug {
					fmt.Println("field[0].val=", field[0].val, "keyPath: ", keyPath)
				}
				node, _ = mv.ValuesForPath(keyPath)
			} else {
				w := make(map[string]interface{}, len(field))
				var merr error
				for _, whl := range field {
					var keyName string
					if strings.Contains(whl.val, "|") {
						keyParts := strings.Split(whl.val, "|")
						if e.Debug {
							pp.Println(keyParts)
						}
						keyName = keyParts[len(keyParts)-1]
						whl.val = keyParts[0]
						if e.Debug {
							fmt.Println("keyName alias: ", keyName)
						}
					} else {
						keyParts := strings.Split(whl.val, ".")
						if e.Debug {
							pp.Println(keyParts)
						}
						keyName = keyParts[len(keyParts)-1]
						if e.Debug {
							fmt.Println("keyName alias", keyName)
						}
					}
					keyPath = fmt.Sprintf("%#s[%#d].%#s", items, i, whl.val)
					if e.Debug {
						fmt.Println("keyName: ", keyName, ", whl.vall=", whl.val, "keyPath: ", keyPath)
					}
					node, merr = mv.ValuesForPath(keyPath)
					if merr != nil {
						fmt.Println("Error: ", merr)
					}
					if node != nil {
						// conform.Strings(&node)
						if len(node) == 1 {
							w[keyName] = node[0]
						} else if len(node) > 1 {
							w[keyName] = node
						}
					}
				}
				if e.Debug {
					fmt.Println("subkeys whitelisted and mapped: ")
					pp.Println(w)
				}
				l[attr] = w
				continue
			}
			if len(node) == 1 {
				l[attr] = node[0]
			} else if len(node) > 1 {
				l[attr] = node
			}
			// conform.Strings(l)
		}
		r = append(r, l)
	}
	return r
}

func (e *Endpoint) extractXpath(node *html.Node, fields map[string]Extractors) Result { //extract 1 result using this endpoints extractor map
	if e.Debug {
		pp.Print(e)
	}
	r := Result{}
	// ParseExtractor()
	for field, ext := range fields {
		xpathRule := GetExtractorValue(ext)
		if e.Debug {
			logf("xpathRule: %s", xpathRule)
		}
		if v := htmlquery.FindOne(node, xpathRule); v != nil {
			t := htmlquery.InnerText(v)
			if e.Debug {
				logf("field %s, InnerText: %s", field, t) // fmt.Printf("field: %s \n", field)
			}
			switch field {
			case "page_url":
			case "link":
			case "url":
				url := htmlquery.SelectAttr(v, "href")
				if url == "" {
					return nil
				}
				if !strings.HasPrefix(url, "http") {
					url = fmt.Sprintf("%s%s", e.BaseURL, url)
				}
				testGoose(url)
				r["external_links"] = e.getLinks(url)
				if field == "url" && !strings.HasPrefix(url, "http") {
					r[field] = sanitizeText(strings.Trim(fmt.Sprintf("%s%s", e.BaseURL, url), " "))
				} else {
					r[field] = sanitizeText(strings.Trim(url, " "))
				}
			default:
				r[field] = sanitizeText(strings.Trim(t, " "))
			}
			if field == "desc" {
				r["rake"] = e.getRake(t)
			}
		} else if e.Debug {
			logf("missing field: %s", field)
		}
	}
	return r
}

type CacheResult struct {
	Id      uint64
	Name    string
	Payload int
}

func (cacheRes *CacheResult) CreateOne(name string, payload int) (uint64, error) {
	tx, err := BoltDB.Begin(true)
	if err != nil {
		log.Printf("[create] begin txn error: %v", err)
		return 0, err
	}
	defer tx.Rollback()

	bucket := tx.Bucket([]byte("user"))

	id, err := bucket.NextSequence()
	if err != nil {
		return 0, err
	}

	user := CacheResult{
		Id:      id,
		Name:    name,
		Payload: payload,
	}

	if data, err := json.Marshal(&user); err != nil {
		log.Printf("marshal error: %v", err)
		return 0, err
	} else if err := bucket.Put(intToByte(int(id)), data); err != nil {
		log.Printf("put error: %v", err)
		return 0, err
	}

	return id, tx.Commit()
}

func (cacheRes *CacheResult) Create(name string, payload int) error {
	var (
		payloadByte []byte
		err         error
	)

	tx, err := BoltDB.Begin(true)
	if err != nil {
		log.Printf("[create] begin txn error: %v", err)
		return err
	}
	defer tx.Rollback()

	b := tx.Bucket([]byte("user"))

	if payloadByte, err = json.Marshal(&payload); err != nil {
		log.Printf("[create] marshal error: %v", err)
		return err
	}

	nameByte, err := json.Marshal(&name)
	if err != nil {
		log.Printf("[create] marshal error: %v", err)
		return err
	}

	err = b.Put(nameByte, payloadByte)
	if err != nil {
		log.Printf("[create] put error: %v", err)
		return err
	}

	return tx.Commit()
}

func (cacheRes *CacheResult) GetOne(id uint64) (*CacheResult, error) {
	tx, err := BoltDB.Begin(false)
	if err != nil {
		log.Printf("[get] begin txn error: %v", err)
		return nil, err
	}
	defer tx.Rollback()

	var a CacheResult

	if v := tx.Bucket([]byte("user")).Get(intToByte(int(id))); v == nil {
		log.Print("get no record")
		return nil, nil
	} else if err := json.Unmarshal(v, &a); err != nil {
		log.Printf("unmarshal error: %v", err)
		return nil, err
	}

	return &a, nil
}

func (cacheRes *CacheResult) Get(name string) (int, error) {
	var (
		payload int
	)

	tx, err := BoltDB.Begin(false)
	if err != nil {
		log.Printf("[get] begin txn error: %v", err)
		return 0, err
	}
	defer tx.Rollback()

	nameByte, err := json.Marshal(&name)
	if err != nil {
		log.Printf("[get] marshal error: %v", err)
		return 0, err
	}

	if v := tx.Bucket([]byte("user")).Get(nameByte); v == nil {
		log.Print("[get] return nil value")
		return 0, nil
	} else if err := json.Unmarshal(v, &payload); err != nil {
		log.Printf("[get] unmarshal error: %v", err)
		return 0, err
	}

	return payload, nil
}
