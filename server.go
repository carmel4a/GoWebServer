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
}

func (p *Server) init() {
	p.port = 3000

	p.database.init()

	p.fileMap.init("./src/e-journal-frontend")
	p.router = chi.NewRouter()

	p.httpHandler.init(p)

	p.loadFiles()
	p.run(3000)
}

func (p Server) getPort() uint {
	return p.port
}

func (p Server) setPort() uint {
	return p.port
}

func (p *Server) getFileMap() *FileMap {
	return &p.fileMap
}

func (p *Server) loadFiles() {
	filesToLoad := []string{
		"./src/e-journal-frontend/favicon.ico.png",
	}

	p.fileMap.loadFilesRecursively("./src/e-journal-frontend/")
	for _, file := range filesToLoad {
		p.fileMap.loadFromFullPath(file)
	}
}

func (p *Server) run(port int) {
	fmt.Println("Started server at: ", port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), p.router))
}
