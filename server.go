package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

type Server struct {
	port        uint
	fileMap     FileMap
	httpHandler HTTPHandler
	router      *chi.Mux
}

func (p *Server) init() {
	p.port = 3000

	p.fileMap.init("./src")
	p.router = chi.NewRouter()

	p.httpHandler.init(p)

	filesToLoad := []string{
		"./src/favicon.ico.png",
	}

	p.fileMap.loadFilesRecursively("./src/")
	for _, file := range filesToLoad {
		p.fileMap.loadFromFullPath(file)
	}

	const portNumber int = 3000
	fmt.Println("Started server at: ", portNumber)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(portNumber), p.router))
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
