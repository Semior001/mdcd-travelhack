package user

import (
	"github.com/go-pkgz/mongo"
)

// MongoStorage implements all storage methods, defined in Store
type MongoStorage struct {
	*mongo.Connection
}

func (m MongoStorage) Migrate(force bool) error {
	panic("implement me")
}

func (m MongoStorage) putUser(user User) (id uint64, err error) {
	panic("implement me")
}

func (m MongoStorage) UpdateUser(user User) (err error) {
	panic("implement me")
}

func (m MongoStorage) GetUser(id uint64) (user *User, err error) {
	panic("implement me")
}

func (m MongoStorage) GetUsers() (users []User, err error) {
	panic("implement me")
}

func (m MongoStorage) GetUserCredentials(email string) (user *User, err error) {
	panic("implement me")
}

func (m MongoStorage) getBasicUserInfo(id uint64) (user *User, err error) {
	panic("implement me")
}

func (m MongoStorage) DeleteUser(id uint64) (err error) {
	panic("implement me")
}
