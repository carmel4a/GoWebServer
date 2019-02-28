package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"

	_ "github.com/mattn/go-sqlite3"
)

type Server struct {
	port        uint
	fileMap     FileMap
	httpHandler HTTPHandler
	router      *chi.Mux
	database    DatabaseHandler
	srcDir      string
	websiteDir  string
}

func (p *Server) init() {
	p.port = 3000
	p.srcDir = "src"
	p.websiteDir = "e-journal-frontend"

	p.database.init(p)
	p.fileMap.init(p.getWebsiteDir(), []string{".scss", ".git"})
	p.router = chi.NewRouter()
	p.httpHandler.init(p)

	p.loadFiles()
	p.run(3000)
}

func (p Server) getWebsiteDir() string {
	return p.srcDir + "/" + p.websiteDir + "/"
}

func (p *Server) loadFiles() {
	p.fileMap.loadFilesRecursively(p.getWebsiteDir())
}

func (p *Server) run(port int) {
	fmt.Println("Started server at: ", port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), p.router))
}
