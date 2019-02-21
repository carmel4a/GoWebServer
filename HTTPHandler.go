package main

import (
	"fmt"
	"net/http"
	"strings"
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
	router.Get("/favicon.ico", p.getFavIcon)
	router.NotFound(p.getNotFund)
	router.Get("/*", p.getDefault)
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
	// if known user
	/*
		createContent(responseWriter, "./src/e-journal-frontend", []string{
			"/partials/baseof",
			"/content/index"}, EmptyPage{})
	*/
	// else
	http.Redirect(responseWriter, request, "/login/", http.StatusFound)
}

func (p *HTTPHandler) getFavIcon(responseWriter http.ResponseWriter,
	request *http.Request) {

	responseWriter.WriteHeader(http.StatusOK)
	responseWriter.Write([]byte(p.server.fileMap.files["/favicon.ico.png"]))
}

func (p *HTTPHandler) getNotFund(responseWriter http.ResponseWriter,
	request *http.Request) {

	responseWriter.WriteHeader(http.StatusOK)
	createPageFromTemplate(responseWriter, PageTemplateSetup{
		baseDir: "./src/e-journal-frontend",
		templateFileList: []string{
			"/partials/baseof",
			"/content/404"},
		content: LoginPage{}.get()},
		p.server.fileMap)
}

// DEBUG/DEV method
func (p *HTTPHandler) getDefault(responseWriter http.ResponseWriter,
	request *http.Request) {
	responseWriter.WriteHeader(http.StatusOK)
	url := request.RequestURI[1:]
	createPageFromTemplate(responseWriter, PageTemplateSetup{
		baseDir: "./src/e-journal-frontend",
		templateFileList: []string{
			"/partials/baseof",
			"/content/" + url},
		content: getDefaultPage(url).get()},
		p.server.fileMap)
}

func (p *HTTPHandler) getLogin(responseWriter http.ResponseWriter,
	request *http.Request) {

	responseWriter.WriteHeader(http.StatusOK)
	createPageFromTemplate(responseWriter, PageTemplateSetup{
		baseDir: "./src/e-journal-frontend",
		templateFileList: []string{
			"/partials/baseof",
			"/content/login"},
		content: LoginPage{}.get()},
		p.server.fileMap)
}

func (p *HTTPHandler) postLogin(responseWriter http.ResponseWriter,
	request *http.Request) {

	request.ParseForm()
	login := request.FormValue("login")
	pass := request.FormValue("password")
	var lm LoginMethod
	if strings.Contains("@", login) {
		lm = EmailLoginMethod
	} else {
		lm = LoginLoginMethod
	}
	if p.server.database.login(login, pass, lm) {
		createPageFromTemplate(responseWriter,
			PageTemplateSetup{
				baseDir: "./src/e-journal-frontend",
				templateFileList: []string{
					"/partials/baseof",
					"/content/index"},
				content: RegisterPage{}.get()},
			p.server.fileMap)
	} else {
		createPageFromTemplate(responseWriter,
			PageTemplateSetup{
				baseDir: "./src/e-journal-frontend",
				templateFileList: []string{
					"/partials/baseof",
					"/content/login"},
				content: LoginPage{}.get()},
			p.server.fileMap)
	}
}

func (p *HTTPHandler) getRegister(responseWriter http.ResponseWriter,
	request *http.Request) {

	responseWriter.WriteHeader(http.StatusOK)
	createPageFromTemplate(responseWriter,
		PageTemplateSetup{
			baseDir: "./src/e-journal-frontend",
			templateFileList: []string{
				"/partials/baseof",
				"/content/register"},
			content: RegisterPage{}.get()},
		p.server.fileMap)
}

func (p *HTTPHandler) postRegister(responseWriter http.ResponseWriter,
	request *http.Request) {

	request.ParseForm()

	login := request.FormValue("login")
	pass := request.FormValue("password")
	email := request.FormValue("email")

	result := p.server.database.register(login, email, pass)
	if result == RegisterOK {
		http.Redirect(responseWriter, request, "/login/", http.StatusFound)
	} else {
		http.Redirect(responseWriter, request, "/register/", http.StatusFound)
	}
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
