package internal

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"imgproc/internal/processing"
)

type fakeTools struct {
	storePath string
}

func (f *fakeTools) SourceWriter() processing.ImageWriter {
	return nil
}

func (f *fakeTools) PreviewWriter() processing.PreviewWriter {
	return nil
}

func (f *fakeTools) StorePath() string {
	return f.storePath
}

func TestNewServer(t *testing.T) {
	s := newServer("testing.com", 80, &fakeTools{})
	require.Equal(t, "testing.com:80", s.address)

	routes := s.srv.Routes()
	require.NotEmpty(t, routes)

	tbl := []struct {
		path   string
		method string
	}{
		{
			path:   "/ping",
			method: http.MethodGet,
		},
		{
			path:   "/process",
			method: http.MethodPost,
		},
	}

	var routesCount = 0
	for _, r := range routes {
		for _, data := range tbl {
			if data.path == r.Path && data.method == r.Method {
				routesCount++
				break
			}
		}
	}

	require.Equal(t, len(tbl), routesCount)
}

func checkExists(aFilePath string) bool {
	if _, err := os.Stat(aFilePath); os.IsNotExist(err) {
		return false
	}
	return true
}

func TestJSONRequest(t *testing.T) {
	s := newServer("testing.com", 80, processing.NewTools("testdata", 100))

	data, _ := ioutil.ReadFile("testdata/img.json")
	req := httptest.NewRequest(http.MethodPost, "/process", bytes.NewReader(data))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	s.srv.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)
	fName := rec.Body.String()
	require.NotEmpty(t, fName)

	f := path.Join("testdata", fName)
	require.True(t, checkExists(f))
	os.Remove(f)

	f = path.Join("testdata", fmt.Sprintf("preview_%s", fName))
	require.True(t, checkExists(f))
	os.Remove(f)
}

func TestBadRequests(t *testing.T) {
	s := newServer("testing.com", 80, processing.NewTools("testdata", 100))

	data, _ := ioutil.ReadFile("testdata/bad.json")
	req := httptest.NewRequest(http.MethodPost, "/process", bytes.NewReader(data))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	s.srv.ServeHTTP(rec, req)
	require.Equal(t, http.StatusBadRequest, rec.Code)

	req = httptest.NewRequest(http.MethodPost, "/process", nil)
	rec = httptest.NewRecorder()

	s.srv.ServeHTTP(rec, req)
	require.Equal(t, http.StatusBadRequest, rec.Code)

	formBuf := new(bytes.Buffer)
	mr := multipart.NewWriter(formBuf)
	w, err := mr.CreateFormField("username")
	if assert.NoError(t, err) {
		w.Write([]byte("Some user"))
	}
	mr.Close()

	req = httptest.NewRequest(http.MethodPost, "/process", formBuf)
	req.Header.Set(echo.HeaderContentType, mr.FormDataContentType())
	rec = httptest.NewRecorder()

	s.srv.ServeHTTP(rec, req)
	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestFormRequest(t *testing.T) {
	s := newServer("testing.com", 80, processing.NewTools("testdata", 100))

	fileSrc, err := os.Open("testdata/image.jpg")

	formBuf := new(bytes.Buffer)
	mr := multipart.NewWriter(formBuf)
	w, err := mr.CreateFormFile("image", "test.jpg")
	if assert.NoError(t, err) {
		io.Copy(w, fileSrc)
	}
	mr.Close()

	req := httptest.NewRequest(http.MethodPost, "/process", formBuf)
	req.Header.Set(echo.HeaderContentType, mr.FormDataContentType())
	rec := httptest.NewRecorder()

	s.srv.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)
	fName := rec.Body.String()
	require.NotEmpty(t, fName)

	f := path.Join("testdata", fName)
	require.True(t, checkExists(f))
	os.Remove(f)

	f = path.Join("testdata", fmt.Sprintf("preview_%s", fName))
	require.True(t, checkExists(f))
	os.Remove(f)
}
