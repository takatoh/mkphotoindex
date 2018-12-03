package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	pattern := os.Args[1]
	filenames, _ := filepath.Glob(pattern)

	for _, filename := range filenames {
		fmt.Println(filename)
	}
}
