package main

import (
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"net/http"
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
			fmt.Println(cwd + f.Name())
			loadFileFromFullPath(files, cwd+f.Name())
		}
	}
}

func loadFileFromFullPath(files *map[string]string, fullPath string) {
	dat, err := ioutil.ReadFile(fullPath)
	check(err)
	(*files)[fullPath] = string(dat)

	// dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	// fmt.Print("Path: ", dir)
}

func loadFile(files *map[string]string, path string, ext string) {
	loadFileFromFullPath(files, "./src/"+path[1:]+"."+ext)

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

	loadFile(&files, "/index", "html")
	loadFile(&files, "/404", "html")
	loadFile(&files, "/favicon.ico", "png")

	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		// fmt.Fprintf(rw, "Hello, %q", html.EscapeString(r.URL.Path))
		fmt.Fprintf(rw, getFile(&files, html.EscapeString(r.URL.Path)))
		fmt.Println("Recivied request: ", html.EscapeString(r.URL.Path))
	})

	const portNumber int = 3000
	fmt.Println("Started server at: ", portNumber)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(portNumber), nil))
}
