package scrapper

import (
	// "bytes"
	"fmt"
	// "io/ioutil"
	"github.com/PuerkitoBio/goquery"
	model "github.com/SurgeNews/SurgeServer/model"
	"log"
	"net/url"
)

type Client struct {
	rootAddress string
	path        string
	queries     []string
	ExpirySec   int
}

func NewClient() *Client {
	client := new(Client)
	client.rootAddress = "https://news.google.co.in"
	client.path = "news/section"
	client.queries = []string{"n", "s", "b"}
	return client
}

func (self *Client) Request(index int) error {

	var err error

	uri, err := url.ParseRequestURI(self.rootAddress)
	if err != nil {
		return err
	}

	uri.Path = self.path

	q := url.Values{}
	q.Set("topic", self.queries[index])
	uri.RawQuery = q.Encode()

	doc, err := goquery.NewDocument(uri.String())
	if err != nil {
		log.Fatal(err)
	}

	// Find the review items
	doc.Find(".esc-layout-article-cell").Each(func(i int, s *goquery.Selection) {

		// For each item found, get the band and title
		link, _ := s.Find("a").Attr("href")
		title := s.Find(".esc-lead-article-title span").Text()
		source := s.Find(".al-attribution-source").Text()
		// title := s.Find("h2").Text()
		fmt.Printf("Review %d: %s \n %s\n\t\t %s\n\n", i, link, title, source)
	})

	return nil
}
