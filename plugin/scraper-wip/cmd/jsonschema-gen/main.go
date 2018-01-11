package main

import (
	"fmt"
	"github.com/mcuadros/go-jsonschema-generator"
	// "github.com/roscopecoltran/scraper/scraper"
)

type EmbeddedType struct {
	Zoo string
}

type Item struct {
	Value string
}

type ExampleBasic struct {
	Foo bool   `json:"foo"`
	Bar string `json:",omitempty"`
	Qux int8
	Baz []string
	EmbeddedType
	List []Item
}

func main() {
	s := &jsonschema.Document{}
	s.Read(&ExampleBasic{})
	fmt.Println(s)

	// scraper

}
