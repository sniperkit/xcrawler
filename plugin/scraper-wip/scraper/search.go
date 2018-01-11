package scraper

import (
	"log"

	"github.com/go-ego/riot"
	"github.com/go-ego/riot/types"
)

// https://github.com/blevesearch/bleve#indexing
// https://github.com/blevesearch/bleve#querying

var (
	// searcher is coroutine safe
	searcher = riot.Engine{}
)

func testRiot() {
	// Init
	searcher.Init(types.EngineInitOptions{
		Using:             5,
		NotUsingSegmenter: true})
	defer searcher.Close()

	// Add the document to the index, docId starts at 1
	searcher.IndexDocument(1, types.DocIndexData{Content: "Google Is Experimenting With Virtual Reality Advertising"}, false)
	searcher.IndexDocument(2, types.DocIndexData{Content: "Google accidentally pushed Bluetooth update for Home speaker early"}, false)
	searcher.IndexDocument(3, types.DocIndexData{Content: "Google is testing another Search results layout with rounded cards, new colors, and the 4 mysterious colored dots again"}, false)

	// Wait for the index to refresh
	searcher.FlushIndex()

	// The search output format is found in the types.SearchResponse structure
	log.Print(searcher.Search(types.SearchRequest{Text: "google testing"}))
}
