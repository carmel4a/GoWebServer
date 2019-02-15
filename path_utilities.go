package main

import "path/filepath"

func getBaseAndExt(path string) (string, string) {
	ext := filepath.Ext(path)
	baseName := path[0 : len(path)-len(ext)]
	return baseName, ext
}
