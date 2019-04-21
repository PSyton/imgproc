package internal

import (
	"bytes"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateJSONProcessor(t *testing.T) {
	e := echo.New()
	data, _ := ioutil.ReadFile("testdata/img.json")
	req := httptest.NewRequest(http.MethodGet, "/test", bytes.NewReader(data))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	srv := Server{}

	p, err := srv.createProcessor(ctx)

	require.NoError(t, err)
	require.NotNil(t, p)
}

func TestCreateFormProcessor(t *testing.T) {
	buf := new(bytes.Buffer)
	mr := multipart.NewWriter(buf)
	w, err := mr.CreateFormFile("image", "image.jpg")
	if assert.NoError(t, err) {
		w.Write([]byte("somedata"))
	}
	mr.Close()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/test", buf)
	req.Header.Set(echo.HeaderContentType, mr.FormDataContentType())
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	srv := Server{}

	p, err := srv.createProcessor(ctx)

	require.NoError(t, err)
	require.NotNil(t, p)
}

func TestCreateURLProcessor(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/test?url=http://google.com/image.jpg", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	srv := Server{}

	p, err := srv.createProcessor(ctx)

	require.NoError(t, err)
	require.NotNil(t, p)
}
