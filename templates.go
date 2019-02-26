package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
)

type PageTemplateSetup struct {
	baseDir          string
	templateFileList []string
	content          interface{}
}

func getPartials(fm FileMap) []string {
	var result []string
	for key := range fm.files {
		if len(key) < 10 {
			continue
		}
		if key[:8] == "partials" {
			result = append(result, fm.skipPath+key+".html")
		}
	}
	return result
}

func getTemplateFileList(fm FileMap, page PageTemplateSetup) []string {
	fileNames := getPartials(fm)
	for _, name := range page.templateFileList {
		fileNames = append(fileNames, page.baseDir+name+".html")
	}
	return fileNames
}

func createTemplate(absoluteTemplatePathes []string) (*template.Template, error) {
	return template.ParseFiles(absoluteTemplatePathes...)
}

func createPageFromTemplate(wr io.Writer, page PageTemplateSetup, fm FileMap) {
	t, err := createTemplate(getTemplateFileList(fm, page))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	buf := new(bytes.Buffer)
	err = t.ExecuteTemplate(buf, "baseof", page.content)
	if err != nil {
		println(err.Error())
	}
	fmt.Println("Temlpate: " + buf.String())
	wr.Write(buf.Bytes())
}
