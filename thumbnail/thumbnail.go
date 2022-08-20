package thumbnail

import (
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/nfnt/resize"
)

func MakeThumbnail(srcfile, dir string, size uint) string {
	src, _ := os.Open(srcfile)
	defer src.Close()

	config, _, _ := image.DecodeConfig(src)
	src.Seek(0, 0)
	img, _, _ := image.Decode(src)

	var resizedImg image.Image
	if config.Width >= config.Height {
		resizedImg = resize.Resize(size, 0, img, resize.Lanczos3)
	} else {
		resizedImg = resize.Resize(0, size, img, resize.Lanczos3)
	}
	filename := filepath.Base(srcfile)
	var thumbFile string
	if dir != "" {
		thumbFile = dir + "/thumbs/thumb_" + filename
	} else {
		thumbFile = "thumbs/thumb_" + filename
	}
	thumb, _ := os.Create(thumbFile)
	ext := strings.ToLower(filepath.Ext(filename))
	if ext == ".jpg" || ext == ".jpeg" {
		jpeg.Encode(thumb, resizedImg, nil)
	} else {
		png.Encode(thumb, resizedImg)
	}
	thumb.Close()

	return "thumbs/thumb_" + filename
}
