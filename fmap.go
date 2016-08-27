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
	Abs     bool
}{}

var stdout = log.New(os.Stdout, "", 0)
var stderr = log.New(os.Stderr, "", 0)

const helpString = `fmap generates a go source file containing a map[string][]byte
for the specified directory trees. The keys are the paths
of files and the values are the contents of the file at that path.

The generated file is printed to stdout. Empty directories are
ignored, and symlinks are not followed.

usage:
  fmap [flags] path/to/dir path/to/dir2 ...

flags:
  -package  package name to use in generated file (default: "main")
  -var      variable name of the map (default: "files")
  -abs      use absolute paths for map keys

example:
  fmap static/css static/js | gofmt > static_files.go`

const fileTmpl = `package << .Package >>

var << .Var >> = map[string][]byte{
	<<- range $k, $v := .M >>
	"<< $k >>": << $v | printf "%#v" | trimPrefix >>,
	<<- end >>
}`

func main() {
	flag.StringVar(&flags.Package, "package", "main", "package name")
	flag.StringVar(&flags.Var, "var", "files", "variable name of map")
	flag.BoolVar(&flags.Abs, "abs", false, "use absolute path for keys")
	flag.Usage = func() {
		stderr.Println(helpString)
		os.Exit(2)
	}
	flag.Parse()

	roots := flag.Args()
	for _, dirRoot := range roots {
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
			stderr.Println(errors.New(`fmap: error: path argument must be a directory`))
			os.Exit(1)
		}
	}

	tmpl, err := template.New("file").Funcs(template.FuncMap{
		"trimPrefix": func(s string) string { return strings.TrimPrefix(s, "[]byte") },
	}).Delims("<<", ">>").Parse(fileTmpl)
	if err != nil {
		stderr.Println(err)
		os.Exit(1)
	}

	mLock := sync.Mutex{}
	m := make(map[string][]byte)

	for _, dirRoot := range roots {
		out := fileContents(dirRoot)
		for f := range out {
			if f.Err != nil {
				stderr.Println(f.Err)
				os.Exit(1)
			}

			k := f.Path
			if flags.Abs {
				absP, err := filepath.Abs(f.Path)
				if err != nil {
					stderr.Println(err)
					os.Exit(1)
				}
				k = absP
			}
			mLock.Lock()
			m[filepath.ToSlash(k)] = f.Content
			mLock.Unlock()
		}
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

// File represents the path and contents of a file.
// Err has no real relation, but makes it convenient to
// pass any errors encountered in fileContents.
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
