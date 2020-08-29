package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/go-bindata/go-bindata/v3"
)

const (
	bindataOutput = "bindata.go"
	bindataIgnore = `.*\.go`

	staticsName   = "statics"
	templatesName = "templates"
)

// Generate bindata.go
func generateBindata() {
	mustGenerateBindata(staticsName, true)
	mustGenerateBindata(templatesName, false)
}

func mustGenerateBindata(path string, fs bool) {
	cleanPath := filepath.Clean(path)
	config := bindata.Config{
		Package: cleanPath,
		Output:  filepath.Join(cleanPath, bindataOutput),
		Prefix:  cleanPath + "/",
		Input: []bindata.InputConfig{
			bindata.InputConfig{
				Path:      cleanPath,
				Recursive: false,
			},
		},
		Ignore: []*regexp.Regexp{
			regexp.MustCompile(bindataIgnore),
		},
		HttpFileSystem: fs,
	}

	err := bindata.Translate(&config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "bindata %v, %v: %v\n", path, fs, err)
		os.Exit(1)
	}
}
