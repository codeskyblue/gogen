package main

import (
	"bytes"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"text/template"

	"github.com/shxsun/flags"
)

type Option struct {
	Template string `short:"t" long:"template" default:"shxsun/gails-default" description:"use which template(templale should on github)"`
}

var mycnf = &Option{}

// ORMName, Type, Name
var tmpl = `
		/* XXXX
		{{ range . }}v.{{.Name}} = this.Get{{.Type|title}}("{{.ORMName}}") 
		{{ end }}

		type Book struct {
			ID int64 json
			{{ range . }}{{.Name}}  {{.Type}} xorm:"{{.ORMName}}" 
			{{ end }}
		}
		XXXX */
		`

type Column struct {
	Name    string
	Type    string
	ORMName string
}

func main() {
	if len(os.Args) != 1 {
		shPath := filepath.Join(filepath.Dir(os.Args[0]), "clone.sh")
		err := syscall.Exec(shPath, os.Args, os.Environ())
		log.Fatal(err)
	}
	args, err := flags.Parse(mycnf)
	if err != nil {
		os.Exit(1)
	}
	patten := regexp.MustCompile(`^(\w+):(string|int)$`)
	vs := make(map[string]interface{}, 0)
	cols := make([]Column, 0)
	for _, s := range args {
		vs := patten.FindStringSubmatch(s)
		if vs == nil {
			log.Fatalf("invalid format: %s", strconv.Quote(s))
		}
		c := Column{}
		c.ORMName, c.Type = vs[1], vs[2]
		c.Name = strings.Title(c.ORMName)
		cols = append(cols, c)
	}
	vs["Cols"] = cols
	vs["PWD"], _ = os.Getwd()

	// render template
	funcMap := template.FuncMap{
		"title": strings.Title,
	}
	buf := bytes.NewBuffer(nil)
	io.Copy(buf, os.Stdin)
	t, err := template.New("test").Funcs(funcMap).Parse(string(buf.Bytes()))
	if err != nil {
		log.Fatal(err)
	}
	t.Execute(os.Stdout, cols)
}
