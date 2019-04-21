package processing

import (
	"context"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/stretchr/testify/require"
)

type toolsMock struct {
	w mockWriter
}

type mockWriter struct {
	Error error
}

func (m *mockWriter) Write(aSrc io.Reader, aStorePath string) (string, error) {
	if m.Error == nil {
		return "filename.png", nil
	}
	return "", m.Error
}

func (m *mockWriter) CreatePreview(aSrcPath, aDstPath string) error {
	if m.Error == nil {
		return nil
	}
	return m.Error
}

func (t *toolsMock) SourceWriter() ImageWriter {
	return &t.w
}

func (t *toolsMock) PreviewWriter() PreviewWriter {
	return &t.w
}

func (t *toolsMock) StorePath() string {
	return ""
}

func TestInvalidNewURLProcessor(t *testing.T) {

	p, err := NewURLProcessor("InvalidUrl", &toolsMock{})

	require.Error(t, err)
	require.Nil(t, p)
}

func TestURLProcessor(t *testing.T) {
	e := echo.New()
	e.Use(middleware.Static("../testdata"))

	go func() {
		e.Start("localhost:8099")
	}()

	// Success story
	p, err := NewURLProcessor("http://localhost:8099/image.jpg", &toolsMock{})

	require.NoError(t, err)

	id, err := p.Process()

	require.NoError(t, err)
	require.Equal(t, "filename.png", id)

	// Error story
	p, err = NewURLProcessor("http://localhost:8099/image.jpg", &toolsMock{
		w: mockWriter{
			Error: fmt.Errorf("Some error"),
		},
	})

	id, err = p.Process()

	require.Error(t, err)
	require.Empty(t, id)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	e.Shutdown(ctx)
}
