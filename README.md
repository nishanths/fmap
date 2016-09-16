# fmap

Generate a `map[string][]byte` of directory contents.

[![wercker status](https://app.wercker.com/status/d946d386cadef972e6dc50cef520b6a1/s/master "wercker status")](https://app.wercker.com/project/byKey/d946d386cadef972e6dc50cef520b6a1)

## Install

```
go get -u github.com/nishanths/fmap
```

## What

```
fmap generates a go source file containing a map[string][]byte
for the specified paths. The keys in the generated map are file paths
and the values are the contents of the file at that path.

The generated go file is printed to stdout (typically you would need to run
the output through gofmt). Empty directories are ignored, and symlinks are not followed.

usage: fmap [flags] [file ...]

flags:
  -package  package name to use in generated file (default: "main")
  -var      variable name of the map (default: "files")
  -abs      use absolute paths for map keys

example:
  fmap static/css static/js favicon.ico | gofmt > files.go
```

## Test

```
go test -race ./...
```

## License

[MIT](https://nishanths.mit-license.org)
