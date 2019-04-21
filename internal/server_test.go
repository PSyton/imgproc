package internal

import (
	"net/http"
	"testing"

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
