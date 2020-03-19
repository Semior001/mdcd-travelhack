package image

import (
	"github.com/Semior001/mdcd-travelhack/app/store/user"
	"time"
)

const (
	ImgTypeBackground = "background"
	ImgTypeSrc        = "source"
	ImgTypeDerived    = "derived"
	ImgTypeCommited   = "commited"
)

// Image describes a particular image added by any user
type Image struct {
	ID            int
	BarCode       string
	ImgType       string
	Mime          string // mime type of an image, e.g. "gif", "jpg", etc.
	LocalFilename string
	UserID        int
	AddedBy       *user.User
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// Store defines an interface to put and load images from the database
type Store interface {
	putImage(imgMetaData Image) (imgID int, err error)
	getImage(id int) (imgMetaData Image, err error)
	GetBackgrounds() (ids []int, err error)
	CheckBarcode(barcode string) (ok bool, err error)
	getImageByBarcode(barcode string) (imgMetaData Image, err error)
}
