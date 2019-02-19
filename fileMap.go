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
	exclude := []string{".scss", ".git"}
	for _, val := range exclude {
		if ext == val {
			return "", nil
		}
	}
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

func (p *FileMap) loadFromFullPath(path string) (string, error) {
	base, ext := getBaseAndExt(path)
	return p.load(base, ext)
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

func (p *FileMap) loadFilesRecursively(cwd string) error {
	fileList, err := ioutil.ReadDir(cwd)
	if err != nil {
		err = StringError{s: "ERROR: Can't open \"" + cwd + "\" directory!"}
		return err
	}

	for _, f := range fileList {
		fileName := f.Name()
		if f.IsDir() {
			err := p.loadFilesRecursively(cwd + fileName + "/")
			if err != nil {
				fmt.Println(err.Error())
				return err
			}
		} else {
			baseName, ext := getBaseAndExt(fileName)

			_, err := p.load(cwd+baseName, ext)
			if err != nil {
				fmt.Println(err.Error())
				return err
			}

			fmt.Print("INFO: Loaded file: ")
			fmt.Println(cwd + filepath.Base(fileName))
		}
	}
	return nil
}
