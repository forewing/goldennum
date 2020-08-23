package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-bindata/go-bindata/v3"
)

var (
	config = bindata.Config{
		Package: "main",
		Output:  "./bindata.go",

		Prefix: "statics/",
		Input: []bindata.InputConfig{
			bindata.InputConfig{
				Path:      filepath.Clean("statics/"),
				Recursive: false,
			},
			bindata.InputConfig{
				Path:      filepath.Clean("templates/"),
				Recursive: false,
			},
		},

		HttpFileSystem: true,
	}
)

func main() {
	generate()
}

// generate bindata.go
func generate() {
	err := bindata.Translate(&config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "bindata: %v\n", err)
		os.Exit(1)
	}
}
