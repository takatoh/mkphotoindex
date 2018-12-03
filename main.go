package main

import (
	"fmt"
	"os"
	"path/filepath"
	"html/template"
)

func main() {
	var photos []*Photo

	pattern := os.Args[1]
	filenames, _ := filepath.Glob(pattern)

	t := template.Must(template.ParseFiles("index.tmpl"))
	w, err := os.OpenFile("index.html", os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Println(err)
	}

	for _, filename := range filenames {
		photos = append(photos, newPhoto(filename, filename, filename))
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
