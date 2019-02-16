package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
)

func createContent(wr io.Writer, baseDir string, contentNames []string, content interface{}) {
	var fileNames []string
	for _, fileName := range contentNames {
		fileNames = append(fileNames, baseDir+fileName+".html")
	}
	t, err := template.ParseFiles(fileNames...)
	if err != nil {
		fmt.Println(err.Error())
	}

	buf := new(bytes.Buffer)
	t.ExecuteTemplate(buf, "baseof", content)
	fmt.Println("Temlpate: " + buf.String())
	wr.Write(buf.Bytes())
}
