package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"text/template"
)

const (
	path        = "../ast/"
	visitorTmpl = `
	// Code generated by visitor_generator. DO NOT EDIT.
	package ast

	//T is a generic type
	type Something generic.Type
	type Another generic.Type
	type Ctx generic.Type

	//SomethingAnotherVisitor is an AST visitor that returns Another for expressions and Something for everything elsd
	type SomethingAnotherVisitor interface {

		{{range $name := .Ifaces}}
		//Visit{{$name}} visits AST node {{$name}}
		Visit{{- $name}}(node {{$name}}, ctx Ctx) {{if or (eq $name "Expression") (eq $name "RHS")}} Another {{else}} Something {{end}}
		{{end}}

		{{range $i, $name := .Structs}}
		//Visit{{$name}} visits AST node {{$name}}
		Visit{{- $name}}(node {{$name}}, ctx Ctx) {{index $.Types $i}}
		{{end}}
	}
	`
	acceptableTmpl = `
	// Code generated by visitor_generator. DO NOT EDIT.
	package ast
	
	type SomethingAcceptor interface {
		AcceptSomething(v SomethingAnotherVisitor, ctx Ctx) Something
	}

	type AnotherAcceptor interface {
		AcceptAnother(v SomethingAnotherVisitor, ctx Ctx) Another
	}

	{{range $i, $name := .Structs}}
	{{$f := fst $name}}
	{{$t := index $.Types $i}}
	//Accept calls v.Visit{{$name}}({{$f}})
	func ({{$f}} {{$name}}) Accept{{$t}}(v SomethingAnotherVisitor, ctx Ctx) {{$t}} {
		return v.Visit{{$name}}({{$f}}, ctx)
	}
	{{end}}
	`
)

var (
	nodeMatch = regexp.MustCompile(`(var )?_ (.*) = &?([A-Z][a-zA-Z]*)`)
)

type namesList struct {
	Structs []string
	Ifaces  []string
	Types   []string
}

func getASTNodes() (structs []string, ifaces []string, types []string) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	ifacesUnique := make(map[string]struct{})

	for _, f := range files {
		if strings.Contains(f.Name(), "visitor") || strings.Contains(f.Name(), "acceptor") {
			continue
		}
		content, _ := ioutil.ReadFile(path + f.Name())
		matches := nodeMatch.FindAllStringSubmatch(string(content), -1)

		const iface = 2
		const nodeName = 3
		for _, match := range matches {
			structs = append(structs, match[nodeName])
			ifacesUnique[match[iface]] = struct{}{}
			if match[iface] == "Expression" || match[iface] == "RHS" {
				types = append(types, "Another")
			} else {
				types = append(types, "Something")
			}
		}
	}

	for k := range ifacesUnique {
		ifaces = append(ifaces, k)
	}

	return structs, ifaces, types
}

func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}

func templateToFile(filename string, tmpl *template.Template, data namesList) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	err = tmpl.Execute(file, data)
	if err != nil {
		return err
	}
	err = exec.Command("gofmt", "-w", filename).Run()
	if err != nil {
		return err
	}
	return exec.Command("../../lib/wacc-goimports", "-w", filename).Run()
}

func main() {
	visitorTemplate, err := template.New("").Parse(visitorTmpl)

	panicErr(err)

	names := namesList{}
	names.Structs, names.Ifaces, names.Types = getASTNodes()

	panicErr(templateToFile("../ast/visitor.go", visitorTemplate, names))

	acceptableTemplate, err := template.New("").Funcs(template.FuncMap{
		"fst": func(str string) string {
			return strings.ToLower(str[:1])
		},
	}).Parse(acceptableTmpl)

	panicErr(err)
	panicErr(templateToFile("../ast/acceptor.go", acceptableTemplate, names))

}
