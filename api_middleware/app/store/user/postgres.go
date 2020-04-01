package user

import (
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
func (p *PgStore) put(user User) (uint64, error) {
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

	row := tx.QueryRow("INSERT INTO "+
		"users(email, password, privileges, created_at, updated_at) "+
		"VALUES ($1, $2, $3, $4, $5) "+
		"RETURNING id",
		user.Email,
		user.Password,
		user.Privileges,
		time.Now(),
		time.Now(),
	)

	var id uint64
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
func (p *PgStore) Get(id uint64) (User, error) {
	user := User{ID: id}

	row := p.connPool.QueryRow("SELECT email, password, privileges, created_at, updated_at "+
		"FROM users WHERE id = $1", id)

	err := row.Scan(&user.Email, &user.Password, &user.Privileges, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return User{}, errors.Wrapf(err, "failed to scan user (id = %d) in Get method", id)
	}

	return user, nil
}

// List returns the list of all user stored in the database
func (p *PgStore) List() ([]User, error) {
	rows, err := p.connPool.Query("SELECT id, email, password, " +
		"privileges, created_at, updated_at FROM users")
	if err != nil {
		return nil, errors.Wrap(err, "failed to select all users")
	}

	defer rows.Close()

	var users []User

	for rows.Next() {
		user := User{}

		err := rows.Scan(&user.ID, &user.Email, &user.Password, &user.Privileges, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan one of users")
		}

		users = append(users, user)
	}

	if rows.Err() != nil {
		return nil, errors.Wrap(rows.Err(), "failed to process all rows in List method")
	}

	return users, nil
}

// GetByEmail returns user's data by its email
func (p *PgStore) GetByEmail(email string) (User, error) {
	user := User{Email: email}

	row := p.connPool.QueryRow("SELECT id, password, privileges, created_at, updated_at "+
		"FROM users WHERE email = $1", email)

	err := row.Scan(&user.ID, &user.Password, &user.Privileges, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return User{}, errors.Wrapf(err, "failed to scan user (email = %s) in GetByEmail", email)
	}

	return user, nil
}

// GetAuthData provides basic auth information about user, such as email, hashed password and
// its privileges in format map[privilege]given
func (p *PgStore) GetAuthData(id uint64) (string, string, map[string]bool, error) {
	user := User{}

	row := p.connPool.QueryRow("SELECT email, password, privileges "+
		"FROM users WHERE id = $1", id)

	err := row.Scan(&user.Email, &user.Password, &user.Privileges)
	if err != nil {
		return "", "", nil, errors.Wrap(err, "failed to scan user in GetAuthData")
	}

	return user.Email, user.Password, user.Privileges, nil
}

// Delete deletes user, given by its ID, from the postgres database
func (p *PgStore) Delete(id uint64) error {
	commandTag, err := p.connPool.Exec("DELETE FROM users "+
		"WHERE id = $1",
		id)

	if err != nil {
		return errors.Wrap(err, "failed to delete user")
	}

	if commandTag.RowsAffected() != 1 {
		log.Printf("[DEBUG] didn't found record to delete")
	}

	return nil
}

// Update updates user information stored in the database (everything except ID, email and password)
func (p *PgStore) Update(user User) error {
	tx, err := p.connPool.Begin()
	if err != nil {
		return errors.Wrap(err, "failed to start update transaction")
	}

	defer func() {
		errNested := tx.Rollback()
		if errNested != nil && errNested != pgx.ErrTxClosed {
			log.Printf("[ERROR] failed to rollback the transaction (Update in pgUsr): %+v", err)
		}
	}()

	commandTag, err := p.connPool.Exec("UPDATE users SET "+
		"privileges = $1, updated_at = $2 "+
		"WHERE id = $3",
		user.Privileges,
		time.Now(),
		user.ID)

	if err != nil {
		return errors.Wrap(err, "failed to update user")
	}

	if commandTag.RowsAffected() != 1 {
		log.Printf("[DEBUG] didn't found record to update")
	}

	err = tx.Commit()

	if err != nil {
		return errors.Wrap(err, "failed to commit transaction while updating user")
	}

	return nil
}
