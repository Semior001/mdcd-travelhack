package user

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestService_HashPwdAndPut(t *testing.T) {
	store := &MockStore{}
	var storedUsr *User = nil
	store.On(
		"put",
		mock.MatchedBy(func(user User) bool {
			assert.Nil(t, storedUsr)
			storedUsr = &user
			return true
		}),
	).Return(uint64(1), nil)

	srv := Service{
		BCryptCost: 8,
		Store:      store,
	}
	id, err := srv.HashPwdAndPut(User{
		ID:       0,
		Email:    "foo@bar.com",
		Password: "qwerty12345",
	})
	require.NoError(t, err)
	assert.Equal(t, uint64(1), id)
	err = bcrypt.CompareHashAndPassword([]byte(storedUsr.Password), []byte("qwerty12345"))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		require.Fail(t, "hash does not match with pwd")
	}
	require.NoError(t, err)
}

func TestService_CheckUserCredentials(t *testing.T) {
	store := &MockStore{}
	store.On(
		"GetByEmail",
		mock.Anything,
	).Return(User{
		Email:    "foo@bar.com",
		Password: "$2y$08$bEpqwi8ylxW9a1i8iQwV2OFs8tGKUjajbFRAGOSnsnWhubnjpcOzW",
	}, nil)

	srv := Service{
		BCryptCost: 8,
		Store:      store,
	}

	// correct credentials
	ok, err := srv.CheckUserCredentials("foo@bar.com", "qwerty12345")
	require.NoError(t, err)
	require.True(t, ok)

	// false credentials
	ok, err = srv.CheckUserCredentials("foo@bar.com", "blahblahwrong")
	require.NoError(t, err)
	require.False(t, ok)

}
