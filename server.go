package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"strconv"
)

type Server struct {
	port uint
}

func (p *Server) init() {
	p.port = 3000

	files := make(map[string]string)

	loadFilesRecursively(&files, "./src/")

	loadFile(&files, "./src/favicon.ico", ".png")

	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(rw, getFile(&files, html.EscapeString(r.URL.Path)))
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
