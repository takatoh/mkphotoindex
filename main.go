package main

import (
	"flag"
	"fmt"
	"text/template"
	"image"
	"image/jpeg"
	"os"
	"path/filepath"
	"strings"

	"github.com/nfnt/resize"
)

const (
	progVersion = "v0.3.0"
	tmpl = `<!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8" />
    <title>Index of photos</title>
  </head>
  <body>
    <h1>Index of photos</h1>
    <ul style="list-style-type: none;">
      {{range .}}
      <li>
        <figure style="float: left;">
          <a href="{{.File}}" target="_blank">
            <img src="{{.Thumb}}" />
          </a>
          <figcaption>{{.Caption}}</figcaption>
        </figure>
      </li>
      {{end}}
    </ul>
  </body>
</html>
`
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
	opt_version := flag.Bool("version", false, "Show version")
	flag.Parse()

	if *opt_version {
		fmt.Println(progVersion)
		os.Exit(0)
	}

	var photos []*Photo
	var imgFiles []string
	var dir string
	var pattern string
	var thumbsDir string
	var indexFile string

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
		if strings.ToLower(ext) == ".jpg" {
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
	t, _ := template.New("index").Parse(tmpl)
	w, err := os.OpenFile(indexFile, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Println(err)
	}

	for _, imgFile := range imgFiles {
		thumb := makeThumbnail(imgFile, dir)
		filename := filepath.Base(imgFile)
		photos = append(photos, newPhoto(filename, thumb, filename))
	}

	err = t.ExecuteTemplate(w, "index", photos)
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

func makeThumbnail(srcfile, dir string) string {
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
	filename := filepath.Base(srcfile)
	var thumbFile string
	if dir != "" {
		thumbFile = dir + "/thumbs/thumb_" + filename
	} else {
		thumbFile = "thumbs/thumb_" + filename
	}
	thumb, _ := os.Create(thumbFile)
	jpeg.Encode(thumb, resizedImg, nil)
	thumb.Close()

	return "thumbs/thumb_" + filename
}
