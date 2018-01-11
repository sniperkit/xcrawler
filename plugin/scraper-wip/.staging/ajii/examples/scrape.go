package main

import (
	// "encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/narzeja/ajii"
	"log"
	// "os"
	"strings"
	"time"
)

type Buffer struct {
	Content []string
}

func (b *Buffer) Flush() string {
	body := strings.Join(b.Content[:], "\n")
	b.Content = b.Content[0:0] // clear buffer
	return body
}

func (b *Buffer) Push(elm string) {
	b.Content = append(b.Content, elm)
}

type Section struct {
	Title    string `json:"title"`
	TitleTag string `json:"title_tag"`
	Body     string `json:"body"`
}

type Result struct {
	Url       string    `json:"url"`
	Title     string    `json:"title"`
	Timestamp time.Time `json:"timestamp"`
	Sections  []Section `json:"sections"`
}

func Blacklisted(node *goquery.Selection) bool {
	list := []string{"script", "noscript", "meta", "link"}
	empty_text := strings.Replace(string(node.Text()), " ", "", -1)
	if len(empty_text) == 0 {
		return true
	}
	for _, elm := range list {
		if node.Is(elm) {
			return true
		}
	}
	return false
}

func NoTraversal(node *goquery.Selection) bool {
	nodes := []string{"p", "span", "td", "pre"}
	for _, elm := range nodes {
		if node.Is(elm) {
			return true
		}
	}
	return false
}

func Headline(node *goquery.Selection) bool {
	if string(goquery.NodeName(node)[0]) == "h" {
		nodes := []string{"h1", "h2", "h3", "h4", "h5", "h6"}
		for _, elm := range nodes {
			if node.Is(elm) {
				return true
			}
		}
	}
	return false
}

func FlushBuffer(result *Result, buffer *Buffer) {
	sec_len := len(result.Sections)
	if sec_len > 0 {
		prev_section := result.Sections[sec_len-1] // get previous struct
		prev_section.Body = buffer.Flush()
		result.Sections[sec_len-1] = prev_section
	}
}

func NewSection(current *goquery.Selection, result *Result) {
	title := strings.TrimSpace(current.Text())

	children := current.ChildrenFiltered("a")
	if children.Length() > 0 {
		href, _ := children.First().Attr("href")
		title = fmt.Sprintf("[%s](%s)", title, href)
	}

	section := Section{
		Title:    title,
		TitleTag: goquery.NodeName(current),
	}
	result.Sections = append(result.Sections, section)
}

func StripTags(current *goquery.Selection, result *Result) {
	body := Buffer{
		Content: []string{},
	}
	queue := []*goquery.Selection{current}

	for {
		if len(queue) == 0 {
			FlushBuffer(result, &body)
			break // reached end of file, no more to parse, flush the buffer
		}
		item_text := strings.TrimSpace(current.Text())
		if Blacklisted(current) {
			// skip meaningless nodes, like script and shit
		} else if Headline(current) { // don't care about children, just get text
			FlushBuffer(result, &body)
			NewSection(current, result)
		} else if current.Is("a") {
			href, _ := current.Attr("href")
			body.Push(fmt.Sprintf("[%s](%s)", item_text, href))
		} else if NoTraversal(current) {
			body.Push(item_text)
		} else {
			children := current.Children()
			if children.Length() == 0 { // no children for current node, just get content
				body.Push(item_text)
			} else {
				// So we do have children, lets queue them
				var newqueue []*goquery.Selection
				children.Each(func(index int, item *goquery.Selection) {
					newqueue = append(newqueue, item)
				})
				for _, old_itm := range queue {
					newqueue = append(newqueue, old_itm)
				}
				queue = newqueue
			}
		}
		current = queue[0]
		queue = queue[1:]
	}
	// result.Body = strings.Join(body[:], "\n")
}

// Perform depth first search in DOM of target page, strips out tags, and generates a non-html result of the page body
func Scrape(ctx *ajii.CallCtx) interface{} {
	url := ctx.Context.Query("url")
	if url == "" {
		fmt.Println("No query found")
		return struct {
			Error string `json:"error"`
		}{Error: "No url provided, use `url`-param"}
	}
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}
	title := doc.Find("title").Text()
	rp := Result{
		Url:       url,
		Title:     title,
		Timestamp: time.Now(),
		Sections: []Section{
			Section{
				Title: "EMPTY SECTION, PRE-BODY (PROBABLY MENU)"},
		},
	}

	current := doc.Find("body")
	StripTags(current, &rp)

	return rp
}

// func main() {
// 	config := ajii.NewConfig()
// 	service := ajii.NewService(config, "scrape", Scrape)
// 	cli := ajii.NewCli(config, service)
// 	cli.Run(os.Args)
// }

func main() {
	service := ajii.NewService("scrape", Scrape)
	service.Run()
}
