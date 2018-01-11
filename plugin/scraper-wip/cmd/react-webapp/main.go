package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/wanliu/goflow"
	// "github.com/k0kubun/pp"
	// "golang.org/x/net/context"
)

var _components = make(map[string]func() interface{})

var Info = Package{
	Name:        "goflow-components",
	Description: "goflow builtin components",
	Version:     "0.0.1",
}

// struct Greeter used communication
type Greeter struct {
	flow.Component
	Name <-chan string
	Res  chan<- string
}

//
func (g *Greeter) OnName(name string) {
	greeting := fmt.Sprintf("Hello, %s!", name)
	g.Res <- greeting
}

type Logger struct {
	flow.Component
	Line <-chan string
	Res  chan<- string
}

func (p *Logger) OnLine(line string) {
	fmt.Println(line)
	p.Res <- line
}

type GreetingFlow struct {
	flow.Graph
}

func NewGreetingFlow() *GreetingFlow {
	n := new(GreetingFlow)
	n.InitGraphState()
	n.Add(new(Greeter), "greeter")
	// n.Add(new(Logger), "logger")
	n.Add(new(Output), "output")

	// n.Connect("greeter", "Res", "logger", "Line")
	n.Connect("greeter", "Res", "output", "In")

	n.MapInPort("In", "greeter", "Name")
	n.MapOutPort("Out", "output", "Out")
	return n
}

func main() {
	net := NewGreetingFlow()
	in := make(chan string)
	out := make(chan string)
	net.SetInPort("In", in)
	net.SetOutPort("Out", out)
	flow.RunNet(net)

	_components["dom/GetElement"] = NewGetElement
	_components["core/Split"] = NewSplit
	_components["Output"] = NewOutput

	router := gin.Default()
	router.GET("/:name", func(c *gin.Context) {
		name := c.Params.ByName("name")
		in <- name
		var resp struct {
			Title string `json:"greeter"`
			Value string
		}
		resp.Title = "golang"
		resp.Value = <-out
		c.JSON(200, resp)
	})
	router.Run(":8080")
	close(in)
	close(out)
	<-net.Wait()
}
