package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"github.com/shxsun/goyaml"

	"github.com/shxsun/flags"
	"github.com/shxsun/go-sh"
)

type Option struct {
	Template string `short:"t" long:"template" default:"shxsun/gails-default" description:"use which template(templale should on github)"`
	Path     string `short:"p" long:"path" default:"." description:"code generate path"`
	Depth    int    `short:"d" long:"depth" default:"4" description:"code depth to copy"`
}

var (
	mycnf   = &Option{}
	args    []string
	funcMap template.FuncMap
	gyml    *GenYaml
)

/*
// ORMName, Type, Name
{{ range . }}v.{{.Name}} = this.Get{{.Type|title}}("{{.ORMName}}")
{{ end }}

type Book struct {
	ID int64 json
	{{ range . }}{{.Name}}  {{.Type}} xorm:"{{.ORMName}}"
	{{ end }}
}
*/

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

var Usage = `Usage:
	gogen repo [table] [col:type, col:type ...]

example:
	gogen rest book name:string
	gogen bootstrap

	// chone from spec git repo
	gogen github.com/gogenerate/rest book name:string
`

// modify global args[] and mycnf
func argParse() (vs map[string]interface{}) {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	var err error
	if args, err = flags.Parse(mycnf); err != nil {
		os.Exit(1)
	}
	if len(args) < 1 {
		fmt.Println(Usage)
		os.Exit(1)
	}
	repo := args[0]
	switch {
	case regexp.MustCompile(`^[\w-]+$`).MatchString(repo):
		repo = "github.com/gogenerate/" + repo
	case regexp.MustCompile(`^[\w-]+/[\w-]+$`).MatchString(repo):
		repo = "github.com/" + repo
	}
	mycnf.Template = repo

	args = args[1:]
	// table argument parse
	patten := regexp.MustCompile(`^(\w+):(string|int)$`)
	vs = make(map[string]interface{}, 0)
	cols := make([]Column, 0)
	if len(args) > 0 {
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
	}
	cwd, _ := os.Getwd()
	tgtwd := filepath.Clean(filepath.Join(cwd, mycnf.Path))
	vs["PWD"] = filepath.Clean(filepath.Join(tgtwd, mycnf.Path))
	vs["PkgPath"] = tgtwd[len(os.Getenv("GOPATH")+"/src/"):]
	vs["AppName"] = filepath.Base(tgtwd)
	return
}

func main() {
	// prepare arguments
	vs := argParse()

	// render template
	tmpdir, err := ioutil.TempDir("./", "tmp.gails.")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tmpdir)
	log.Println("use template", mycnf.Template)
	sh.Command("git", "clone", "https://"+mycnf.Template, tmpdir).Run()

	files, err := pathWalk(tmpdir, mycnf.Depth)
	if err != nil {
		log.Fatal(err)
	}

	// newName func, eg: user.go -- rename -> book.go
	pjoin := func(p string) string { return filepath.Join(tmpdir, p) }
	gyml = readGailsYml(pjoin(".gails.yml"))
	newName := func(file string) string {
		t, ok := gyml.FileRename[file]
		if !ok {
			return filepath.Join(mycnf.Path, file)
		}
		s, err := renderString(t, vs)
		if err != nil {
			log.Fatal(err)
		}
		return filepath.Join(mycnf.Path, s)
	}

	// format code
	for _, src := range files {
		orig := src[len(tmpdir)+1:]
		dst := newName(orig)
		if _, exists := fileExists(dst); !exists {
			dstDir := filepath.Dir(dst)
			os.MkdirAll(dstDir, 0755)
			fmt.Println("git://"+orig, "-->", dst)
			if err = renderFile(dst, src, vs); err != nil {
				log.Println(err)
				return
			}
			// format code
			if strings.HasSuffix(dst, ".go") {
				exec.Command("go", "fmt", dst).Run()
			}
		}
	}
	fmt.Println("---------- template notice -------------\n" + gyml.Notice)
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

type GenYaml struct {
	FileRename   map[string]string `yaml:"file-rename"`
	StringRename map[string]string `yaml:"string-rename"`
	Notice       string            `yaml:"notice"`
}

func readGailsYml(file string) (gyml *GenYaml) {
	var err error
	defer func() {
		gyml = new(GenYaml)
		if err != nil {
			gyml.FileRename = make(map[string]string)
			gyml.StringRename = make(map[string]string)
		}
	}()
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}
	err = goyaml.Unmarshal(data, gyml)
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

// TODO, change remove prefix
func renderFile(dst string, src string, v interface{}) (err error) {
	trimLeft, trimRight := "<<X", "X>>"
	_ = trimLeft
	_ = trimRight

	xxxx := regexp.MustCompile(`[^\n]*XXXX[^\n]*\n`)
	nnnn := regexp.MustCompile(`.*NNNN`)
	s, err := ioutil.ReadFile(src)
	if err != nil {
		return
	}
	out := xxxx.ReplaceAll(s, []byte(""))
	out = nnnn.ReplaceAll(out, []byte(""))

	t, err := template.New("rdfile").Funcs(funcMap).Parse(string(out))
	if err != nil {
		return
	}
	buf := bytes.NewBuffer(nil)
	err = t.Execute(buf, v)
	if err != nil {
		log.Println(err)
		return
	}
	out = buf.Bytes()
	for key, val := range gyml.StringRename {
		nkey, err := renderString(val, v)
		if err != nil {
			return err
		}
		out = bytes.Replace(out, []byte(key), []byte(nkey), -1)
	}

	fi, _ := os.Stat(src)
	err = ioutil.WriteFile(dst, out, fi.Mode())
	return
}
