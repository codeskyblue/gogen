package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"launchpad.net/goyaml"

	"github.com/shxsun/flags"
	"github.com/shxsun/go-sh"
)

type Option struct {
	Template string `short:"t" long:"template" default:"shxsun/gails-default" description:"use which template(templale should on github)"`
}

var (
	mycnf   = &Option{}
	args    []string
	funcMap template.FuncMap
)

// ORMName, Type, Name
/*
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

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	var err error
	args, err = flags.Parse(mycnf)
	if len(args) < 2 {
		program := filepath.Base(os.Args[0])
		fmt.Printf(`Usage:
	%s <table> <col:type> [col:type ...]
Example:
	%s book name:string
`, program, program)
		os.Exit(1)
	}
	if err != nil {
		os.Exit(1)
	}
}

func main() {
	// prepare arguments
	patten := regexp.MustCompile(`^(\w+):(string|int)$`)
	vs := make(map[string]interface{}, 0)
	cols := make([]Column, 0)
	vs["Table"] = args[0]
	for _, s := range args[1:] {
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
	//sh.Command("git", "clone", "https://github.com/"+mycnf.Template, tmpdir).Run()
	sh.Command("git", "clone", "/Users/skyblue/goproj/src/github.com/shxsun/gails-default", tmpdir).Run()

	files, err := pathWalk(tmpdir, 1)
	if err != nil {
		log.Fatal(err)
	}

	// newName func, eg: user.go -- rename -> book.go
	pjoin := func(p string) string { return filepath.Join(tmpdir, p) }
	filemap := readGailsYml(pjoin(".gails.yml"))
	newName := func(file string) string {
		t, ok := filemap[file]
		if !ok {
			return file
		}
		s, err := renderString(t, vs)
		if err != nil {
			log.Fatal(err)
		}
		return s
	}

	// format code
	session := sh.NewSession()
	session.ShowCMD = true
	for _, src := range files {
		//fmt.Println(file)
		orig := src[len(tmpdir)+1:]
		dst := newName(orig)
		if _, exists := fileExists(dst); !exists {
			fmt.Println("git://"+orig, "-->", dst)
			if err = renderFile(dst, src, vs); err != nil {
				sh.Command("cp", "-v", src, dst)
			}
			// format code
			if strings.HasSuffix(dst, ".go") {
				//session.Command("go", "fmt", file).Run()
			}
		}
	}
}

func init() {
	funcMap = template.FuncMap{
		"title": strings.Title,
	}
}

func fileExists(file string) (os.FileInfo, bool) {
	fi, err := os.Stat(file)
	return fi, err == nil
}
func readGailsYml(file string) (m map[string]string) {
	m = make(map[string]string)
	var err error
	defer func() {
		if err != nil {
			m[".gitignore"] = ".gitignore"
		}
	}()
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}
	err = goyaml.Unmarshal(data, m)
	return
}

func renderString(tmplstr string, v interface{}) (out string, err error) {
	mytmpl := template.New("rdstr").Funcs(funcMap)
	t, err := mytmpl.Parse(tmplstr)
	if err != nil {
		return
	}
	buf := bytes.NewBuffer(nil)
	err = t.Execute(buf, v)
	return string(buf.Bytes()), err
}

func renderFile(dst string, src string, v interface{}) (err error) {
	mytmpl := template.New("rdfile").Funcs(funcMap)
	t, err := mytmpl.ParseFiles(src)
	if err != nil {
		log.Fatal(err)
	}
	return t.Execute(os.Stdout, v)
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
