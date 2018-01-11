package scraper

import (
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
	"unsafe"

	"github.com/PuerkitoBio/goquery"
	"github.com/jinzhu/now"
	"github.com/k0kubun/pp"
	"github.com/mholt/archiver"
	"github.com/plar/go-adaptive-radix-tree"
	"github.com/schollz/progressbar"
	// "github.com/jinzhu/inflection"
	// "github.com/jinzhu/copier"
	// "github.com/linkosmos/mapdecor"
	// "github.com/toukii/jsnm"
	// "github.com/byrnedo/mapcast"
	// "github.com/spf13/cast"
	// "github.com/rakanalh/goscrape"
	// "github.com/rakanalh/goscrape/extract"
	// "github.com/rakanalh/goscrape/processors"
	// "github.com/gorilla/css/scanner"
	// "github.com/moovweb/gokogiri"
	// "github.com/moovweb/gokogiri/html"
	// "github.com/moovweb/gokogiri/xml"
	// "github.com/moovweb/gokogiri/xpath"
	// "github.com/roscopecoltran/css2xpath"
)

/*
	Refs:
	- https://github.com/osvik/txttransformer
	- https://github.com/hotei/tempfile
	- https://github.com/gonyi/yIndex
	- https://github.com/nightmouse/funcsplit
	- https://github.com/danward79/csvtool
	- https://github.com/doganov/filesort-example
	- https://github.com/vbauerster/mpb/blob/master/examples/io/multiple/main.go
	- https://github.com/vbauerster/mpb/blob/master/examples/cancel/main.go
	- https://github.com/vbauerster/mpb/blob/master/examples/remove/main.go
	- https://github.com/vbauerster/mpb/blob/master/examples/sort/main.go
	- https://github.com/vbauerster/mpb/blob/master/examples/prependETA/main.go
	- https://github.com/alin-sinpalean/concurrent-writer
*/

var templateRe = regexp.MustCompile(`\{\{\s*(\w+)\s*(:(\w+))?\s*\}\}`)

func testAdaptiveRadixTree() {

	tree := art.New()

	tree.Insert(art.Key("Hi, I'm Key"), "Nice to meet you, I'm Value")
	value, found := tree.Search(art.Key("Hi, I'm Key"))
	if found {
		fmt.Printf("Search value=%v\n", value)
	}

	tree.ForEach(func(node art.Node) bool {
		fmt.Printf("Callback value=%v\n", node.Value())
		return true
	})

	for it := tree.Iterator(); it.HasNext(); {
		value, _ := it.Next()
		fmt.Printf("Iterator value=%v\n", value.Value())
	}
}

func testProgressBar() {
	bar := progressbar.New(100)
	for i := 0; i < 100; i++ {
		bar.Add(1)
		time.Sleep(10 * time.Millisecond)
	}
}

func createArchive(outputFile string, files []string) {
	if err := archiver.Zip.Make(outputFile, files); err != nil {
		fmt.Println("Could not create archive...")
	}
}

func extractArchive(inputFile string, outputFolder string) {
	if err := archiver.Zip.Open(inputFile, outputFolder); err != nil {
		fmt.Println("Could not extract archive...")
	}
}

func intToByte(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

func str2bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

func bytes2str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func stack() []byte {
	buf := make([]byte, 10240)
	n := runtime.Stack(buf, false)
	if n > 710 {
		copy(buf, buf[710:n])
		return buf[:n-710]
	}
	return buf[:n]
}

func contains(input []string, match string) bool {
	for _, value := range input {
		if value == match {
			return true
		}
	}
	return false
}

func dedup(input []string) []string {
	var output []string
	for _, value := range input {
		if !contains(output, value) {
			output = append(output, value)
		}
	}
	return output
}

func ParseJson(body io.ReadCloser) (jsonBody map[string]interface{}, err error) {
	bytes, err := ioutil.ReadAll(body)

	if err != nil {
		return nil, err
	}

	json.Unmarshal(bytes, &jsonBody)

	return
}

func template(isurl bool, str string, vars map[string]string) (out string, err error) {
	out = templateRe.ReplaceAllStringFunc(str, func(key string) string {
		m := templateRe.FindStringSubmatch(key)
		k := m[1]
		value, ok := vars[k]
		if !ok { //missing - apply defaults or error
			if m[3] != "" {
				value = m[3]
			} else {
				err = errors.New("Missing param: " + k)
			}
		}
		if isurl { //determine if we need to escape
			queryi := strings.Index(str, "?")
			keyi := strings.Index(str, key)
			if queryi != -1 && keyi > queryi {
				value = url.QueryEscape(value)
			}
		}
		return value
	})
	return
}

func checkSelector(s string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	doc, _ := goquery.NewDocumentFromReader(bytes.NewBufferString(`<html>
		<body>
			<h3>foo bar</h3>
		</body>
	</html>`))
	doc.Find(s)
	return
}

func jsonerr(err error) []byte {
	return []byte(`{"error":"` + err.Error() + `"}`)
}

func jsoncache(content []byte) []byte {
	return []byte(content)
}

func logf(format string, args ...interface{}) {
	log.Printf("[scraper] "+format, args...)
}

func HasElem(s interface{}, elem interface{}) bool {
	arrV := reflect.ValueOf(s)
	if arrV.Kind() == reflect.Slice {
		for i := 0; i < arrV.Len(); i++ {
			// XXX - panics if slice element points to an unexported struct field
			// see https://golang.org/pkg/reflect/#Value.Interface
			if arrV.Index(i).Interface() == elem {
				return true
			}
		}
	}

	return false
}

func isJSONString(s string) bool {
	var js string
	return json.Unmarshal([]byte(s), &js) == nil

}

func isJSON(s string) bool {
	var js map[string]interface{}
	return json.Unmarshal([]byte(s), &js) == nil

}

func isJsonArray(s string) bool {
	var js []interface{}
	return json.Unmarshal([]byte(s), &js) == nil
}

func debugHttpReqResp(req *http.Request, resp *http.Response) {
	reqDump, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		log.Fatalln("error while loging request", err)
	}
	fmt.Printf("--- REQUEST START ---\n%s\n--- REQUEST END ---", reqDump)
	respDump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		log.Fatalln("error while loging response", err)
	}
	fmt.Printf("--- RESPONSE START ---\n%s\n--- RESPONSE END ---", respDump)
}

func generateCacheKey(req *http.Request, debug bool) (string, error) {
	reqBytes, err := httputil.DumpRequest(req, true)
	if err != nil {
		return "", errors.New("dump request")
	}
	if debug {
		pp.Println(string(reqBytes))
	}
	return fmt.Sprintf("%s-%s-%x", req.Method, req.URL.String(), md5.Sum(reqBytes)), nil
}

func testtime() {
	time.Now() // 2013-11-18 17:51:49.123456789 Mon

	now.BeginningOfMinute()   // 2013-11-18 17:51:00 Mon
	now.BeginningOfHour()     // 2013-11-18 17:00:00 Mon
	now.BeginningOfDay()      // 2013-11-18 00:00:00 Mon
	now.BeginningOfWeek()     // 2013-11-17 00:00:00 Sun
	now.FirstDayMonday = true // Set Monday as first day, default is Sunday
	now.BeginningOfWeek()     // 2013-11-18 00:00:00 Mon
	now.BeginningOfMonth()    // 2013-11-01 00:00:00 Fri
	now.BeginningOfQuarter()  // 2013-10-01 00:00:00 Tue
	now.BeginningOfYear()     // 2013-01-01 00:00:00 Tue

	now.EndOfMinute()         // 2013-11-18 17:51:59.999999999 Mon
	now.EndOfHour()           // 2013-11-18 17:59:59.999999999 Mon
	now.EndOfDay()            // 2013-11-18 23:59:59.999999999 Mon
	now.EndOfWeek()           // 2013-11-23 23:59:59.999999999 Sat
	now.FirstDayMonday = true // Set Monday as first day, default is Sunday
	now.EndOfWeek()           // 2013-11-24 23:59:59.999999999 Sun
	now.EndOfMonth()          // 2013-11-30 23:59:59.999999999 Sat
	now.EndOfQuarter()        // 2013-12-31 23:59:59.999999999 Tue
	now.EndOfYear()           // 2013-12-31 23:59:59.999999999 Tue

	// Use another time
	t1 := time.Date(2013, 02, 18, 17, 51, 49, 123456789, time.Now().Location())
	now.New(t1).EndOfMonth() // 2013-02-28 23:59:59.999999999 Thu

	// Don't want be bothered with the First Day setting, Use Monday, Sunday
	now.Monday()      // 2013-11-18 00:00:00 Mon
	now.Sunday()      // 2013-11-24 00:00:00 Sun (Next Sunday)
	now.EndOfSunday() // 2013-11-24 23:59:59.999999999 Sun (End of next Sunday)

	t2 := time.Date(2013, 11, 24, 17, 51, 49, 123456789, time.Now().Location()) // 2013-11-24 17:51:49.123456789 Sun
	now.New(t2).Monday()                                                        // 2013-11-18 00:00:00 Sun (Last Monday if today is Sunday)
	now.New(t2).Sunday()                                                        // 2013-11-24 00:00:00 Sun (Beginning Of Today if today is Sunday)
	now.New(t2).EndOfSunday()                                                   // 2013-11-24 23:59:59.999999999 Sun (End of Today if today is Sunday)

	/*
	   time.Now() // 2013-11-18 17:51:49.123456789 Mon

	   // Parse(string) (time.Time, error)
	   t, err := now.Parse("12:20")            // 2013-11-18 12:20:00, nil
	   t, err := now.Parse("1999-12-12 12:20") // 1999-12-12 12:20:00, nil
	   t, err := now.Parse("99:99")            // 2013-11-18 12:20:00, Can't parse string as time: 99:99

	   // MustParse(string) time.Time
	   now.MustParse("2013-01-13")             // 2013-01-13 00:00:00
	   now.MustParse("02-17")                  // 2013-02-17 00:00:00
	   now.MustParse("2-17")                   // 2013-02-17 00:00:00
	   now.MustParse("8")                      // 2013-11-18 08:00:00
	   now.MustParse("2002-10-12 22:14")       // 2002-10-12 22:14:00
	   now.MustParse("99:99")                  // panic: Can't parse string as time: 99:99
	*/
}

func gzipFast(a *[]byte) []byte {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	if _, err := gz.Write(*a); err != nil {
		gz.Close()
		panic(err)
	}
	gz.Close()
	return b.Bytes()
}

func ReadJson(r *http.Request, data interface{}) bool {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(data)
	if err != nil {
		panic(err)
	}
	return true
}

func WriteJson(w http.ResponseWriter, data interface{}) {
	prefix := []byte(")]}',\n")
	if d, err := json.Marshal(data); err != nil {
		log.Printf("Error marshalling json: %v", err)
	} else {
		d = append(prefix, d...)
		w.Header().Set("Content-Length", strconv.Itoa(len(d)))
		w.Header().Set("Content-Type", "application/json")
		w.Write(d)
	}
}

func ErrorJson(w http.ResponseWriter, data interface{}, errcode int) {
	if d, err := json.Marshal(data); err != nil {
		log.Printf("Error marshalling json: %v", err)
	} else {
		w.WriteHeader(errcode)
		w.Header().Set("Content-Length", strconv.Itoa(len(d)))
		w.Header().Set("Content-Type", "application/json")
		w.Write(d)
	}
}

func indentLen(s string) int {
	i := 0
	for i < len(s) && (s[i] == ' ' || s[i] == '\t') {
		i++
	}
	return i
}

func isBlank(s string) bool {
	return len(s) == 0 || (len(s) == 1 && s[0] == '\n')
}

func commonPrefix(a, b string) string {
	i := 0
	for i < len(a) && i < len(b) && a[i] == b[i] {
		i++
	}
	return a[0:i]
}

func unindent(block []string) {
	if len(block) == 0 {
		return
	}

	// compute maximum common white prefix
	prefix := block[0][0:indentLen(block[0])]
	for _, line := range block {
		if !isBlank(line) {
			prefix = commonPrefix(prefix, line[0:indentLen(line)])
		}
	}
	n := len(prefix)

	// remove
	for i, line := range block {
		if !isBlank(line) {
			block[i] = line[n:]
		}
	}
}

// heading returns the trimmed line if it passes as a section heading;
// otherwise it returns the empty string.
func heading(line string) string {
	line = strings.TrimSpace(line)
	if len(line) == 0 {
		return ""
	}

	// a heading must start with an uppercase letter
	r, _ := utf8.DecodeRuneInString(line)
	if !unicode.IsLetter(r) || !unicode.IsUpper(r) {
		return ""
	}

	// it must end in a letter or digit:
	r, _ = utf8.DecodeLastRuneInString(line)
	if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
		return ""
	}

	// exclude lines with illegal characters
	if strings.IndexAny(line, ",.;:!?+*/=()[]{}_^°&§~%#@<\">\\") >= 0 {
		return ""
	}

	// allow "'" for possessive "'s" only
	for b := line; ; {
		i := strings.IndexRune(b, '\'')
		if i < 0 {
			break
		}
		if i+1 >= len(b) || b[i+1] != 's' || (i+2 < len(b) && b[i+2] != ' ') {
			return "" // not followed by "s "
		}
		b = b[i+2:]
	}

	return line
}

func Atoi(s string) int32 {
	i, _ := strconv.ParseInt(s, 10, 32)
	return int32(i)
}

func Itoa(i int32) string {
	return strconv.FormatInt(int64(i), 10)
}

func Atoi64(s string) int64 {
	i, _ := strconv.ParseInt(s, 10, 64)
	return i
}

// https://github.com/thbourlove/restc/blob/master/transport.go
// https://github.com/lox/package-proxy/blob/master/cache/http.go#L39
// https://golang.org/pkg/net/http/httputil/#example_DumpRequest

/*
func mapdecor() {
    input := map[string]interface{}{
      "key1": nil,
      "key2": "with",
      "val1": "1",
      "val2": "2",
      "val3": "3",
      "val4": "4",
    }

    decorFunc := func(input map[string]interface{}) (output map[string]interface{}) {
      partitonFunc := func(s string, i interface{}) bool {
        return strings.Contains(s, "val")
      }

      // For first (inputPartitioned[0]) partition we get key-values containing 'val' in key
      // For second (inputPartitioned[1]) partition what is left
      inputPartitioned := mapop.Partition(partitonFunc, input)

      // Assigning values key to first partition
      inputPartitioned[1]["values"] = inputPartitioned[0]

      return inputPartitioned[1]
    }


    got := Decorate(input, decorFunc)

  // Got
  // map[string]interface{}{
  //   "key1": nil,
  //   "key2": "with",
  //   "values": map[string]interface{}{
  //     "val1": "1",
  //     "val2": "2",
  //     "val3": "3",
  //     "val4": "4",
  //   }
  // }
}
*/

/*
// https://github.com/byrnedo/mapcast
type myStruct struct {
    Field int `json:"input_name" bson:"output_name"`
    Hidden float32 `json:"-" bson:"hidden_field"`
}

func mapcast() {
	myInput := map[string]string{"input_name": "345"}

	Cast(myInput, myStruct{})
	// returns map{"Field" : 345}

	CastViaJson(myInput, myStruct{})
	// returns map{"input_name" : 345}

	CastViaJsonToBson(myInput, myStruct{})
	// returns map{"output_name" : 345}

	myMultiInput := map[string][]string{"input_name" : []string{"345}}

	CastMultiViaJsonToBson(myMultiInput, myStruct{})
	// returns map{"output_name" : []int[345]}
}
*/

/*
func xPathToCss(xpath []string, xtype string) []string {
	fmt.Printf("xpath type: %s \n", xtype)
	var result []string
	for _, css := range xpath {
		switch xtype {
		case "local":
			result = append(result, css2xpath.Convert(css, css2xpath.LOCAL))
		default:
			result = append(result, css2xpath.Convert(css, css2xpath.GLOBAL))
		}
	}
	return result
}
*/
