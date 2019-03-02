package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type EmptyPage struct{}

// HTTPHandler is HTTP entry point.
// Gets, posts etc.
type HTTPHandler struct {
	// server - server reference.
	server *Server
}

func (p *HTTPHandler) init(server *Server) {
	p.server = server
	var router = p.server.router

	router.Use(middleware.Timeout(60 * time.Second))
	router.Use(middleware.Logger)
	router.Use(isLogged)

	router.With(shouldBeLogged).Route("/", p.indexRoute)
	router.With(redirectToIndexIfLogged).Route("/login", p.loginRoute)
	router.With(redirectToIndexIfLogged).Route("/register", p.registerRoute)
	router.Route("/resources", p.resourcesRoute)
}

func (p *HTTPHandler) createNormalPage(responseWriter http.ResponseWriter, page interface{}, templatePathes []string) {
	responseWriter.WriteHeader(http.StatusOK)
	createPageFromTemplate(responseWriter, PageTemplateSetup{
		baseDir:          p.server.getWebsiteDir(),
		templateFileList: templatePathes,
		content:          page},
		p.server.fileMap)
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

	responseWriter.Write([]byte("Siema!"))
}

func (p *HTTPHandler) getFavIcon(responseWriter http.ResponseWriter,
	request *http.Request) {

	responseWriter.WriteHeader(http.StatusOK)
	responseWriter.Write([]byte(p.server.fileMap.files["favicon.ico.png"]))
}

func (p *HTTPHandler) getNotFund(responseWriter http.ResponseWriter,
	request *http.Request) {

	p.createNormalPage(responseWriter, Page{}.get(), []string{
		"partials/baseof",
		"content/404"})
}

// DEBUG/DEV method
func (p *HTTPHandler) getDefault(responseWriter http.ResponseWriter,
	request *http.Request) {
	responseWriter.WriteHeader(http.StatusOK)
	url := request.RequestURI[1:]
	createPageFromTemplate(responseWriter, PageTemplateSetup{
		baseDir: p.server.getWebsiteDir(),
		templateFileList: []string{
			"partials/baseof",
			"content/" + url},
		content: getDefaultPage(url).get()},
		p.server.fileMap)
}

func (p *HTTPHandler) getLogin(responseWriter http.ResponseWriter,
	request *http.Request) {

	p.createNormalPage(responseWriter, LoginPage{}.get(), []string{
		"partials/baseof",
		"content/login"})
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
		p.createNormalPage(responseWriter, Page{}.get(), []string{
			"partials/baseof",
			"content/index"})
	} else {
		p.createNormalPage(responseWriter, LoginPage{}.get(), []string{
			"partials/baseof",
			"content/login"})
	}
}

func (p *HTTPHandler) getRegister(responseWriter http.ResponseWriter,
	request *http.Request) {

	p.createNormalPage(responseWriter, RegisterPage{}.get(), []string{
		"partials/baseof",
		"content/register"})
}

func (p *HTTPHandler) postRegister(responseWriter http.ResponseWriter,
	request *http.Request) {

	request.ParseForm()

	login := request.FormValue("login")
	pass := request.FormValue("password")
	email := request.FormValue("email")

	result := p.server.database.register(login, email, pass)
	if result == RegisterOK {
		redirectTo(responseWriter, request, "/login/", http.StatusFound)
	} else {
		redirectTo(responseWriter, request, "/register/", http.StatusFound)
	}
}

func (p *HTTPHandler) getCSS(responseWriter http.ResponseWriter,
	request *http.Request) {

	base, _ := getBaseAndExt(request.URL.EscapedPath())

	val, err := p.server.fileMap.get(base[1:])
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

	val, err := p.server.fileMap.get(base[1:])
	if err != nil {
		fmt.Println(err.Error())
	}
	responseWriter.Header().Add("Content-Type", "text/javascript")
	responseWriter.WriteHeader(http.StatusOK)
	responseWriter.Write([]byte(val))
}

func isLogged(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "isLogged", true)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func shouldBeLogged(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if ctx.Value("isLogged").(bool) {
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			redirectTo(w, r, "/login/", http.StatusFound)
		}
	})
}

func redirectToIndexIfLogged(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if ctx.Value("isLogged").(bool) {
			redirectTo(w, r, "/index", http.StatusFound)
		} else {
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}

func redirectTo(responseWriter http.ResponseWriter,
	request *http.Request, where string, HTTPStatusCode int) {
	http.Redirect(responseWriter, request, where, HTTPStatusCode)
}
