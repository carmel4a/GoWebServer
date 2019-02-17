package main

import (
	"fmt"
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
	router.Route("/register", p.registerRoute)
	router.Route("/resources", p.resourcesRoute)
}

func (p *HTTPHandler) indexRoute(router chi.Router) {
	router.Get("/", p.getIndex)
	router.Get("/index", p.getIndex)
}

func (p *HTTPHandler) loginRoute(router chi.Router) {
	router.Get("/", p.getLogin)
	router.Post("/", p.postLogin)
}

func (p *HTTPHandler) resourcesRoute(router chi.Router) {
	router.Get("/javascript/*", p.getJS)
	router.Get("/css/*", p.getCSS)
}

func (p *HTTPHandler) registerRoute(router chi.Router) {
	router.Get("/", p.getRegister)
	router.Post("/", p.postRegister)
}

func (p *HTTPHandler) getIndex(responseWriter http.ResponseWriter,
	request *http.Request) {

	createContent(responseWriter, "./src/e-journal-frontend", []string{
		"/partials/baseof",
		"/content/index"}, EmptyPage{})
}

func (p *HTTPHandler) getLogin(responseWriter http.ResponseWriter,
	request *http.Request) {

	responseWriter.WriteHeader(http.StatusOK)
	createContent(responseWriter, "./src/e-journal-frontend", []string{
		"/partials/baseof",
		"/content/login"}, LoginPage{}.get())
}

func (p *HTTPHandler) postLogin(responseWriter http.ResponseWriter,
	request *http.Request) {

	request.ParseForm()

	/*
		page := Page{Name: request.FormValue("firstname")}
		createContent(responseWriter, "./src", []string{
			"/partials/baseof",
			"/content/greet"}, page)
	*/
	// fmt.Fprintf(responseWriter, val)
}

func (p *HTTPHandler) getRegister(responseWriter http.ResponseWriter,
	request *http.Request) {

	responseWriter.WriteHeader(http.StatusOK)
	createContent(responseWriter, "./src/e-journal-frontend", []string{
		"/partials/baseof",
		"/content/register"}, RegisterPage{}.get())
}

func (p *HTTPHandler) postRegister(responseWriter http.ResponseWriter,
	request *http.Request) {

	request.ParseForm()

	/*
		page := Page{Name: request.FormValue("firstname")}
		createContent(responseWriter, "./src", []string{
			"/partials/baseof",
			"/content/greet"}, page)
	*/
	// fmt.Fprintf(responseWriter, val)
}

func (p *HTTPHandler) getCSS(responseWriter http.ResponseWriter,
	request *http.Request) {

	base, _ := getBaseAndExt(request.URL.EscapedPath())

	val, err := p.server.fileMap.get(base)
	if err != nil {
		fmt.Println(err.Error())
	}
	responseWriter.Header().Add("Content-Type", "text/css")
	responseWriter.WriteHeader(http.StatusOK)
	responseWriter.Write([]byte(val))
}

func (p *HTTPHandler) getJS(responseWriter http.ResponseWriter,
	request *http.Request) {

	base, _ := getBaseAndExt(request.URL.EscapedPath())

	val, err := p.server.fileMap.get(base)
	if err != nil {
		fmt.Println(err.Error())
	}
	responseWriter.Header().Add("Content-Type", "text/javascript")
	responseWriter.WriteHeader(http.StatusOK)
	responseWriter.Write([]byte(val))
}
