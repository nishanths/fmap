# fmap

```
fmap generates a go source file containing a map[string][]byte that
corresponds to the specified directory tree. The keys are the paths
of files and the values are the contents of the file at the path.

The generated file is printed to stdout. Empty directories are
ignored, and symlinks are not followed.

usage:
  fmap [flags] /path/to/dir

flags:
  -package  package name to use in generated file (default: "main")
  -var      variable name of the map (default: "files")
```

## Install

```
go get -u github.com/nishanths/fmap
```

## Test

```
go test -race ./...
```

## License

[MIT](https://nishanths.mit-license.org)
