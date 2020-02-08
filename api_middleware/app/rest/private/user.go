package private

import (
	"github.com/Semior001/mdcd-travelhack/app/rest/http_errors"
	"github.com/Semior001/mdcd-travelhack/app/store/user"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	R "github.com/go-pkgz/rest"
	"net/http"
	"strconv"
)

const hardBodyLimit = 1024 * 64

type UserController struct {
	ServiceUsr user.Service
}

type UserRest interface {
	GetUsers(w http.ResponseWriter, r *http.Request)
	GetUserById(w http.ResponseWriter, r *http.Request)
	PostUser(w http.ResponseWriter, r *http.Request)
	UpdateUser(w http.ResponseWriter, r *http.Request)
	DeleteUser(w http.ResponseWriter, r *http.Request)
}

func (u UserController) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := u.ServiceUsr.GetUsers()
	if err != nil {
		http_errors.SendJSONError(w, r, http.StatusInternalServerError, err, "", http_errors.ErrUserNotFound)
		return
	}
	render.JSON(w, r, users)
}

func (u UserController) GetUserById(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 0)
	if err != nil {
		http_errors.SendJSONError(w, r, http.StatusInternalServerError, err, "", http_errors.ErrUserNotFound)
		return
	}
	userInfo, err := u.ServiceUsr.GetUser(userId)
	if err != nil {
		http_errors.SendJSONError(w, r, http.StatusInternalServerError, err, "", http_errors.ErrUserNotFound)
		return
	}

	render.JSON(w, r, userInfo)
}

// todo add correct privilege
func (u UserController) PostUser(w http.ResponseWriter, r *http.Request) {
	var newUser user.User
	err := render.DecodeJSON(http.MaxBytesReader(w, r.Body, hardBodyLimit), &newUser)
	if err != nil {
		http_errors.SendJSONError(w, r, http.StatusInternalServerError, err, "", http_errors.ErrUserNotFound)
		return
	}

	userId, err := u.ServiceUsr.PutUser(newUser)
	if err != nil {
		http_errors.SendJSONError(w, r, http.StatusInternalServerError, err, "", http_errors.ErrUserNotFound)
		return
	}
	render.JSON(w, r, R.JSON{
		"ID": userId,
	})
}

func (u UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var edit user.User
	err := render.DecodeJSON(http.MaxBytesReader(w, r.Body, hardBodyLimit), &edit)
	if err != nil {
		http_errors.SendJSONError(w, r, http.StatusInternalServerError, err, "", http_errors.ErrUserNotFound)
		return
	}

	userId, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 0)
	if err != nil {
		http_errors.SendJSONError(w, r, http.StatusInternalServerError, err, "", http_errors.ErrUserNotFound)
		return
	}
	userInfo, err := u.ServiceUsr.GetUser(userId)
	if err != nil {
		http_errors.SendJSONError(w, r, http.StatusInternalServerError, err, "", http_errors.ErrUserNotFound)
		return
	}
	edit.ID = userId
	if edit.Email == "" {
		edit.Email = userInfo.Email
	}
	if edit.Password == "" {
		edit.Password = userInfo.Password
	}
	if edit.Privileges == nil {
		edit.Privileges = userInfo.Privileges
	}

	err = u.ServiceUsr.UpdateUser(edit)
	if err != nil {
		http_errors.SendJSONError(w, r, http.StatusInternalServerError, err, "", http_errors.ErrUserNotFound)
		return
	}
}

func (u UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 0)
	if err != nil {
		http_errors.SendJSONError(w, r, http.StatusInternalServerError, err, "", http_errors.ErrUserNotFound)
		return
	}
	err = u.ServiceUsr.DeleteUser(userId)
	if err != nil {
		http_errors.SendJSONError(w, r, http.StatusInternalServerError, err, "", http_errors.ErrUserNotFound)
		return
	}
	render.JSON(w, r, R.JSON{
		"message": "confirmed",
	})
}
