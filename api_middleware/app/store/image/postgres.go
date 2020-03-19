package image

import (
	"github.com/jackc/pgx"
	"github.com/pkg/errors"
	"time"
)

// PgStore is a store interface implementation over PostgresSQL
type PgStore struct {
	ConnStr string

	connPool *pgx.ConnPool
}

// NewPgStore creates a connection pool to the postgres storage
func NewPgStore(connStr string) (*PgStore, error) {
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

	return &PgStore{
		ConnStr:  connStr,
		connPool: p,
	}, nil
}

func (p *PgStore) putImage(imgMetaData Image) (int, error) {
	panic("implement me")
}

func (p *PgStore) getImage(id int) (Image, error) {
	panic("implement me")
}

func (p *PgStore) GetBackgroundIds() ([]int, error) {
	panic("implement me")
}

func (p *PgStore) CheckBarcode(barcode string) (bool, error) {
	panic("implement me")
}

func (p *PgStore) getImgByBarcode(barcode string) (Image, error) {
	panic("implement me")
}
