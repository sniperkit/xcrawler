package main

import (
	"net/http"
	"strconv"
	"text/template"

	"github.com/pkg/browser"
)

func webCrawlerMain(w http.ResponseWriter, r *http.Request) {
	htmlBody := `<h1>Crawler</h1>
    <p><h2><a href="/crawlerQuora">Quora</a></p>
    <p><h2><a href="/crawlerMedium">Medium</a></p>
    `

	w.Write([]byte(htmlBody))
}

func webMain(w http.ResponseWriter, r *http.Request) {
	htmlBody := `<h1>Mail Classifier</h1>
    <p><h2><a href="/gmailFetch">E-Mails from Gmail</a></p>
    <p><h2><a href="/crawlerMain">Crawler</a></p>
    `

	w.Write([]byte(htmlBody))
}

func main() {
	http.HandleFunc("/", webMain)
	http.HandleFunc("/crawlerMain", webCrawlerMedium)
	http.ListenAndServe(":8080", nil)
	browser.OpenURL("http://localhost:8080")
}
