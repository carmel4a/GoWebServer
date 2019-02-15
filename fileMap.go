package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
)

type FileMap struct {
	files    map[string]string
	skipPath string // Used to omit part of loaded file path
}

func (p *FileMap) init(skipPath string) {
	p.files = make(map[string]string)
	p.skipPath = skipPath
}

func (p *FileMap) load(path string, ext string) {
	dat, err := ioutil.ReadFile(path + ext)
	check(err)

	// hack for now
	p.files[path[len(p.skipPath):]] = string(dat)
}

func (p FileMap) get(path string) (string, bool) {
	val, ok := p.files[path]
	return val, ok
}

func (p *FileMap) loadFilesRecursively(cwd string) {
	fileList, err := ioutil.ReadDir(cwd)
	check(err)

	for _, f := range fileList {
		fileName := f.Name()
		if f.IsDir() {
			p.loadFilesRecursively(cwd + fileName + "/")
		} else {
			ext := filepath.Ext(fileName)
			baseName := fileName[0 : len(fileName)-len(ext)]
			fmt.Println(cwd + filepath.Base(fileName))
			p.load(cwd+baseName, ext)
		}
	}
}
