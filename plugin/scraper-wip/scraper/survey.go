package scraper

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-mangos/mangos"
	"github.com/go-mangos/mangos/protocol/rep"
	"github.com/go-mangos/mangos/protocol/respondent"
	"github.com/go-mangos/mangos/protocol/surveyor"
	"github.com/go-mangos/mangos/transport/ipc"
	"github.com/go-mangos/mangos/transport/tcp"
	"github.com/go-mangos/mangos/transport/ws"
)

func die(format string, v ...interface{}) {
	fmt.Fprintln(os.Stderr, fmt.Sprintf(format, v...))
	os.Exit(1)
}

func date() string {
	return time.Now().Format(time.ANSIC)
}

func serverNN(url string) {
	var sock mangos.Socket
	var err error
	var msg []byte
	if sock, err = surveyor.NewSocket(); err != nil {
		die("can't get new surveyor socket: %s", err)
	}
	sock.AddTransport(ipc.NewTransport())
	sock.AddTransport(tcp.NewTransport())
	if err = sock.Listen(url); err != nil {
		die("can't listen on surveyor socket: %s", err.Error())
	}
	err = sock.SetOption(mangos.OptionSurveyTime, time.Second/2)
	if err != nil {
		die("SetOption(): %s", err.Error())
	}
	for {
		time.Sleep(time.Second)
		fmt.Println("SERVER: SENDING DATE SURVEY REQUEST")
		if err = sock.Send([]byte("DATE")); err != nil {
			die("Failed sending survey: %s", err.Error())
		}
		for {
			if msg, err = sock.Recv(); err != nil {
				break
			}
			fmt.Printf("SERVER: RECEIVED \"%s\" SURVEY RESPONSE\n",
				string(msg))
		}
		fmt.Println("SERVER: SURVEY OVER")
	}
}

func clientNN(url string, name string) {
	var sock mangos.Socket
	var err error
	var msg []byte

	if sock, err = respondent.NewSocket(); err != nil {
		die("can't get new respondent socket: %s", err.Error())
	}
	sock.AddTransport(ipc.NewTransport())
	sock.AddTransport(tcp.NewTransport())
	if err = sock.Dial(url); err != nil {
		die("can't dial on respondent socket: %s", err.Error())
	}
	for {
		if msg, err = sock.Recv(); err != nil {
			die("Cannot recv: %s", err.Error())
		}
		fmt.Printf("CLIENT(%s): RECEIVED \"%s\" SURVEY REQUEST\n",
			name, string(msg))

		d := date()
		fmt.Printf("CLIENT(%s): SENDING DATE SURVEY RESPONSE\n", name)
		if err = sock.Send([]byte(d)); err != nil {
			die("Cannot send: %s", err.Error())
		}
	}
}

func testSurvey() {
	if len(os.Args) > 2 && os.Args[1] == "server" {
		serverNN(os.Args[2])
		os.Exit(0)
	}
	if len(os.Args) > 3 && os.Args[1] == "client" {
		clientNN(os.Args[2], os.Args[3])
		os.Exit(0)
	}
	fmt.Fprintf(os.Stderr, "Usage: survey server|client <URL> <ARG>\n")
	os.Exit(1)
}

// reqHandler just spins on the socket and reads messages.  It replies
// with "REPLY <time>".  Not very interesting...

func reqHandler(sock mangos.Socket) {
	count := 0
	for {
		// don't care about the content of received message
		_, e := sock.Recv()
		if e != nil {
			die("Cannot get request: %v", e)
		}
		reply := fmt.Sprintf("REPLY #%d %s", count, time.Now().String())
		if e := sock.Send([]byte(reply)); e != nil {
			die("Cannot send reply: %v", e)
		}
		count++
	}
}

func addReqHandler(mux *http.ServeMux, port int) {
	sock, _ := rep.NewSocket()

	sock.AddTransport(ws.NewTransport())

	url := fmt.Sprintf("ws://127.0.0.1:%d/req", port)

	if l, e := sock.NewListener(url, nil); e != nil {
		die("bad listener: %v", e)
	} else if h, e := l.GetOption(ws.OptionWebSocketHandler); e != nil {
		die("bad handler: %v", e)
	} else {
		mux.Handle("/req", h.(http.Handler))
		l.Listen()
	}
	go reqHandler(sock)
}
