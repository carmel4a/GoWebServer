package main

import (
	"fmt"
	"html"
	"net/http"
)

type HTTPHandler struct {
	server *Server
}

func (p HTTPHandler) init(server *Server) {
	p.server = server

	p.route("/", p.indexHandle)
}

func (p HTTPHandler) route(url string,
	function func(http.ResponseWriter, *http.Request)) {
	http.HandleFunc(url, function)
}

func (p HTTPHandler) indexHandle(responseWriter http.ResponseWriter,
	request *http.Request) {
	url := html.EscapeString(request.URL.Path)
	filecontent, err := p.server.getFileMap().get(url)

	if err == nil {
		fmt.Fprintf(responseWriter, filecontent)
	} else {
		filecontent, _ := p.server.getFileMap().get("/404")
		fmt.Fprintf(responseWriter, filecontent)
	}
}
