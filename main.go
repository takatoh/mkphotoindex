package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/takatoh/mkphotoindex/contactsheet"
	"github.com/takatoh/mkphotoindex/core"
	"github.com/takatoh/mkphotoindex/html"
	"github.com/takatoh/mkphotoindex/thumbnail"
)

const (
	progVersion = "v0.7.0"
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
	opt_csheet := flag.Bool("contact-sheet", false, "Genarete contact sheet, instead")
	opt_title := flag.String("title", "", "Specify title")
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
	var photoSet *core.PhotoSet

	if flag.NArg() > 0 {
		dir = flag.Arg(0)
		pattern = dir + "/*.*"
	} else {
		dir = ""
		pattern = "*.*"
	}

	filenames, _ := filepath.Glob(pattern)
	for _, f := range filenames {
		ext := filepath.Ext(f)
		if contains(photoTypes, strings.ToLower(ext)) {
			imgFiles = append(imgFiles, f)
		}
	}

	if *opt_csheet {
		// Generate contact sheet
		var title string
		if *opt_title != "" {
			title = *opt_title
		} else {
			title = dir
		}
		thumbsDir = thumbnail.MakeDirectory(dir, "_csheet_thumbs")
		photoSet = thumbnail.MakeThumbnails(imgFiles, thumbsDir, 200)
		contactsheet.Generate(photoSet, thumbsDir, title)
		err := os.RemoveAll(thumbsDir)
		if err != nil {
			panic(err)
		}
	} else {
		// Or generate index.html
		thumbsDir = thumbnail.MakeDirectory(dir, "thumbs")
		photoSet = thumbnail.MakeThumbnails(imgFiles, thumbsDir, *opt_size)
		indexFile = html.IndexFilePath(dir)
		w, err := os.OpenFile(indexFile, os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			fmt.Println(err)
		}
		err = html.MakeIndex(w, photoSet)
		if err != nil {
			panic(err)
		}
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
