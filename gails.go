package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	"github.com/shxsun/flags"
	"github.com/shxsun/go-sh"
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

func ignore(info os.FileInfo) bool {
	if info.IsDir() {
		if info.Name() != "." && info.Name() != ".." &&
			strings.HasPrefix(info.Name(), ".") { // ignore hidden dir
			return true
		}
	} else {
		return strings.HasPrefix(info.Name(), ".")
	}
	return false
}

func pathWalk(path string, depth int) (files []string, err error) {
	files = make([]string, 0)
	baseNumSeps := strings.Count(path, string(os.PathSeparator))
	err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			pathDepth := strings.Count(path, string(os.PathSeparator)) - baseNumSeps
			if pathDepth > depth {
				return filepath.SkipDir
			}
			if ignore(info) {
				return filepath.SkipDir
			}
		} else if info.Mode().IsRegular() && !ignore(info) {
			files = append(files, path)
			//if matched, _ := regexp.Match(mycnf.Include, []byte(info.Name())); matched { //}
		}
		return nil
	})
	return
}

var args []string

func init() {
	var err error
	args, err = flags.Parse(mycnf)
	if err != nil {
		os.Exit(1)
	}
}

func main() {
	// prepare arguments
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
	tmpdir, err := ioutil.TempDir("./", "tmp.gails.")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tmpdir)
	log.Println("use template", mycnf.Template)
	sh.Command("git", "clone", "https://github/"+mycnf.Template, tmpdir).Run()

	files, err := pathWalk("./", 1)
	if err != nil {
		log.Fatal(err)
	}

	// format code
	session := sh.NewSession()
	session.ShowCMD = true
	for _, file := range files {
		if strings.HasSuffix(file, ".go") {
			session.Command("go", "fmt", file).Run()
		}
	}
}

func render(src, dst string, v interface{}) {
	funcMap := template.FuncMap{
		"title": strings.Title,
	}
	t, err := template.New("test").Funcs(funcMap).ParseFiles(src)
	if err != nil {
		log.Fatal(err)
	}
	t.Execute(os.Stdout, v)
}

func deleteXXXX(filename string) (err error) {
	xxxx := regexp.MustCompile(`.*\s+XXXX.*\n`)
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	fi, _ := os.Stat(filename)
	out := xxxx.ReplaceAll(content, []byte(""))
	err = ioutil.WriteFile(filename, out, fi.Mode())
	return
}
