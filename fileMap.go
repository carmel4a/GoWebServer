package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
)

// FileMap is structure that hold (relative) names of files with their content.
type FileMap struct {
	// files is a map where keys are relative pathes (wit skipped `skipPath`
	// part). Values are raw file content.
	files map[string]string
	// skipPath is used to omit part of loaded file path
	skipPath string
	// excludedExt stores extentions witch mustn't be loaded. Changing this in
	// runtime has no effect on already loaded files.
	excludedExt []string
}

// init creates FileMap. Allocates map.
func (p *FileMap) init(skipPath string, excludedExt []string) {
	p.files = make(map[string]string)
	p.skipPath = skipPath
	p.excludedExt = excludedExt
}

// load loads file from given `path` and `extention`.
// Returns content of file as string. Error is returned if load weren't
// successfull OR given extention should be skipped.
func (p *FileMap) load(path string, ext string) (string, error) {
	for _, val := range p.excludedExt {
		if ext == val {
			err := StringError{"INFO: Skipped loading: \"" + path + ext + "\"."}
			return "", err
		}
	}
	dat, err := ioutil.ReadFile(path + ext)

	if err != nil {
		fmt.Printf(err.Error())
		err = StringError{"ERROR: Can't load file \"" + path + ext + "\"!"}
		return "", err
	}

	content := string(dat)

	// XXX clear that
	p.files[path[len(p.skipPath):]] = content

	return content, nil
}

// loadFromFullPath loads file from full path (relative to executable).
// For details see `load`.
func (p *FileMap) loadFromFullPath(path string) (string, error) {
	base, ext := getBaseAndExt(path)
	return p.load(base, ext)
}

// get returns content of file with given path. Note that `skipPath` shouldn't
// be passed. Returns content of file. Error is returned if no file were loaded
// (or were skipped).
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

// loadFilesRecursively loads every file from `cwd`, and subdirectories.
// Returns error if `cwd` cant' be opened. Sub-errors about files and
// subdirectories are ignored.
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
				continue
			}
		} else {
			baseName, ext := getBaseAndExt(fileName)

			_, err := p.load(cwd+baseName, ext)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}

			fmt.Println("INFO: Loaded file: " + cwd + filepath.Base(fileName))
		}
	}
	return nil
}
