// This program must be run from `../` as `go run generate/main.go`
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/go-bindata/go-bindata/v3"
)

const (
	outputName = "bindata.go"
	ignoreFile = `.*\.go`

	staticsName   = "statics"
	templatesName = "templates"
)

func main() {
	generate()
}

// Generate bindata.go
func generate() {
	mustGenerate(staticsName, true)
	mustGenerate(templatesName, false)
}

func mustGenerate(path string, fs bool) {
	cleanPath := filepath.Clean(path)
	config := bindata.Config{
		Package: cleanPath,
		Output:  filepath.Join(cleanPath, outputName),
		Prefix:  cleanPath + "/",
		Input: []bindata.InputConfig{
			bindata.InputConfig{
				Path:      cleanPath,
				Recursive: false,
			},
		},
		Ignore: []*regexp.Regexp{
			regexp.MustCompile(ignoreFile),
		},
		HttpFileSystem: fs,
	}

	err := bindata.Translate(&config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "bindata %v, %v: %v\n", path, fs, err)
		os.Exit(1)
	}
}
