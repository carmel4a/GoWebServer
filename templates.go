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

func getTemplateFileList(page PageTemplateSetup) []string {
	var fileNames []string
	for _, name := range page.templateFileList {
		fileNames = append(fileNames, page.baseDir+name+".html")
	}
	return fileNames
}

func createTemplate(absoluteTemplatePathes []string) (*template.Template, error) {
	return template.ParseFiles(absoluteTemplatePathes...)
}

func createPageFromTemplate(wr io.Writer, page PageTemplateSetup) {
	t, err := createTemplate(getTemplateFileList(page))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	buf := new(bytes.Buffer)
	t.ExecuteTemplate(buf, "baseof", page.content)
	fmt.Println("Temlpate: " + buf.String())
	wr.Write(buf.Bytes())
}
