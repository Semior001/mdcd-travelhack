package image

import (
	"github.com/Semior001/mdcd-travelhack/app/store/user"
	"github.com/pkg/errors"
	"io"
	"os"
	"strings"
	"time"
)

//go:generate mockery -inpkg -name Store -case snake

const (
	ImgTypeBackground = "background"
	ImgTypeSrc        = "source"
	ImgTypeDerived    = "derived"
	ImgTypeCommitted  = "committed"
)

// Image describes a particular image added by any user
type Image struct {
	ID            int
	Barcode       string
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
	GetBackgroundIds() (ids []int, err error)
	CheckBarcode(barcode string) (ok bool, err error)
	getImgByBarcode(barcode string) (imgMetaData Image, err error)
}

// Service provides methods for operating, processing and storing images
type Service struct {
	LocalStoragePath string
	Store
}

// ServiceOpts defines options to create connection with storage
//
// ConnStr should be provided in format: dbdriver://uname:password@address:port/dbname?[param1=][&param2=][...etc]
type ServiceOpts struct {
	ConnStr          string
	LocalStoragePath string
}

func NewService(opts ServiceOpts) (*Service, error) {
	if err := os.MkdirAll(opts.LocalStoragePath, 0750); err != nil {
		return nil, errors.Wrap(err, "failed to mkdir local media path")
	}

	driverEndIdx := strings.Index(opts.ConnStr, "://")
	driver := opts.ConnStr[0:driverEndIdx]

	var db Store
	var err error

	switch driver {
	case "postgres":
		db, err = NewPgStore(opts.ConnStr)
	}

	if err != nil {
		return nil, err
	}
	return &Service{
		Store:            db,
		LocalStoragePath: opts.LocalStoragePath,
	}, nil
}

// Save stores image in the local path, given to the Service
// and stores all metadata in the database via calling Store.putImage
func (s *Service) Save(userId uint64, barcode string, imgType string, reader io.Reader) (imgId uint64, err error) {
	panic("implement me")
}

func (s *Service) Load(imgId uint64) (Image, io.ReadCloser, error) {
	panic("implement me")
}

func (s *Service) GetImgByBarcode(barcode string) (Image, io.ReadCloser, error) {
	panic("implement me")
}
