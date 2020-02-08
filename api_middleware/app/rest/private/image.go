package private

import (
	"github.com/Semior001/mdcd-travelhack/app/rest/http_errors"
	"github.com/Semior001/mdcd-travelhack/app/store/image"
	"github.com/Semior001/mdcd-travelhack/app/store/user"
	"github.com/go-chi/render"
	"github.com/go-pkgz/auth/token"
	R "github.com/go-pkgz/rest"
	"net/http"
)

// ImageController defines some parameters that are necessary to
// mount controller methods
type ImageController struct {
	ServiceImg image.Service
	ServiceUsr user.Service
}

// ImageRest defines methods to mount to the web-server
type ImageRest interface {
	SaveImage(w http.ResponseWriter, r *http.Request)
	GetImage(w http.ResponseWriter, r *http.Request)
	PostFilter(w http.ResponseWriter, r *http.Request)
	CommitImage(w http.ResponseWriter, r *http.Request)
	//GetBackgrounds(w http.ResponseWriter, r *http.Request)
}

func (i ImageController) SaveImage(w http.ResponseWriter, r *http.Request) {
	imgType := r.URL.Query().Get("imgType")
	usrToken, err := token.GetUserInfo(r)
	if err != nil {
		http_errors.SendJSONError(w, r, http.StatusInternalServerError, err, "", http_errors.ErrDBStoring)
		return
	}
	email := usrToken.Name
	userCredentials, err := i.ServiceUsr.GetUserCredentials(email)
	if err != nil {
		http_errors.SendJSONError(w, r, http.StatusInternalServerError, err, "", http_errors.ErrUserNotFound)
		return
	}
	userId := userCredentials.ID
	err = r.ParseMultipartForm(20480)
	if err != nil {

	}
	fh := r.MultipartForm.File["image"]
	reader, err := fh[0].Open()

	imgId, err := i.ServiceImg.PutImage(userId, imgType, reader)
	if err != nil {
		http_errors.SendJSONError(w, r, http.StatusInternalServerError, err, "", http_errors.ErrPutImage)
		return
	}
	render.JSON(w, r, R.JSON{
		"ID": imgId,
	})
}

func (i ImageController) GetImage(w http.ResponseWriter, r *http.Request) {
	//imgId := r.URL.Query().Get("imgId")

}

func (i ImageController) PostFilter(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}

func (i ImageController) CommitImage(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}
