package main

import (
	"fmt"
	"github.com/ekzhu/go-fasttext"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	ft := fasttext.NewFastText("./sqlite3")
	fmt.Println("vim-go")
	fmt.Println(ft)
}
