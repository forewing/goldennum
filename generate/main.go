// This program must be run from `../` as `go run generate/main.go`
package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"runtime"

	"github.com/go-bindata/go-bindata/v3"
)

const (
	output = "./bindata.go"
)

var (
	config = bindata.Config{
		Package: "main",
		Output:  output,

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
	if runtime.GOOS == "windows" {
		fixWindowsPathDelimiter()
	}
}

// Generate bindata.go
func generate() {
	err := bindata.Translate(&config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "bindata: %v\n", err)
		os.Exit(1)
	}
}

// Fix go-bindata's BUG on windows
func fixWindowsPathDelimiter() {
	data, err := ioutil.ReadFile(output)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fix: %v\n", err)
		os.Exit(1)
	}
	text := string(data)

	text = regexp.MustCompile(`templates\\{1,2}`).ReplaceAllString(text, "templates/")

	err = ioutil.WriteFile(output, []byte(text), 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fix: %v\n", err)
		os.Exit(1)
	}
}
