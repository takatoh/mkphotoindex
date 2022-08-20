package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/takatoh/mkphotoindex/core"
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

	var photos []*core.Photo
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

	for _, imgFile := range imgFiles {
		thumb := thumbnail.MakeThumbnail(imgFile, dir, *opt_size)
		filename := filepath.Base(imgFile)
		ext := filepath.Ext(filename)
		caption := strings.Replace(filename, ext, "", 1)
		photos = append(photos, core.NewPhoto(filename, thumb, caption))
	}

	photoSet := core.NewPhotoSet(photos, *opt_size)

	err = html.MakeIndex(w, photoSet)
	if err != nil {
		panic(err)
	}
}

//func makeThumbnail(srcfile, dir string, size uint) string {
//	src, _ := os.Open(srcfile)
//	defer src.Close()
//
//	config, _, _ := image.DecodeConfig(src)
//	src.Seek(0, 0)
//	img, _, _ := image.Decode(src)
//
//	var resizedImg image.Image
//	if config.Width >= config.Height {
//		resizedImg = resize.Resize(size, 0, img, resize.Lanczos3)
//	} else {
//		resizedImg = resize.Resize(0, size, img, resize.Lanczos3)
//	}
//	filename := filepath.Base(srcfile)
//	var thumbFile string
//	if dir != "" {
//		thumbFile = dir + "/thumbs/thumb_" + filename
//	} else {
//		thumbFile = "thumbs/thumb_" + filename
//	}
//	thumb, _ := os.Create(thumbFile)
//	ext := strings.ToLower(filepath.Ext(filename))
//	if ext == ".jpg" || ext == ".jpeg" {
//		jpeg.Encode(thumb, resizedImg, nil)
//	} else {
//		png.Encode(thumb, resizedImg)
//	}
//	thumb.Close()
//
//	return "thumbs/thumb_" + filename
//}

func contains(s []string, e string) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}
