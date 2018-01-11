package main

import (
	"fmt"
	"log"

	"github.com/PuerkitoBio/goquery"
)

/*

refs:
- https://raw.githubusercontent.com/suntong/kuaizh/master/crawler/duckduckgo.go
- https://github.com/frizinak/tiif/blob/master/engine/duckduckgo.go

eg: https://duckduckgo.com/html/?q=whitesmith

 <h2 class="result__title">

   <a rel="nofollow" class="result__a" href="https://en.wikipedia.org/wiki/Whitesmith">Tinsmith - Wikipedia, the free encyclopedia</a>

 </h2>

*/

func scrape(url string) {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find("a.result__a").Each(func(i int, s *goquery.Selection) {
		link, exists := s.Attr("href")
		if !exists {
			return
		}
		fmt.Println(link)
	})
}

func concurrentQueries() {
	go scrape("http://duckduckgo.com/html/?q=whitesmith")
	go scrape("http://duckduckgo.com/html/?q=coimbra")
	go scrape("http://duckduckgo.com/html/?q=adbjesus")
}

func main() {
	scrape("http://duckduckgo.com/html/?q=whitesmith")
}
