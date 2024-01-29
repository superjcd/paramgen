package main

import (
	"bytes"
	"fmt"
	"go/format"
	"html/template"
	"log"
	"os"
	"regexp"
	"strings"
)

func genCode(file string, scm map[string]structComponets) {
	fileName := strings.TrimRight(file, ".go")
	src := fmt.Sprintf("package %s", pkgName)

	for _, v := range scm {
		forms, err := genForm(v)
		if err != nil {
			log.Fatal(err)
		}
		src += "\n" + forms

		jsons, err := genJson(v)
		if err != nil {
			log.Fatal(err)
		}
		src += "\n" + jsons
	}
	finalSrc, err := format.Source([]byte(src))
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println(string(finalSrc))
	outFile := fmt.Sprintf("%s_param.go", fileName)
	os.WriteFile(outFile, finalSrc, 0644)
}

func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

var (
	matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap   = regexp.MustCompile("([a-z0-9])([A-Z])")
	funcMap       = template.FuncMap{
		"ToLower": strings.ToLower,
		"ToSnake": ToSnakeCase,
	}
)

func genForm(sc structComponets) (string, error) {
	tempString := `
	type {{.Name}}Form struct {
		 {{range $i, $e := .FieldNames}} 
		    {{$e}}	{{index $.FieldTypes $i}}` + "   `form:\"{{$e|ToSnake}}\"`" + `{{end}}` + `
			}`

	temp, err := template.New("controller").Funcs(funcMap).Parse(tempString)
	if err != nil {
		return "", err
	}
	buff := bytes.NewBufferString("")

	err = temp.Execute(buff, sc)

	if err != nil {
		return "", err
	}
	return buff.String(), nil
}

func genJson(sc structComponets) (string, error) {
	tempString := `
	type {{.Name}}Json struct {
		 {{range $i, $e := .FieldNames}} 
		    {{$e}}	{{index $.FieldTypes $i}}` + "   `json:\"{{$e|ToSnake}}\"`" + `{{end}}` + `
			}`

	temp, err := template.New("controller").Funcs(funcMap).Parse(tempString)
	if err != nil {
		return "", err
	}
	buff := bytes.NewBufferString("")

	err = temp.Execute(buff, sc)

	if err != nil {
		return "", err
	}
	return buff.String(), nil

}
