package scraper

import (
	"encoding/csv"
	"log"
	"os"
	"path/filepath"

	"github.com/k0kubun/pp"
	// "github.com/agrison/go-tablib"
	// "github.com/inotom/csv2json"
	// "github.com/mmcloughlin/databundler"
	// "github.com/tsak/concurrent-csv-writer"
)

// https://github.com/angeldm/optiqor/blob/master/db/seeds/product.go
// https://github.com/angeldm/optiqor/blob/master/config/admin/admin.go

var Seeds = struct {
	Topics []struct {
		Name string
	}
	Groups []struct {
		Name string
	}
}{}

func SeedAlexaTop1M() {
	filepaths, _ := filepath.Abs("./shared/seeds/top-1m.csv")
	f, _ := os.Open(filepaths)
	r := csv.NewReader(f)
	r.Comma = '|'
	r.Comment = '#'
	r.LazyQuotes = true

	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < 100; i++ {
		pp.Println("provider: ", records[i])
		// for i := 0; i < len(records); i++ {
		// fmt.Println("provider: ", records[i])
		// FindOrCreateProviderByName(records[i])
		// createProduct(records[i])
	}

	log.Printf("Imported %d providers.", len(records))

}
