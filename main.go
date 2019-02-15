package main

import (
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func loadFilesRecursively(files *map[string]string, cwd string) {
	fileList, err := ioutil.ReadDir(cwd)
	check(err)

	for _, f := range fileList {
		if f.IsDir() {
			loadFilesRecursively(files, cwd+f.Name()+"/")
		} else {
			basename := f.Name()[0 : len(f.Name())-len(filepath.Ext(f.Name()))]
			fmt.Println(basename)
			fmt.Println(cwd + filepath.Base(f.Name()))
			loadFile(files, cwd+basename, filepath.Ext(f.Name()))
		}
	}
}

func loadFile(files *map[string]string, path string, ext string) {
	dat, err := ioutil.ReadFile(path + ext)
	check(err)
	// hack for now
	(*files)[path[5:]] = string(dat)
}

func getFile(files *map[string]string, path string) string {
	if val, ok := (*files)[path]; ok {
		return val
	}
	return (*files)["/404"]
}

func main() {
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
