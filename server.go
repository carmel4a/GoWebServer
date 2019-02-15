package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"strconv"
)

type Server struct {
	port    uint
	fileMap FileMap
}

func (p *Server) init() {
	p.port = 3000

	p.fileMap.init("./src")
	p.fileMap.loadFilesRecursively("./src/")
	p.fileMap.load("./src/favicon.ico", ".png")

	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		filecontent, ok := p.fileMap.get(html.EscapeString(r.URL.Path))
		if ok {
			fmt.Fprintf(rw, filecontent)
		} else {
			filecontent, _ := p.fileMap.get("/404")
			fmt.Fprintf(rw, filecontent)
		}
		fmt.Println("Recivied request: ", html.EscapeString(r.URL.Path))
	})

	const portNumber int = 3000
	fmt.Println("Started server at: ", portNumber)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(portNumber), nil))
}

func (p Server) getPort() uint {
	return p.port
}

func (p Server) setPort() uint {
	return p.port
}
