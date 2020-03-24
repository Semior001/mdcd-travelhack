package user

import (
	"strings"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

//go:generate mockery -inpkg -name Store -case snake

const (
	// PrivilegeInviteUsers privilege to invite users via link
	PrivilegeInviteUsers = "invite_users"
	// PrivilegeEditUsers privilege to editing users
	PrivilegeEditUsers = "edit_users"
	// PrivilegeAdmin privilege gives all privileges above
	PrivilegeAdmin = "admin"
)

// User describes basic user
//
// Note: if privilege is not present as key in the map, golang will return false in such case
type User struct {
	ID         int
	Email      string
	Password   string          `json:"-"`
	Privileges map[string]bool // in format "privilege: given"
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// Store defines an interface to put and load users from the database
type Store interface {
	put(user User) (id int, err error)
	Update(user User) (err error)
	Get(id int) (user User, err error)
	List() (users []User, err error)
	GetByEmail(email string) (user User, err error)
	GetAuthData(id int) (email string, pwdHash string, privs map[string]bool, err error)
	Delete(id int) (err error)
}

// Service wraps Store interface providing methods that
// are needed for any Store implementation
type Service struct {
	BCryptCost int
	Store
}

// ServiceOpts defines options to create connection with storage
//
// ConnStr should be provided in format: dbdriver://uname:password@address:port/dbname?[param1=][&param2=][...etc]
type ServiceOpts struct {
	ConnStr    string
	BCryptCost int
}

// NewService creates a new user service with specified parameters and returns it
func NewService(opts ServiceOpts) (*Service, error) {
	driverEndIdx := strings.Index(opts.ConnStr, "://")
	driver := opts.ConnStr[0:driverEndIdx]

	var db Store
	var err error

	switch driver {
	case "postgres":
		db, err = NewPgStore(opts.ConnStr)
	}

	if err != nil {
		return nil, errors.Wrapf(err, "failed to initialize user Store with connstr %s", opts.ConnStr)
	}

	return &Service{
		Store:      db,
		BCryptCost: opts.BCryptCost,
	}, nil
}

// CheckUserCredentials matches given user password with the stored (by email) hash
func (s *Service) CheckUserCredentials(email string, password string) (bool, error) {
	user, err := s.GetByEmail(email)
	if err != nil {
		return false, errors.Wrapf(err, "email %s not listed in db", email)
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err == bcrypt.ErrMismatchedHashAndPassword {
		return false, nil
	}

	return err == nil, err
}

// HashPwdAndPut hashes user password before saving it to the database
func (s *Service) HashPwdAndPut(user User) (int, error) {
	// hashing password
	bytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), s.BCryptCost)
	if err != nil {
		return 0, errors.Wrapf(err, "unable to hash given password")
	}
	user.Password = string(bytes)
	return s.put(user)
}
