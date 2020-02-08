package rest

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSendHTMLErrorPage(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/error" {
			t.Log("http err request", r.URL)
			SendHTMLErrorPage(w, r, 500, errors.New("error 500"), "error details 123456", 987)
			return
		}
		w.WriteHeader(404)
	}))

	defer ts.Close()

	resp, err := http.Get(ts.URL + "/error")
	require.NoError(t, err)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Equal(t, 500, resp.StatusCode)

	assert.NotContains(t, string(body), `987`, "user html should not contain internal error code")
	assert.Contains(t, string(body), `error details 123456`)
	assert.Contains(t, string(body), `error 500`)
}

func TestSendJSONError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/error" {
			t.Log("http err request", r.URL)
			SendJSONError(w, r, 500, errors.New("error 500"), "error details 654321", 293)
			return
		}
		w.WriteHeader(404)
	}))

	defer ts.Close()

	resp, err := http.Get(ts.URL + "/error")
	require.NoError(t, err)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, func() string {
		j, _ := json.Marshal(errData{
			details: "error details 654321",
			error:   "error 500",
		})
		return string(j) + "\n"
	}(), string(body))
}