package main

import (
	"log"

	flow "github.com/wanliu/goflow"
)

type GetElement struct {
	flow.Component
}

type Split struct {
	flow.Component
}

type Output struct {
	flow.Component
	In      <-chan string
	Options <-chan map[string]interface{}
	Out     chan<- string
}

func NewGetElement() interface{} {
	return new(GetElement)
}

func NewSplit() interface{} {
	return new(Split)
}

func NewOutput() interface{} {
	return new(Output)
}

func (o *Output) OnIn(msg string) {
	log.Printf("output: %s", msg)
	o.Out <- msg
}
