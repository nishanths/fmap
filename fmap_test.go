package main

import "testing"

func shouldResemble(t *testing.T, m, n map[string]File) {
	for k := range m {
		if n[k].Err != m[k].Err || n[k].Path != m[k].Path ||
			string(n[k].Content) != string(m[k].Content) {
			t.Fatalf("shouldResemble: m: %v, n: %v", m[k], n[k])
		}
	}
}

func TestFileContents(t *testing.T) {
	ch := fileContents("testdata")
	m := make(map[string]File)
	for f := range ch {
		m[f.Path] = f
	}
	expected := map[string]File{
		"testdata/bar/bar.txt": File{
			Path:    "testdata/bar/bar.txt",
			Content: []byte("boo baa\ncool\n"),
		},
		"testdata/hello/hello.txt": File{
			Path:    "testdata/hello/hello.txt",
			Content: []byte("hello, world\n\n"),
		},
		"testdata/hello/.gitignore": File{
			Path:    "testdata/hello/.gitignore",
			Content: []byte("# this is a comment.\n"),
		},
		"testdata/foo/se7en": File{
			Path:    "testdata/foo/se7en",
			Content: []byte("foo 42\nvim\n"),
		},
		"testdata/index.html": File{
			Path: "testdata/index.html",
			Content: []byte(`<!doctype html>
<title>The C Programming Language</title>
<h1>
    K&R
</h1>
`),
		},
	}
	shouldResemble(t, m, expected)
}
