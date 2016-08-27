# fmap

Generate a Go map of static file contents.

## Install

```
go get -u github.com/nishanths/fmap
```

## What

```
fmap generates a go source file containing a map[string][]byte
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
  fmap static/css static/js | gofmt > static_files.go
```

## Test

```
go test -race ./...
```

## License

[MIT](https://nishanths.mit-license.org)
