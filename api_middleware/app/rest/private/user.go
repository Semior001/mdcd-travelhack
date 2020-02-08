package private

import (
	"github.com/Semior001/mdcd-travelhack/app/rest/http_errors"
	"github.com/Semior001/mdcd-travelhack/app/store/user"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"net/http"
	"strconv"
)

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

func (u UserController) PostUser(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}

func (u UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}

func (u UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}
