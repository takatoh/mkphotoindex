package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/takatoh/mkphotoindex/html"
	"github.com/takatoh/mkphotoindex/thumbnail"
)

const (
	progVersion = "v0.5.2"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr,
			`Usage:
  %s [options] [photo_dir]

Options:
`, os.Args[0])
		flag.PrintDefaults()
	}
	opt_size := flag.Uint("size", 320, "Specify thubnail size")
	opt_version := flag.Bool("version", false, "Show version")
	flag.Parse()

	if *opt_version {
		fmt.Println(progVersion)
		os.Exit(0)
	}

	var imgFiles []string
	var dir string
	var pattern string
	var thumbsDir string
	var indexFile string
	var photoTypes []string = []string{".jpg", ".jpeg", ".png"}

	if flag.NArg() > 0 {
		dir = flag.Arg(0)
	} else {
		dir = ""
	}

	if dir != "" {
		pattern = dir + "/*.*"
	} else {
		pattern = "*.*"
	}

	filenames, _ := filepath.Glob(pattern)
	for _, f := range filenames {
		ext := filepath.Ext(f)
		if contains(photoTypes, strings.ToLower(ext)) {
			imgFiles = append(imgFiles, f)
		}
	}
	if dir != "" {
		thumbsDir = dir + "/thumbs"
	} else {
		thumbsDir = "thumbs"
	}
	if _, err := os.Stat(thumbsDir); os.IsNotExist(err) {
		os.Mkdir(thumbsDir, 0777)
	}

	if dir != "" {
		indexFile = dir + "/index.html"
	} else {
		indexFile = "index.html"
	}
	w, err := os.OpenFile(indexFile, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Println(err)
	}

	photoSet := thumbnail.MakeThumbnails(imgFiles, thumbsDir, *opt_size)

	err = html.MakeIndex(w, photoSet)
	if err != nil {
		panic(err)
	}
}

func contains(s []string, e string) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}
