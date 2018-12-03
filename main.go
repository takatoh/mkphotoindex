package main

import (
	"fmt"
	"os"
	"path/filepath"
	"html/template"
	"image"
	"image/jpeg"

	"github.com/nfnt/resize"
)

func main() {
	var photos []*Photo

	pattern := "*.jpg"
	filenames, _ := filepath.Glob(pattern)

	t := template.Must(template.ParseFiles("index.tmpl"))
	w, err := os.OpenFile("index.html", os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Println(err)
	}

	for _, filename := range filenames {
		thumb := makeThumbnail(filename)
		photos = append(photos, newPhoto(filename, thumb, filename))
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
