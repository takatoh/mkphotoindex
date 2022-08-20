package core

type Photo struct {
	File    string
	Thumb   string
	Caption string
}

func NewPhoto(file, thumb, caption string) *Photo {
	p := new(Photo)
	p.File = file
	p.Thumb = thumb
	p.Caption = caption
	return p
}

type PhotoSet struct {
	Photos []*Photo
	Size   uint
}

func NewPhotoSet(photos []*Photo, size uint) *PhotoSet {
	p := new(PhotoSet)
	p.Photos = photos
	p.Size = size
	return p
}
