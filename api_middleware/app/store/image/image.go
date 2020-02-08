package image

import (
	"github.com/Semior001/mdcd-travelhack/app/store/user"
	"github.com/go-pg/pg/v9"
	"log"
	"os"
	"strings"
	"time"
)

type Image struct {
	ID          uint64
	LocalPath   string
	AddedUserId uint64
	AddedBy     *user.User
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Store interface {
	Migrate(force bool) error
	putImage(image *Image) (id uint64, err error)
	GetImage(id uint64) (image *Image, err error)
}

type Service struct {
	Store
}

// ServiceOpts defines options to create connection with storage
type ServiceOpts struct {
	Driver      string
	User        string
	Password    string
	Source      string
	LoggerFlags int
}

func NewService(opts ServiceOpts) (*Service, error) {
	var db Store
	var err error

	switch opts.Driver {
	case "postgres":
		db, err = NewPgStorage(pg.Options{
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
