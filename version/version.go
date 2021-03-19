package version

import (
	"fmt"
	"runtime"
)

var (
	Version = "dev"
	Hash    = "unknown"
)

// Display version
func Display() {
	fmt.Println("github.com/forewing/goldennum")
	fmt.Printf("version\t%v\n", Version)
	fmt.Printf("commit\t%v\n", Hash)
	fmt.Printf("runtime\t%v\n", runtime.Version())
}
