package image

import (
	"github.com/Semior001/mdcd-travelhack/app/store/user"
	"github.com/go-pg/pg/v9"
	R "github.com/go-pkgz/rest"
	"github.com/pkg/errors"
	"github.com/segmentio/ksuid"
	"io"
	"log"
	"os"
	"path"
	"strings"
	"time"
)

type Image struct {
	ID            uint64
	BarCode       string
	ImgType       string
	LocalFilename string
	UserId        uint64
	AddedBy       *user.User
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type Store interface {
	Migrate(force bool) error
	putImage(imgMetaData Image) (imgId uint64, err error)
	getImage(id uint64) (imgMetaData Image, err error)
	GetBackgrounds() (ids []uint64, err error)
	CheckBarcode(barcode string) (json R.JSON, err error)
}

type Service struct {
	Store

	LocalStoragePath string
}

// ServiceOpts defines options to create connection with storage
type ServiceOpts struct {
	Driver           string
	User             string
	Password         string
	Source           string
	LoggerFlags      int
	LocalStoragePath string
}

func NewService(opts ServiceOpts) (*Service, error) {
	var db Store
	var err error

	if err := os.MkdirAll(opts.LocalStoragePath, 0750); err != nil {
		return nil, errors.Wrap(err, "failed to mkdir local media path")
	}

	switch opts.Driver {
	case "postgres":
		db, err = NewPgImageStorage(pg.Options{
			User:     opts.User,
			Password: opts.Password,
			Database: strings.Split(opts.Source, "@")[0],
			Addr:     strings.Split(opts.Source, "@")[1],
		}, log.New(os.Stdout, "pgstorage", opts.LoggerFlags))
	}

	if err != nil {
		return nil, err
	}
	return &Service{
		Store: db,
	}, nil
}

// PutImage stores image in the local path, given to the Service
// and stores all metadata in the database via calling Store.putImage
func (s *Service) PutImage(userId uint64, barcode string, imgType string, reader io.Reader) (imgId uint64, err error) {
	fileName := ksuid.New().String()
	imFile := path.Join(s.LocalStoragePath, fileName)

	fh, err := os.Create(imFile)
	if err != nil {
		return 0, errors.Wrapf(err, "can't create image file")
	}

	defer func() {
		if err := fh.Close(); err != nil {
			log.Printf("[DEBUG] can't close image file")
		}
	}()

	_, err = io.Copy(fh, reader)
	if err != nil {
		return 0, errors.Wrapf(err, "can't save image file")
	}
	imgId, err = s.putImage(Image{
		LocalFilename: fileName,
		UserId:        userId,
		ImgType:       imgType,
		BarCode:       barcode,
	})
	if err != nil {
		return 0, errors.Wrapf(err, "can't put image metadata into db")
	}
	return imgId, err
}

func (s *Service) GetImage(imgId uint64) (imgMetaData Image, reader io.ReadCloser, err error) {
	imgMetaData, err = s.getImage(imgId)
	if err != nil {
		return imgMetaData, nil, errors.Wrap(err, "can't find image in db")
	}

	fh, err := os.Open(path.Join(s.LocalStoragePath, imgMetaData.LocalFilename))
	if err != nil {
		return imgMetaData, nil, errors.Wrap(err, "can't load image from local media path")
	}
	return imgMetaData, fh, nil
}
