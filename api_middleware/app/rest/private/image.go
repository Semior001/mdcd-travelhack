package private

import (
	"github.com/Semior001/mdcd-travelhack/app/rest/http_errors"
	"github.com/Semior001/mdcd-travelhack/app/store/image"
	"github.com/Semior001/mdcd-travelhack/app/store/user"
	"github.com/go-chi/render"
	"github.com/go-pkgz/auth/token"
	R "github.com/go-pkgz/rest"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
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
	GetBackgrounds(w http.ResponseWriter, r *http.Request)
	GetBackground(w http.ResponseWriter, r *http.Request)

	CheckBarcode(w http.ResponseWriter, r *http.Request)
}

func (i ImageController) CheckBarcode(w http.ResponseWriter, r *http.Request) {
	sctoken := r.URL.Query().Get("sctoken")
	if sctoken != "admin_access" {
		render.JSON(w, r, R.JSON{"ok": false})
		return
	}
	barcode := r.URL.Query().Get("barcode")

	json, err := i.ServiceImg.CheckBarcode(barcode)
	if err != nil {
		http_errors.SendJSONError(w, r, http.StatusInternalServerError, err, "", http_errors.ErrInternal)
		return
	}

	// todo implement call to printsrv

	render.JSON(w, r, json)
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

func (i ImageController) GetBackgrounds(w http.ResponseWriter, r *http.Request) {
	ids, err := i.ServiceImg.GetBackgrounds()
	if err != nil {
		http_errors.SendJSONError(w, r, http.StatusInternalServerError, err, "", http_errors.ErrInternal)
		return
	}
	render.JSON(w, r, ids)
}

func (i ImageController) GetBackground(w http.ResponseWriter, r *http.Request) {
	imgId, err := strconv.ParseUint(r.URL.Query().Get("id"), 10, 0)
	if err != nil {
		http_errors.SendJSONError(w, r, http.StatusInternalServerError, err, "", http_errors.ErrInternal)
		return
	}

	_, fh, err := i.ServiceImg.GetImage(imgId)
	if err != nil {
		http_errors.SendJSONError(w, r, http.StatusInternalServerError, err, "", http_errors.ErrInternal)
		return
	}

	imgContentType := func(img string) string {
		img = strings.ToLower(img)
		switch {
		case strings.HasSuffix(img, ".png"):
			return "image/png"
		case strings.HasSuffix(img, ".jpg") || strings.HasSuffix(img, ".jpeg"):
			return "image/jpeg"
		case strings.HasSuffix(img, ".gif"):
			return "image/gif"
		}
		return "image/*"

	}

	defer func() {
		if err := fh.Close(); err != nil {
			log.Printf("[DEBUG] can't close image file")
		}
	}()

	w.Header().Set("Content-Type", imgContentType(".jpg"))
	//w.Header().Set("Content-Length", strconv.Itoa(int()))
	w.WriteHeader(http.StatusOK)
	if _, err = io.Copy(w, fh); err != nil {
		http_errors.SendJSONError(w, r, http.StatusInternalServerError, err, "", http_errors.ErrInternal)
		return
	}
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
