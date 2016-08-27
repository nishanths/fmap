package main

import (
	"errors"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/template"
)

var flags = struct {
	Package string
	Var     string
}{}

var stdout = log.New(os.Stdout, "", 0)
var stderr = log.New(os.Stderr, "", 0)

const helpString = `fmap generates a go source file containing a map[string][]byte that
corresponds to the specified directory tree. The keys are the paths
of files and the values are the contents of the file at that path.

The generated file is printed to stdout. Empty directories are
ignored, and symlinks are not followed.

usage:
  fmap [flags] /path/to/dir

flags:
  -package  package name to use in generated file (default: "main")
  -var      variable name of the map (default: "files")`

const fileTmpl = `package << .Package >>

var << .Var >> = map[string][]byte{
	<<- range $k, $v := .M >>
	"<< $k >>": << $v | printf "%#v" | trimPrefix >>,
	<<- end >>
}`

func main() {
	flag.StringVar(&flags.Package, "package", "main", "package name")
	flag.StringVar(&flags.Var, "var", "files", "variable name of map")
	flag.Usage = func() {
		stderr.Println(helpString)
		os.Exit(2)
	}
	flag.Parse()

	dirRoot := flag.Arg(0)
	if dirRoot == "" {
		stderr.Println(errors.New(`fmap: error: require path argument`))
		stderr.Println(helpString)
		os.Exit(2)
	}

	info, err := os.Stat(dirRoot)
	if err != nil {
		stderr.Println(err)
		os.Exit(1)
	}
	if !info.IsDir() {
		stderr.Println(errors.New(`fmap: error: "path" argument must be a directory"`))
		os.Exit(1)
	}

	tmpl, err := template.New("file").Funcs(template.FuncMap{
		"trimPrefix": func(s string) string { return strings.TrimPrefix(s, "[]byte") },
	}).Delims("<<", ">>").Parse(fileTmpl)
	if err != nil {
		stderr.Println(err)
		os.Exit(1)
	}

	m := make(map[string][]byte)
	out := fileContents(dirRoot)
	for f := range out {
		if f.Err != nil {
			stderr.Println(f.Err)
			os.Exit(1)
		}
		absRoot, err := filepath.Abs(dirRoot)
		if err != nil {
			stderr.Println(err)
			os.Exit(1)
		}
		absP, err := filepath.Abs(f.Path)
		if err != nil {
			stderr.Println(err)
			os.Exit(1)
		}
		s, err := filepath.Rel(absRoot, absP)
		if err != nil {
			stderr.Println(err)
			os.Exit(1)
		}
		k := filepath.ToSlash(s)
		m[k] = f.Content
	}

	tmpl.Execute(os.Stdout, struct {
		Package string
		Var     string
		M       map[string][]byte
	}{
		flags.Package,
		flags.Var,
		m,
	})
}

type File struct {
	Path    string
	Content []byte
	Err     error
}

func fileContents(root string) <-chan File {
	out := make(chan File)

	go func() {
		wg := sync.WaitGroup{}

		// Ignore error since all errors are sent on the channel anyway.
		_ = filepath.Walk(root, func(p string, info os.FileInfo, err error) error {
			if err != nil {
				out <- File{Err: err}
				return err
			}
			if info.IsDir() {
				return nil
			}
			wg.Add(1)
			go func() {
				defer wg.Done()
				b, err := ioutil.ReadFile(p)
				out <- File{p, b, err}
			}()
			return nil
		})

		go func() {
			wg.Wait()
			close(out)
		}()
	}()

	return out
}
