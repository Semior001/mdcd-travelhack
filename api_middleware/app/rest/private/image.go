package private

import (
	"bytes"
	"github.com/Semior001/mdcd-travelhack/app/rest/http_errors"
	"github.com/Semior001/mdcd-travelhack/app/store/image"
	"github.com/Semior001/mdcd-travelhack/app/store/user"
	"github.com/go-chi/render"
	"github.com/go-pkgz/auth/token"
	R "github.com/go-pkgz/rest"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
)

// ImageController defines some parameters that are necessary to
// mount controller methods
type ImageController struct {
	ServiceImg          image.Service
	ServiceUsr          user.Service
	ImageProcServiceURL string
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

func (i *ImageController) CheckBarcode(w http.ResponseWriter, r *http.Request) {
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

// SaveImage stores image into the database
func (i *ImageController) SaveImage(w http.ResponseWriter, r *http.Request) {
	imgType := r.URL.Query().Get("imgType")
	barcode := r.URL.Query().Get("barcode")
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

	imgId, err := i.ServiceImg.PutImage(userId, barcode, imgType, reader)
	if err != nil {
		http_errors.SendJSONError(w, r, http.StatusInternalServerError, err, "", http_errors.ErrPutImage)
		return
	}
	render.JSON(w, r, R.JSON{
		"ID": imgId,
	})
}

func (i *ImageController) GetBackgrounds(w http.ResponseWriter, r *http.Request) {
	ids, err := i.ServiceImg.GetBackgrounds()
	if err != nil {
		http_errors.SendJSONError(w, r, http.StatusInternalServerError, err, "", http_errors.ErrInternal)
		return
	}
	render.JSON(w, r, ids)
}

func (i *ImageController) GetBackground(w http.ResponseWriter, r *http.Request) {
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

// GetImage returns image itself
func (i *ImageController) GetImage(w http.ResponseWriter, r *http.Request) {
	//imgId := r.URL.Query().Get("imgId")

}

// PostFilter commits filter and bg replacement and returns them to client
func (i *ImageController) PostFilter(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(20480)
	if err != nil {
		log.Printf("[WARN] failed to parse multipart form %+v", err)
		render.JSON(w, r, R.JSON{
			"error": "failed to parse multipart form",
		})
		render.Status(r, 500)
	}

	bgFiles := r.MultipartForm.File["background"]
	bgIdStr := r.MultipartForm.Value["background_id"][0]
	filterName := r.MultipartForm.Value["filter_name"][0]
	barcodeStr := r.MultipartForm.Value["barcode"][0]

	var bg io.Reader

	if len(bgFiles) == 1 {
		bg, err = bgFiles[0].Open()
		if err != nil {
			log.Printf("[WARN] failed to open bgfile %+v", err)
			render.Status(r, 500)
			render.JSON(w, r, R.JSON{
				"error": "failed to open bgfile",
			})
			return
		}
	}

	bgId, err := strconv.Atoi(bgIdStr)
	if err != nil {
		log.Printf("[WARN] failed to atoi %+v", err)
		render.Status(r, 500)
		render.JSON(w, r, R.JSON{
			"error": "failed to atoi",
		})
		return
	}

	if bgId != -1 {
		_, bg, err = i.ServiceImg.GetImage(uint64(bgId))
		if err != nil {
			log.Printf("[WARN] failed to load background by its id %+v", err)
			render.Status(r, 500)
			render.JSON(w, r, R.JSON{
				"error": "failed to load background by its id",
			})
			return
		}
	}

	_, img, err := i.ServiceImg.GetImageByBarcode(barcodeStr)
	if err != nil {
		log.Printf("[WARN] failed to get image %+v", err)
		render.Status(r, 500)
		render.JSON(w, r, R.JSON{
			"error": "failed to get image from db",
		})
		return
	}

	// replacing background
	if bg != nil {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		// writing image
		iowr, err := writer.CreateFormFile("image", "img.jpg")
		if err != nil {
			log.Printf("[WARN] failed to create image form field %+v", err)
			render.Status(r, 500)
			render.JSON(w, r, R.JSON{
				"error": "failed to create multipart body",
			})
			return
		}
		_, err = io.Copy(iowr, img)
		if err != nil {
			log.Printf("[WARN] failed to write image to multipart body %+v", err)
			render.Status(r, 500)
			render.JSON(w, r, R.JSON{
				"error": "failed to create multipart body",
			})
			return
		}

		iowr, err = writer.CreateFormFile("background", "bg.jpg")
		if err != nil {
			log.Printf("[WARN] failed to create multipart body %+v", err)
			render.Status(r, 500)
			render.JSON(w, r, R.JSON{
				"error": "failed to create multipart body",
			})
			return
		}
		_, err = io.Copy(iowr, bg)
		if err != nil {
			log.Printf("[WARN] failed to write background to multipart body %+v", err)
			render.Status(r, 500)
			render.JSON(w, r, R.JSON{
				"error": "failed to create multipart body",
			})
			return
		}

		// calling replace_background
		req, err := http.NewRequest("POST", i.ImageProcServiceURL+"replace-background", body)
		if err != nil {
			log.Printf("[WARN] failed to initialize request to replace background %+v", err)
			render.Status(r, 500)
			render.JSON(w, r, R.JSON{
				"error": "failed to create multipart body",
			})
			return
		}
		req.Header.Set("Content-Type", writer.FormDataContentType())
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("[WARN] failed to send request to replace background %+v", err)
			render.Status(r, 500)
			render.JSON(w, r, R.JSON{
				"error": "failed to send request to replace background",
			})
			return
		} else {
			body := &bytes.Buffer{}
			_, err := body.ReadFrom(resp.Body)
			if err != nil {
				log.Printf("[WARN] failed to read body from response from request to replace background %+v", err)
				render.Status(r, 500)
				render.JSON(w, r, R.JSON{
					"error": "failed to read body from response from request to replace background",
				})
				return
			}
			err = resp.Body.Close()
			if err != nil {
				log.Printf("[WARN] failed to close response body from request to replace background %+v", err)
				render.Status(r, 500)
				render.JSON(w, r, R.JSON{
					"error": "failed to close response body from request to replace background ",
				})
				return
			}
		}
	}

	// gathering new request to filters
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// writing image
	iowr, err := writer.CreateFormFile("image", "img.jpg")
	if err != nil {
		log.Printf("[WARN] failed to create image form field %+v", err)
		render.JSON(w, r, R.JSON{
			"error": "failed to create multipart body",
		})
		render.Status(r, 500)
		return
	}
	_, err = io.Copy(iowr, img)
	if err != nil {
		log.Printf("[WARN] failed to write image to multipart body %+v", err)
		render.JSON(w, r, R.JSON{
			"error": "failed to create multipart body",
		})
		render.Status(r, 500)
		return
	}
	err = writer.WriteField("filter_name", filterName)
	if err != nil {
		log.Printf("[WARN] failed to write filter name to multipart body %+v", err)
		render.JSON(w, r, R.JSON{
			"error": "failed to write filter name to multipart body",
		})
		render.Status(r, 500)
		return
	}
	// calling apply_filter
	req, err := http.NewRequest("POST", i.ImageProcServiceURL+"apply-filter", body)
	if err != nil {
		log.Printf("[WARN] failed to initialize request to replace background %+v", err)
		render.JSON(w, r, R.JSON{
			"error": "failed to create multipart body",
		})
		render.Status(r, 500)
		return
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[WARN] failed to send request to replace background %+v", err)
		render.JSON(w, r, R.JSON{
			"error": "failed to send request to replace background",
		})
		render.Status(r, 500)
		return
	} else {
		body := &bytes.Buffer{}
		_, err := body.ReadFrom(resp.Body)
		if err != nil {
			log.Printf("[WARN] failed to read body from response from request to replace background %+v", err)
			render.JSON(w, r, R.JSON{
				"error": "failed to read body from response from request to replace background",
			})
			render.Status(r, 500)
			return
		}
		err = resp.Body.Close()
		if err != nil {
			log.Printf("[WARN] failed to close response body from request to replace background %+v", err)
			render.JSON(w, r, R.JSON{
				"error": "failed to close response body from request to replace background ",
			})
			render.Status(r, 500)
			return
		}
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

	w.Header().Set("Content-Type", imgContentType(".jpg"))
	//w.Header().Set("Content-Length", strconv.Itoa(int()))
	w.WriteHeader(http.StatusOK)
	if _, err := io.Copy(w, resp.Body); err != nil {
		http_errors.SendJSONError(w, r, http.StatusInternalServerError, err, "", http_errors.ErrInternal)
		return
	}
}

func (i *ImageController) CommitImage(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}
