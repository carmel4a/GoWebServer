package main

import (
	"fmt"
	"html"
	"html/template"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type HTTPHandler struct {
	server *Server
}

func (p *HTTPHandler) init(server *Server) {
	p.server = server
	var router = p.server.router

	router.Use(middleware.Timeout(60 * time.Second))
	router.Use(middleware.Logger)

	router.Route("/", p.indexRoute)
	router.Route("/login", p.loginRoute)
}

func (p *HTTPHandler) indexRoute(router chi.Router) {
	router.Get("/", p.indexHandle)
	router.Get("/index", p.indexHandle)
}

func (p *HTTPHandler) loginRoute(router chi.Router) {
	router.Get("/", p.getLogin)
	router.Post("/", p.postLogin)
}

func (p *HTTPHandler) indexHandle(responseWriter http.ResponseWriter,
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

func (p *HTTPHandler) getLogin(responseWriter http.ResponseWriter,
	request *http.Request) {

	val, _ := p.server.getFileMap().get("/login")
	responseWriter.WriteHeader(http.StatusOK)
	fmt.Fprintf(responseWriter, val)
	// responseWriter.Write([]byte(val))

	/*
		url := html.EscapeString(request.URL.Path)
		filecontent, err := p.server.getFileMap().get(url)

		if err == nil {
			fmt.Fprintf(responseWriter, filecontent)
		} else {
			filecontent, _ := p.server.getFileMap().get("/404")
			fmt.Fprintf(responseWriter, filecontent)
		}
	*/
}

func (p *HTTPHandler) postLogin(responseWriter http.ResponseWriter,
	request *http.Request) {

	// val, _ := p.server.getFileMap().get("/greet")

	request.ParseForm()

	type Page struct {
		Name string
	}

	page := Page{Name: request.FormValue("firstname")}

	t, _ := template.ParseFiles("./src/greet.html")
	responseWriter.WriteHeader(http.StatusOK)
	t.Execute(responseWriter, page)

	// fmt.Fprintf(responseWriter, val)
}
