package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
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
	server := Server{}
	server.init()
}
