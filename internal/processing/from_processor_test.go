package processing

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFormProcessor(t *testing.T) {

	fileSrc, err := os.Open("../testdata/image.jpg")

	buf := new(bytes.Buffer)
	mr := multipart.NewWriter(buf)
	w, err := mr.CreateFormFile("image", "test.jpg")
	if assert.NoError(t, err) {
		io.Copy(w, fileSrc)
	}
	mr.Close()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/test", buf)
	req.Header.Set(echo.HeaderContentType, mr.FormDataContentType())
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	form, err := ctx.MultipartForm()
	require.NoError(t, err)
	files, ok := form.File["image"]
	require.True(t, ok)

	mock := toolsMock{}

	p, err := NewFormProcessor(files[0], &mock)
	require.NoError(t, err)
	require.NotNil(t, p)

	// Success story
	filename, err := p.Process()

	require.NoError(t, err)
	require.Equal(t, "filename.png", filename)

	// Error story
	mock.w.Error = fmt.Errorf("Some process error")
	filename, err = p.Process()

	require.Error(t, err)
	require.Empty(t, filename)
}
