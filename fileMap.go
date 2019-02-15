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

func (p *FileMap) load(path string, ext string) (string, error) {
	dat, err := ioutil.ReadFile(path + ext)
	content := string(dat)

	if err != nil {
		fmt.Printf(err.Error())
		err = StringError{s: "ERROR: Can't load file \"" + path + ext + "\"!"}
		return "", err
	}

	// hack for now
	p.files[path[len(p.skipPath):]] = content

	return content, nil
}

func (p FileMap) get(path string) (string, error) {
	val, ok := p.files[path]

	if ok {
		return val, nil
	} else {
		err := StringError{s: "WARNING: No file named \"" +
			path + "\" was loaded."}
		return "", err
	}
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
			p.load(cwd+baseName, ext) // note - no error handling
		}
	}
}
