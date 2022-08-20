package contactsheet

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/signintech/gopdf"

	"github.com/takatoh/mkphotoindex/core"
)

const (
	defaultCols = 3
	defaultRows = 4
)

func Generate(imageFiles *core.PhotoSet, thumbsDir string) {
	pdf := gopdf.GoPdf{}

	A4 := *gopdf.PageSizeA4

	pdf.Start(gopdf.Config{PageSize: A4})
	pdf.AddPage()

	// Draw grid
	drawGrid(&pdf, &A4)

	// Font
	err := pdf.AddTTFFont("IPAex", "ipaexg.ttf")
	if err != nil {
		panic(err)
	}

	// Title
	pdf.SetFont("IPAex", "", 24)
	drawText(&pdf, 100, 25, "Index of photos")

	// Drow images
	pages, totalPage := paginate(imageFiles.Photos, defaultRows*defaultCols)
	for j, page := range pages {
		for i, img := range page {
			x := 100.0 + 150.0*float64(i%defaultCols)
			y := 80.0 + 150.0*float64((i/defaultCols))
			thumb := strings.Replace(img.Thumb, "thumbs", thumbsDir, 1)
			drawImage(&pdf, x, y, thumb)
			pdf.SetFont("IPAex", "", 10)
			basename := filepath.Base(img.File)
			drawText(&pdf, x, y+120.0, basename)
		}
		pdf.SetFont("IPAex", "", 12)
		drawText(&pdf, 265, 800, "page "+strconv.Itoa(j+1)+" of "+strconv.Itoa(totalPage))
		if j < totalPage-1 {
			pdf.AddPage()
		}
	}

	// Write PDF
	pdf.WritePdf("contactsheet.pdf")

}

func drawText(pdf *gopdf.GoPdf, x float64, y float64, s string) {
	pdf.SetX(x)
	pdf.SetY(y)
	pdf.Cell(nil, s)
}

func drawImage(pdf *gopdf.GoPdf, x float64, y float64, filename string) {
	pdf.Image(filename, x, y, nil)
}

func drawGrid(pdf *gopdf.GoPdf, page *gopdf.Rect) {
	ww := 10.0
	for i := 1; i < int(page.H/ww); i++ {
		if i%10 == 0 {
			pdf.SetLineWidth(0.8)
			pdf.SetStrokeColor(50, 50, 100)
		} else {
			pdf.SetLineWidth(0.2)
			pdf.SetStrokeColor(100, 100, 130)
		}
		x, y := float64(i)*ww, float64(i)*ww
		pdf.Line(x, 0, x, page.H)
		pdf.Line(0, y, page.W, y)
	}
}

func MakeDirectory(baseDir string) string {
	var thumbsDir string

	if baseDir != "" {
		thumbsDir = baseDir + "/_csheet_thumbs"
	} else {
		thumbsDir = "_csheet_thumbs"
	}
	if _, err := os.Stat(thumbsDir); os.IsNotExist(err) {
		os.Mkdir(thumbsDir, 0777)
	}

	return thumbsDir
}

func paginate(photos []*core.Photo, parPage int) ([][]*core.Photo, int) {
	var pages [][]*core.Photo

	numOfPhotos := len(photos)
	totalPage := numOfPhotos / parPage
	if numOfPhotos%parPage != 0 {
		totalPage = totalPage + 1
	}
	for i := 0; i < numOfPhotos; i = i + parPage {
		if i+parPage < numOfPhotos {
			pages = append(pages, photos[i:i+parPage])
		} else {
			pages = append(pages, photos[i:numOfPhotos])
		}
	}

	return pages, totalPage
}
