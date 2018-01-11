package scraper

import (
	"fmt"

	"github.com/mcuadros/go-jsonschema-generator"
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

func ConvertToJsonSchema() {
	s := &jsonschema.Document{}
	s.Read(&ExampleBasic{})
	fmt.Printf("ExampleBasic schema:\n %s \n\n", s)

	e := &jsonschema.Document{}
	e.Read(&Endpoint{})
	fmt.Printf("Endpoint schema:\n %s \n\n", e)

	c := &jsonschema.Document{}
	c.Read(&Config{})
	fmt.Printf("Config schema:\n %s \n\n", c)

}
