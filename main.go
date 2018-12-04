package main

import (
	"flag"
	"fmt"
	"html/template"
	"image"
	"image/jpeg"
	"os"
	"path/filepath"
	"strings"

	"github.com/nfnt/resize"
)

const (
	progVersion = "v0.1.0"
)

func main() {
	opt_version := flag.Bool("version", false, "Show version")
	flag.Parse()

	if *opt_version {
		fmt.Println(progVersion)
		os.Exit(0)
	}

	var photos []*Photo
	var imgFiles []string

	pattern := "*.*"
	filenames, _ := filepath.Glob(pattern)
	for _, f := range filenames {
		ext := filepath.Ext(f)
		if strings.ToLower(ext) == ".jpg" {
			imgFiles = append(imgFiles, f)
		}
	}
	thumbsDir := "thumbs"
	if _, err := os.Stat(thumbsDir); os.IsNotExist(err) {
		os.Mkdir(thumbsDir, 0777)
	}

	t := template.Must(template.ParseFiles("index.tmpl"))
	w, err := os.OpenFile("index.html", os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Println(err)
	}

	for _, imgFile := range imgFiles {
		thumb := makeThumbnail(imgFile)
		photos = append(photos, newPhoto(imgFile, thumb, imgFile))
	}

	err = t.ExecuteTemplate(w, "index.tmpl", photos)
	if err != nil {
		panic(err)
	}
}

type Photo struct {
	File    string
	Thumb   string
	Caption string
}

func newPhoto(file, thumb, caption string) *Photo {
	p := new(Photo)
	p.File = file
	p.Thumb = thumb
	p.Caption = caption
	return p
}

func makeThumbnail(srcfile string) string {
	src, _ := os.Open(srcfile)
	defer src.Close()

	config, _, _ := image.DecodeConfig(src)
	src.Seek(0, 0)
	img, _, _ := image.Decode(src)

	var resizedImg image.Image
	if config.Width >= config.Height {
		resizedImg = resize.Resize(320, 0, img, resize.Lanczos3)
	} else {
		resizedImg = resize.Resize(0, 320, img, resize.Lanczos3)
	}
	thumbFile := "thumbs/thumb_" + srcfile
	thumb, _ := os.Create(thumbFile)
	jpeg.Encode(thumb, resizedImg, nil)
	thumb.Close()

	return thumbFile
}
