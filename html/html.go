package html

import (
	"os"
	"text/template"

	"github.com/takatoh/mkphotoindex/core"
)

const (
	IndexTmpl = `<!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8" />
    <title>Index of photos</title>
    <style>
      ul {list-style-type: none;}
      figure {float: left;}
      div.photo {width: {{.Size}}px; height: {{.Size}}px;}
    </style>
  </head>
  <body>
    <h1>Index of photos</h1>
    <ul>
      {{range .Photos}}
      <li>
        <figure>
          <div class="photo">
            <a href="{{.File}}" target="_blank">
              <img src="{{.Thumb}}" />
            </a>
          </div>
          <figcaption>{{.Caption}}</figcaption>
        </figure>
      </li>
      {{end}}
    </ul>
  </body>
</html>
`
)

func MakeIndex(file *os.File, photoSet *core.PhotoSet) error {
	t, _ := template.New("index").Parse(IndexTmpl)
	err := t.ExecuteTemplate(file, "index", photoSet)
	return err
}
