package scraper

import (
	"fmt"
	"strings"

	"github.com/jdkato/prose/chunk"
	"github.com/jdkato/prose/summarize"
	"github.com/jdkato/prose/tag"
	"github.com/jdkato/prose/tokenize"
	"github.com/jdkato/prose/transform"
)

func testTokenizer() {
	text := "They'll save and invest more."
	tokenizer := tokenize.NewTreebankWordTokenizer()
	for _, word := range tokenizer.Tokenize(text) {
		// [They 'll save and invest more .]
		fmt.Println(word)
	}
}

func testTagger() {
	text := "A fast and accurate part-of-speech tagger for Golang."
	words := tokenize.NewTreebankWordTokenizer().Tokenize(text)

	tagger := tag.NewPerceptronTagger()
	for _, tok := range tagger.Tag(words) {
		fmt.Println(tok.Text, tok.Tag)
	}
}

func testTransformer() {
	text := "the last of the mohicans"
	tc := transform.NewTitleConverter(transform.APStyle)
	fmt.Println(strings.Title(text)) // The Last Of The Mohicans
	fmt.Println(tc.Title(text))      // The Last of the Mohicans
}

func testSummarizer() {
	doc := summarize.NewDocument("This is some interesting text.")
	fmt.Println(doc.SMOG(), doc.FleschKincaid())
}

func testChunker() {
	words := tokenize.TextToWords("Go is an open source programming language created at Google.")
	regex := chunk.TreebankNamedEntities

	tagger := tag.NewPerceptronTagger()
	for _, entity := range chunk.Chunk(tagger.Tag(words), regex) {
		fmt.Println(entity) // [Go Google]
	}
}
