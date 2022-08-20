package contactsheet

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/signintech/gopdf"

	"github.com/takatoh/mkphotoindex/core"
)

func Generate(imageFiles *core.PhotoSet, thumbsDir string) {
	pdf := gopdf.GoPdf{}

	A4 := *gopdf.PageSizeA4

	pdf.Start(gopdf.Config{PageSize: A4})
	pdf.AddPage()

	// Draw grid
	drawGrid(&pdf, &A4)

	// Font
	err := pdf.AddTTFFont("Arial", "C:/windows/fonts/arial.ttf")
	if err != nil {
		panic(err)
	}

	//var imageFiles = []string{
	//	"./photos/thumbs/thumb_01_1024768.jpg",
	//	"./photos/thumbs/thumb_01_12801024.jpg",
	//	"./photos/thumbs/thumb_01_19201080.jpg",
	//	"./photos/thumbs/thumb_01_19201200.jpg",
	//}

	// Drow images
	for i, img := range imageFiles.Photos {
		x := 100.0 + 150.0*float64(i%3)
		y := 50.0 + 150.0*float64((i/3))
		thumb := strings.Replace(img.Thumb, "thumbs", thumbsDir, 1)
		drawImage(&pdf, x, y, thumb)
		pdf.SetFont("Arial", "", 10)
		basename := filepath.Base(img.File)
		drawText(&pdf, x, y+120.0, basename)
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
			pdf.SetLineWidth(0.3)
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
