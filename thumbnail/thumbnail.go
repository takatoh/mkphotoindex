package thumbnail

import (
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/nfnt/resize"

	"github.com/takatoh/mkphotoindex/core"
)

func MakeThumbnails(imgFiles []string, thumbsDir string, opt_size uint) *core.PhotoSet {
	var photos []*core.Photo

	for _, imgFile := range imgFiles {
		thumb := MakeThumbnail(imgFile, thumbsDir, opt_size)
		filename := filepath.Base(imgFile)
		ext := filepath.Ext(filename)
		caption := strings.Replace(filename, ext, "", 1)
		photos = append(photos, core.NewPhoto(filename, thumb, caption))
	}

	photoSet := core.NewPhotoSet(photos, opt_size)
	return photoSet
}

func MakeThumbnail(srcfile, thumbsDir string, size uint) string {
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
	thumbFile = thumbsDir + "/thumb_" + filename
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
