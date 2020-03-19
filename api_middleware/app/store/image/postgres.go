package image

import (
	"github.com/jackc/pgx"
	"github.com/pkg/errors"
	"time"
)

// pgStore is a store interface implementation over PostgresSQL
type pgStore struct {
	ConnStr string

	connPool *pgx.ConnPool
}

// newPgStore creates a connection pool to the postgres storage
func newPgStore(connStr string) (*pgStore, error) {
	connConf, err := pgx.ParseConnectionString(connStr)

	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse pg image store with connstr %s", connStr)
	}

	p, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig:     connConf,
		MaxConnections: 5,
		AfterConnect: func(conn *pgx.Conn) error {
			// todo no-op yet
			return nil
		},
		AcquireTimeout: time.Minute,
	})

	if err != nil {
		return nil, errors.Wrapf(err, "failed to initialize pg image store with connstr %s", connStr)
	}

	return &pgStore{
		ConnStr:  connStr,
		connPool: p,
	}, nil
}

func (p *pgStore) putImage(imgMetaData Image) (int, error) {
	panic("implement me")
}

func (p *pgStore) getImage(id int) (Image, error) {
	panic("implement me")
}

func (p *pgStore) GetBackgrounds() ([]int, error) {
	panic("implement me")
}

func (p *pgStore) CheckBarcode(barcode string) (bool, error) {
	panic("implement me")
}

func (p *pgStore) getImageByBarcode(barcode string) (Image, error) {
	panic("implement me")
}
