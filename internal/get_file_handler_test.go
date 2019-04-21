package internal

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetFile(t *testing.T) {

	s := newServer("testing.com", 80, &fakeTools{
		storePath: "testdata",
	})

	req := httptest.NewRequest(http.MethodGet, "/image.jpg", nil)
	rec := httptest.NewRecorder()
	s.srv.ServeHTTP(rec, req)

	require.Equal(t, 200, rec.Code)
}
