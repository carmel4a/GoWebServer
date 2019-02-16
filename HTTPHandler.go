package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type EmptyPage struct{}

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
	router.Get("/", p.getIndex)
	router.Get("/index", p.getIndex)
}

func (p *HTTPHandler) loginRoute(router chi.Router) {
	router.Get("/", p.getLogin)
	router.Post("/", p.postLogin)
}

func (p *HTTPHandler) getIndex(responseWriter http.ResponseWriter,
	request *http.Request) {

	createContent(responseWriter, "./src", []string{
		"/partials/baseof",
		"/content/index"}, EmptyPage{})
}

func (p *HTTPHandler) getLogin(responseWriter http.ResponseWriter,
	request *http.Request) {

	responseWriter.WriteHeader(http.StatusOK)
	createContent(responseWriter, "./src", []string{
		"/partials/baseof",
		"/content/login"}, EmptyPage{})
}

func (p *HTTPHandler) postLogin(responseWriter http.ResponseWriter,
	request *http.Request) {

	// val, _ := p.server.getFileMap().get("/greet")

	request.ParseForm()

	type Page struct {
		Name string
	}

	page := Page{Name: request.FormValue("firstname")}
	createContent(responseWriter, "./src", []string{
		"/partials/baseof",
		"/content/greet"}, page)

	// fmt.Fprintf(responseWriter, val)
}
