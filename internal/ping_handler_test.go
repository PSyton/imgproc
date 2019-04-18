package internal

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/require"
)

func TestPing(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	srv := Server{}

	require.NoError(t, srv.pingHandler(ctx))
	require.Equal(t, http.StatusOK, rec.Code)
	require.Empty(t, rec.Body)
}
