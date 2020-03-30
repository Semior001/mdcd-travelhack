package image

import (
	"github.com/Semior001/mdcd-travelhack/app/store/user"
	"github.com/pkg/errors"
	"github.com/segmentio/ksuid"
	"io"
	"log"
	"os"
	"path"
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
	ID            uint64
	Barcode       string
	ImgType       string
	Mime          string // mime type of an image, e.g. "gif", "jpg", etc.
	LocalFilename string
	UserID        uint64
	AddedBy       *user.User
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// Store defines an interface to put and load images from the database
type Store interface {
	putImage(imgMetaData Image) (imgID uint64, err error)
	getImage(id uint64) (imgMetaData Image, err error)
	GetBackgroundIds() (ids []uint64, err error)
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
func (s *Service) Save(userID uint64, barcode string, imgType string, mime string, reader io.Reader) (imgID uint64, err error) {
	fileName := ksuid.New().String()
	filePath := path.Join(s.LocalStoragePath, fileName)

	file, err := os.Create(filePath)
	if err != nil {
		return 0, errors.Wrapf(err, "failed to create image file in %s", filePath)
	}

	defer func() {
		if err = file.Close(); err != nil {
			log.Printf("[WARN] failed to close image file %s", fileName)
		}
	}()

	_, err = io.Copy(file, reader)
	if err != nil {
		return 0, errors.Wrap(err, "failed to copy image to file")
	}

	imgID, err = s.putImage(Image{
		Barcode:       barcode,
		ImgType:       imgType,
		Mime:          mime,
		LocalFilename: fileName,
		UserID:        userID,
	})
	if err != nil {
		return 0, errors.Wrap(err, "failed to put image in db")
	}

	return imgID, nil
}

// Load returns image from local storage and its metadata from database by id
func (s *Service) Load(imgID uint64) (Image, io.ReadCloser, error) {
	imgMetaData, err := s.getImage(imgID)
	if err != nil {
		return Image{}, nil, errors.Wrap(err, "failed to find image metadata in db")
	}

	file, err := os.Open(path.Join(s.LocalStoragePath, imgMetaData.LocalFilename))
	if err != nil {
		return Image{}, nil, errors.Wrap(err, "failed to load image from storage")
	}
	return imgMetaData, file, nil
}

// GetImgByBarcode returns image from local storage and its metadata from database by barcode
func (s *Service) GetImgByBarcode(barcode string) (Image, io.ReadCloser, error) {
	imgMetaData, err := s.getImgByBarcode(barcode)
	if err != nil {
		return Image{}, nil, errors.Wrap(err, "failed to find image metadata in db")
	}

	file, err := os.Open(path.Join(s.LocalStoragePath, imgMetaData.LocalFilename))
	if err != nil {
		return Image{}, nil, errors.Wrap(err, "failed to load image from storage")
	}

	return imgMetaData, file, nil
}
