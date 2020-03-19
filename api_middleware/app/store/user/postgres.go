package user

import (
	"encoding/json"
	"github.com/jackc/pgx"
	"github.com/pkg/errors"
	"log"
	"time"
)

// PgStore is a Store interface implementation over PostgresSQL
type PgStore struct {
	ConnStr string

	connPool *pgx.ConnPool
}

// NewPgStore creates a new connection pool to postgres storage with specified parameters
//
// connStr should be provided in format: dbdriver://uname:password@address:port/dbname?[param1=][&param2=][...etc]
func NewPgStore(connStr string) (*PgStore, error) {
	connConf, err := pgx.ParseConnectionString(connStr)

	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse pg user Store with connstr %s", connStr)
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
		return nil, errors.Wrapf(err, "failed to initialize pg user Store with connstr %s", connStr)
	}

	return &PgStore{
		ConnStr:  connStr,
		connPool: p,
	}, nil
}

// put saves user in postgres Store
func (p *PgStore) put(user User) (int, error) {
	tx, err := p.connPool.Begin()
	if err != nil {
		return 0, errors.Wrap(err, "failed to start insert transaction into pg users store")
	}

	defer func() {
		errNested := tx.Rollback()
		if errNested != nil && errNested != pgx.ErrTxClosed {
			log.Printf("[ERROR] failed to rollback the transaction (put in pgUsr): %+v", err)
		}
	}()

	privsMarshalled, err := json.Marshal(user.Privileges)

	if err != nil {
		return 0, errors.Wrapf(err, "failed to insert user: failed to marshal privileges")
	}

	row := tx.QueryRow("INSERT INTO "+
		"users(email, password, privileges, created_at, updated_at) "+
		"VALUES ($1, $2, $3, $4, $5) "+
		"RETURNING id",
		user.Email,
		user.Password,
		privsMarshalled,
		time.Now(),
		time.Now(),
	)

	var id int
	err = row.Scan(&id)
	if err != nil {
		return 0, errors.Wrap(err, "failed to scan user ID while inserting")
	}

	err = tx.Commit()

	if err != nil {
		return 0, errors.Wrap(err, "failed to commit transaction user ID while inserting")
	}

	return id, nil
}

// Get returns user by its ID, if user is present
func (p *PgStore) Get(id int) (User, error) {
	panic("implement me")
}

// List returns the list of all user stored in the database
func (p *PgStore) List() ([]User, error) {
	panic("implement me")
}

// GetByEmail returns user's data by its email
func (p *PgStore) GetByEmail(email string) (User, error) {
	panic("implement me")
}

// GetAuthData provides basic auth information about user, such as email, hashed password and
// its privileges in format map[privilege]given
func (p *PgStore) GetAuthData(id int) (string, string, map[string]bool, error) {
	panic("implement me")
}

// Delete deletes user, given by its ID, from the postgres database
func (p *PgStore) Delete(id int) error {
	panic("implement me")
}

// Update updates user information stored in the database (everything except ID, email and password)
func (p *PgStore) Update(user User) error {
	panic("implement me")
}
