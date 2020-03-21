package user

import (
	"encoding/json"
	"github.com/jackc/pgx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestPgStore_Delete(t *testing.T) {
	srv := preparePgStore(t)

	// inserting users
	usrs := map[string]*User{
		"foo@bar.com": {
			Email:    "foo@bar.com",
			Password: "blahblahblah",
			Privileges: map[string]bool{
				PrivilegeAdmin:       false,
				PrivilegeEditUsers:   false,
				PrivilegeInviteUsers: false,
			},
		},
		"foo1@bar.com": {
			Email:    "foo1@bar.com",
			Password: "blahblahblah1",
			Privileges: map[string]bool{
				PrivilegeAdmin:       true,
				PrivilegeEditUsers:   false,
				PrivilegeInviteUsers: false,
			},
		},
		"foo2@bar.com": {
			Email:    "foo2@bar.com",
			Password: "blahblahblah2",
			Privileges: map[string]bool{
				PrivilegeAdmin:       false,
				PrivilegeEditUsers:   true,
				PrivilegeInviteUsers: false,
			},
		},
		"foo3@bar.com": {
			Email:    "foo3@bar.com",
			Password: "blahblahblah",
			Privileges: map[string]bool{
				PrivilegeAdmin:       true,
				PrivilegeEditUsers:   true,
				PrivilegeInviteUsers: false,
			},
		},
		"foo4@bar.com": {
			Email:    "foo4@bar.com",
			Password: "blahblahblah",
			Privileges: map[string]bool{
				PrivilegeAdmin:       false,
				PrivilegeEditUsers:   false,
				PrivilegeInviteUsers: true,
			},
		},
		"foo5@bar.com": {
			Email:    "foo5@bar.com",
			Password: "blahblahblah",
			Privileges: map[string]bool{
				PrivilegeAdmin:       true,
				PrivilegeEditUsers:   true,
				PrivilegeInviteUsers: false,
			},
		},
	}
	for k := range usrs {
		privsMarshalled, err := json.Marshal(usrs[k].Privileges)

		tx, err := srv.connPool.Begin()
		require.NoError(t, err)

		row := tx.QueryRow("INSERT INTO "+
			"users(email, password, privileges) "+
			"VALUES ($1, $2, $3) "+
			"RETURNING id",
			usrs[k].Email,
			usrs[k].Password,
			privsMarshalled,
		)

		var id int
		err = row.Scan(&id)
		usrs[k].ID = id

		require.NoError(t, err)

		err = tx.Commit()
		require.NoError(t, err)
	}

	// querying
	err := srv.Delete(usrs["foo3@bar.com"].ID)
	require.NoError(t, err)

	// checking
	row := srv.connPool.QueryRow(`SELECT id, email, password, privileges FROM users WHERE id = $1`, usrs["foo3@bar.com"].ID)
	var id int
	var email string
	var pwd string
	var privsStr string
	err = row.Scan(&id, &email, &pwd, &privsStr)
	assert.Equal(t, err, pgx.ErrNoRows)
}

func TestPgStore_Get(t *testing.T) {
	srv := preparePgStore(t)

	// inserting users
	usrs := map[string]*User{
		"foo@bar.com": {
			Email:    "foo@bar.com",
			Password: "blahblahblah",
			Privileges: map[string]bool{
				PrivilegeAdmin:       false,
				PrivilegeEditUsers:   false,
				PrivilegeInviteUsers: false,
			},
		},
		"foo1@bar.com": {
			Email:    "foo1@bar.com",
			Password: "blahblahblah1",
			Privileges: map[string]bool{
				PrivilegeAdmin:       true,
				PrivilegeEditUsers:   false,
				PrivilegeInviteUsers: false,
			},
		},
		"foo2@bar.com": {
			Email:    "foo2@bar.com",
			Password: "blahblahblah2",
			Privileges: map[string]bool{
				PrivilegeAdmin:       false,
				PrivilegeEditUsers:   true,
				PrivilegeInviteUsers: false,
			},
		},
		"foo3@bar.com": {
			Email:    "foo3@bar.com",
			Password: "blahblahblah",
			Privileges: map[string]bool{
				PrivilegeAdmin:       true,
				PrivilegeEditUsers:   true,
				PrivilegeInviteUsers: false,
			},
		},
		"foo4@bar.com": {
			Email:    "foo4@bar.com",
			Password: "blahblahblah",
			Privileges: map[string]bool{
				PrivilegeAdmin:       false,
				PrivilegeEditUsers:   false,
				PrivilegeInviteUsers: true,
			},
		},
		"foo5@bar.com": {
			Email:    "foo5@bar.com",
			Password: "blahblahblah",
			Privileges: map[string]bool{
				PrivilegeAdmin:       true,
				PrivilegeEditUsers:   true,
				PrivilegeInviteUsers: false,
			},
		},
	}
	for k := range usrs {
		privsMarshalled, err := json.Marshal(usrs[k].Privileges)

		tx, err := srv.connPool.Begin()
		require.NoError(t, err)

		row := tx.QueryRow("INSERT INTO "+
			"users(email, password, privileges) "+
			"VALUES ($1, $2, $3) "+
			"RETURNING id",
			usrs[k].Email,
			usrs[k].Password,
			privsMarshalled,
		)

		var id int
		err = row.Scan(&id)
		usrs[k].ID = id

		require.NoError(t, err)

		err = tx.Commit()
		require.NoError(t, err)
	}

	// querying
	usr, err := srv.Get(usrs["foo3@bar.com"].ID)
	require.NoError(t, err)

	// checking
	shouldBe := usrs["foo3@bar.com"]
	assert.Equal(t, shouldBe.Email, usr.Email)
	assert.Equal(t, shouldBe.Password, usr.Password)

	ok := reflect.DeepEqual(shouldBe.Privileges, usr.Privileges)
	assert.True(t, ok)
}

func TestPgStore_GetAuthData(t *testing.T) {
	srv := preparePgStore(t)

	// inserting users
	usrs := map[string]*User{
		"foo@bar.com": {
			Email:    "foo@bar.com",
			Password: "blahblahblah",
			Privileges: map[string]bool{
				PrivilegeAdmin:       false,
				PrivilegeEditUsers:   false,
				PrivilegeInviteUsers: false,
			},
		},
		"foo1@bar.com": {
			Email:    "foo1@bar.com",
			Password: "blahblahblah1",
			Privileges: map[string]bool{
				PrivilegeAdmin:       true,
				PrivilegeEditUsers:   false,
				PrivilegeInviteUsers: false,
			},
		},
		"foo2@bar.com": {
			Email:    "foo2@bar.com",
			Password: "blahblahblah2",
			Privileges: map[string]bool{
				PrivilegeAdmin:       false,
				PrivilegeEditUsers:   true,
				PrivilegeInviteUsers: false,
			},
		},
		"foo3@bar.com": {
			Email:    "foo3@bar.com",
			Password: "blahblahblah",
			Privileges: map[string]bool{
				PrivilegeAdmin:       true,
				PrivilegeEditUsers:   true,
				PrivilegeInviteUsers: false,
			},
		},
		"foo4@bar.com": {
			Email:    "foo4@bar.com",
			Password: "blahblahblah",
			Privileges: map[string]bool{
				PrivilegeAdmin:       false,
				PrivilegeEditUsers:   false,
				PrivilegeInviteUsers: true,
			},
		},
		"foo5@bar.com": {
			Email:    "foo5@bar.com",
			Password: "blahblahblah",
			Privileges: map[string]bool{
				PrivilegeAdmin:       true,
				PrivilegeEditUsers:   true,
				PrivilegeInviteUsers: false,
			},
		},
	}
	for k := range usrs {
		privsMarshalled, err := json.Marshal(usrs[k].Privileges)

		tx, err := srv.connPool.Begin()
		require.NoError(t, err)

		row := tx.QueryRow("INSERT INTO "+
			"users(email, password, privileges) "+
			"VALUES ($1, $2, $3) "+
			"RETURNING id",
			usrs[k].Email,
			usrs[k].Password,
			privsMarshalled,
		)

		var id int
		err = row.Scan(&id)
		usrs[k].ID = id

		require.NoError(t, err)

		err = tx.Commit()
		require.NoError(t, err)
	}

	// querying
	email, pwd, privs, err := srv.GetAuthData(usrs["foo3@bar.com"].ID)
	require.NoError(t, err)

	// checking
	shouldBe := usrs["foo3@bar.com"]
	assert.Equal(t, shouldBe.Email, email)
	assert.Equal(t, shouldBe.Password, pwd)

	ok := reflect.DeepEqual(shouldBe.Privileges, privs)
	assert.True(t, ok)
}

func TestPgStore_GetByEmail(t *testing.T) {
	srv := preparePgStore(t)

	// inserting users
	usrs := map[string]*User{
		"foo@bar.com": {
			Email:    "foo@bar.com",
			Password: "blahblahblah",
			Privileges: map[string]bool{
				PrivilegeAdmin:       false,
				PrivilegeEditUsers:   false,
				PrivilegeInviteUsers: false,
			},
		},
		"foo1@bar.com": {
			Email:    "foo1@bar.com",
			Password: "blahblahblah1",
			Privileges: map[string]bool{
				PrivilegeAdmin:       true,
				PrivilegeEditUsers:   false,
				PrivilegeInviteUsers: false,
			},
		},
		"foo2@bar.com": {
			Email:    "foo2@bar.com",
			Password: "blahblahblah2",
			Privileges: map[string]bool{
				PrivilegeAdmin:       false,
				PrivilegeEditUsers:   true,
				PrivilegeInviteUsers: false,
			},
		},
		"foo3@bar.com": {
			Email:    "foo3@bar.com",
			Password: "blahblahblah",
			Privileges: map[string]bool{
				PrivilegeAdmin:       true,
				PrivilegeEditUsers:   true,
				PrivilegeInviteUsers: false,
			},
		},
		"foo4@bar.com": {
			Email:    "foo4@bar.com",
			Password: "blahblahblah",
			Privileges: map[string]bool{
				PrivilegeAdmin:       false,
				PrivilegeEditUsers:   false,
				PrivilegeInviteUsers: true,
			},
		},
		"foo5@bar.com": {
			Email:    "foo5@bar.com",
			Password: "blahblahblah",
			Privileges: map[string]bool{
				PrivilegeAdmin:       true,
				PrivilegeEditUsers:   true,
				PrivilegeInviteUsers: false,
			},
		},
	}
	for k := range usrs {
		privsMarshalled, err := json.Marshal(usrs[k].Privileges)

		tx, err := srv.connPool.Begin()
		require.NoError(t, err)

		row := tx.QueryRow("INSERT INTO "+
			"users(email, password, privileges) "+
			"VALUES ($1, $2, $3) "+
			"RETURNING id",
			usrs[k].Email,
			usrs[k].Password,
			privsMarshalled,
		)

		var id int
		err = row.Scan(&id)
		usrs[k].ID = id

		require.NoError(t, err)

		err = tx.Commit()
		require.NoError(t, err)
	}

	// querying
	usr, err := srv.GetByEmail("foo2@bar.com")
	require.NoError(t, err)

	// checking
	shouldBe := usrs["foo2@bar.com"]
	assert.Equal(t, usr.ID, shouldBe.ID)
	assert.Equal(t, usr.Email, shouldBe.Email)
	assert.Equal(t, usr.Password, shouldBe.Password)

	ok := reflect.DeepEqual(usr.Privileges, shouldBe.Privileges)
	assert.True(t, ok)
}

func TestPgStore_List(t *testing.T) {
	srv := preparePgStore(t)

	// inserting users
	usrs := map[string]*User{
		"foo@bar.com": {
			Email:    "foo@bar.com",
			Password: "blahblahblah",
			Privileges: map[string]bool{
				PrivilegeAdmin:       false,
				PrivilegeEditUsers:   false,
				PrivilegeInviteUsers: false,
			},
		},
		"foo1@bar.com": {
			Email:    "foo1@bar.com",
			Password: "blahblahblah1",
			Privileges: map[string]bool{
				PrivilegeAdmin:       true,
				PrivilegeEditUsers:   false,
				PrivilegeInviteUsers: false,
			},
		},
		"foo2@bar.com": {
			Email:    "foo2@bar.com",
			Password: "blahblahblah2",
			Privileges: map[string]bool{
				PrivilegeAdmin:       false,
				PrivilegeEditUsers:   true,
				PrivilegeInviteUsers: false,
			},
		},
		"foo3@bar.com": {
			Email:    "foo3@bar.com",
			Password: "blahblahblah",
			Privileges: map[string]bool{
				PrivilegeAdmin:       true,
				PrivilegeEditUsers:   true,
				PrivilegeInviteUsers: false,
			},
		},
		"foo4@bar.com": {
			Email:    "foo4@bar.com",
			Password: "blahblahblah",
			Privileges: map[string]bool{
				PrivilegeAdmin:       false,
				PrivilegeEditUsers:   false,
				PrivilegeInviteUsers: true,
			},
		},
		"foo5@bar.com": {
			Email:    "foo5@bar.com",
			Password: "blahblahblah",
			Privileges: map[string]bool{
				PrivilegeAdmin:       true,
				PrivilegeEditUsers:   true,
				PrivilegeInviteUsers: false,
			},
		},
	}
	for k := range usrs {
		privsMarshalled, err := json.Marshal(usrs[k].Privileges)

		tx, err := srv.connPool.Begin()
		require.NoError(t, err)

		row := tx.QueryRow("INSERT INTO "+
			"users(email, password, privileges) "+
			"VALUES ($1, $2, $3) "+
			"RETURNING id",
			usrs[k].Email,
			usrs[k].Password,
			privsMarshalled,
		)

		var id int
		err = row.Scan(&id)
		usrs[k].ID = id

		require.NoError(t, err)

		err = tx.Commit()
		require.NoError(t, err)
	}

	// checking
	usrsQueried, err := srv.List()
	require.NoError(t, err)

	for _, user := range usrsQueried {
		shouldBe, ok := usrs[user.Email]
		assert.True(t, ok)

		assert.Equal(t, shouldBe.ID, user.ID)

		ok = reflect.DeepEqual(shouldBe.Privileges, user.Privileges)
		assert.True(t, ok, "privileges does not match")

		assert.Equal(t, shouldBe.Password, user.Password)
	}
}

func TestPgStore_Update(t *testing.T) {
	srv := preparePgStore(t)

	user := User{
		Email:    "foo@bar.com",
		Password: "$2y$08$bEpqwi8ylxW9a1i8iQwV2OFs8tGKUjajbFRAGOSnsnWhubnjpcOzW",
		Privileges: map[string]bool{
			PrivilegeAdmin:       true,
			PrivilegeEditUsers:   false,
			PrivilegeInviteUsers: false,
		},
	}

	privsMarshalled, err := json.Marshal(user.Privileges)
	require.NoError(t, err)

	tx, err := srv.connPool.Begin()
	require.NoError(t, err)

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
	require.NoError(t, err)
	user.ID = id

	err = tx.Commit()
	require.NoError(t, err)

	// updating user
	user.Privileges[PrivilegeEditUsers] = true

	err = srv.Update(user)
	require.NoError(t, err)

	// load from db
	var uId int
	var uEmail, uPwd string
	var uPrivsMarshalled string
	row = srv.connPool.QueryRow(`SELECT id, email, password, privileges FROM users`)

	err = row.Scan(&uId, &uEmail, &uPwd, &uPrivsMarshalled)
	require.NoError(t, err)

	var uPrivs map[string]bool
	err = json.Unmarshal([]byte(uPrivsMarshalled), &uPrivs)
	require.NoError(t, err)

	assert.Equal(t, 1, uId)
	assert.Equal(t, "foo@bar.com", uEmail)
	assert.Equal(t, "$2y$08$bEpqwi8ylxW9a1i8iQwV2OFs8tGKUjajbFRAGOSnsnWhubnjpcOzW", uPwd)

	ok := reflect.DeepEqual(map[string]bool{
		PrivilegeAdmin:       true,
		PrivilegeEditUsers:   true,
		PrivilegeInviteUsers: false,
	}, uPrivs)
	assert.True(t, ok)

}

func TestPgStore_put(t *testing.T) {
	srv := preparePgStore(t)

	id, err := srv.put(User{
		Email:    "foo@bar.com",
		Password: "$2y$08$bEpqwi8ylxW9a1i8iQwV2OFs8tGKUjajbFRAGOSnsnWhubnjpcOzW",
		Privileges: map[string]bool{
			PrivilegeAdmin:       true,
			PrivilegeEditUsers:   false,
			PrivilegeInviteUsers: false,
		},
	})

	require.NoError(t, err)
	assert.Equal(t, 1, id)

	// load from db
	var uId int
	var uEmail, uPwd string
	var uPrivsMarshalled string
	row := srv.connPool.QueryRow(`SELECT id, email, password, privileges FROM users`)

	err = row.Scan(&uId, &uEmail, &uPwd, &uPrivsMarshalled)
	require.NoError(t, err)

	var uPrivs map[string]bool
	err = json.Unmarshal([]byte(uPrivsMarshalled), &uPrivs)
	require.NoError(t, err)

	assert.Equal(t, 1, uId)
	assert.Equal(t, "foo@bar.com", uEmail)
	assert.Equal(t, "$2y$08$bEpqwi8ylxW9a1i8iQwV2OFs8tGKUjajbFRAGOSnsnWhubnjpcOzW", uPwd)

	ok := reflect.DeepEqual(map[string]bool{
		PrivilegeAdmin:       true,
		PrivilegeEditUsers:   false,
		PrivilegeInviteUsers: false,
	}, uPrivs)
	assert.True(t, ok)

}

func preparePgStore(t *testing.T) *PgStore {
	connStr := os.Getenv("DB_TEST")

	connConf, err := pgx.ParseConnectionString(connStr)

	const toMilliseconds = 1e6
	connPool, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig:     connConf,
		MaxConnections: 2,
		AfterConnect:   nil,
		AcquireTimeout: 60 * toMilliseconds,
	})
	require.NoError(t, err)

	st := PgStore{
		ConnStr:  connStr,
		connPool: connPool,
	}

	require.NoError(t, err)

	cleanupStorage(t, st.connPool)
	t.Cleanup(func() {
		cleanupStorage(t, st.connPool)
	})

	return &st
}

func cleanupStorage(t *testing.T, pool *pgx.ConnPool) {
	tx, err := pool.Begin()
	require.NoError(t, err)
	_, err = tx.Exec(`TRUNCATE images CASCADE`)
	require.NoError(t, err)
	_, err = tx.Exec(`TRUNCATE users CASCADE`)
	require.NoError(t, err)
	_, err = tx.Exec(`ALTER SEQUENCE images_id_seq RESTART WITH 1`)
	require.NoError(t, err)
	_, err = tx.Exec(`ALTER SEQUENCE users_id_seq RESTART WITH 1`)
	err = tx.Commit()
	require.NoError(t, err)
}
