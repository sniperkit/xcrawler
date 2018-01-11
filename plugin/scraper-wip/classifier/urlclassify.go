package classifier

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"github.com/cloudflare/ahocorasick"
	"github.com/deckarep/golang-set"
	"io"
	"log"
	"os"
	"strings"
)

func createKeywordDictionaryFromCSVFile(name string) map[string]string {
	f, err := os.Open(name)

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	// Create a new reader.
	r := csv.NewReader(bufio.NewReader(f))
	m := make(map[string]string)
	for {
		record, err := r.Read()

		// Stop at EOF.
		if err == io.EOF {
			break
		}

		// Display record.
		//fmt.Printf("%v	%v\n", record[0], record[1])
		// Map the keyword to the category word...
		m[strings.ToLower(record[0])] = strings.ToLower(record[1])
		// Map the category word too...
		m[strings.ToLower(record[1])] = strings.ToLower(record[1])
	}

	return m
}

func createListFromCSVFile(name string) []string {
	f, err := os.Open(name)

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	// Create a new reader.
	r := csv.NewReader(bufio.NewReader(f))
	a := make([]string, 0)
	for {
		record, err := r.Read()

		// Stop at EOF.
		if err == io.EOF {
			break
		}

		a = append(a, record[0])
	}

	return a
}

func getKeys(m map[string]string) []string {
	i := 0
	keys := make([]string, len(m))
	for k := range m {
		keys[i] = k
		i++
	}

	return keys
}

func getKeywords(m map[string]string) []string {
	return getKeys(m)
}

func getUniqueCategoriesFromHits(hits []int, keywords []string, dict map[string]string) []interface{} {
	set := mapset.NewSet()
	for _, v := range hits {
		set.Add(dict[keywords[v]])
	}

	return set.ToSlice()
}

func urlClassifierMain() {
	// Load urls...
	urls := createListFromCSVFile("./shared/dataset/newurls.csv")

	// Load keywords...
	dict := createKeywordDictionaryFromCSVFile("./shared/dataset/keywords.csv")

	// Classify...
	keywords := getKeywords(dict)
	m := ahocorasick.NewStringMatcher(keywords)
	for _, url := range urls {
		hits := m.Match([]byte(strings.ToLower(url)))
		//fmt.Printf("# of hits for %s: %d\n", url, len(hits))
		categories := getUniqueCategoriesFromHits(hits, keywords, dict)
		for _, v := range categories {
			fmt.Printf("%s,%s\n", url, v)
		}
	}
}
