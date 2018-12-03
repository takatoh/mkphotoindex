package main

import (
	"fmt"
	"os"
	"path/filepath"
	"html/template"
)

func main() {
	pattern := os.Args[1]
	filenames, _ := filepath.Glob(pattern)

	t := template.Must(template.ParseFiles("index.tmpl"))
	w, err := os.OpenFile("index.html", os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Println(err)
	}


	err = t.ExecuteTemplate(w, "index.tmpl", filenames)
	if err != nil {
		panic(err)
	}
}
