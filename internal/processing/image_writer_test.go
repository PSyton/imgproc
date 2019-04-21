package processing

import (
	"os"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestImageWriterUnsupported(t *testing.T) {

	req := require.New(t)

	src, err := os.Open("../testdata/img.json")
	req.NoError(err)
	w := newWriter()
	defer src.Close()

	filename, err := w.Write(src, "../testdata/1")

	req.EqualError(err, "unsupported mime type")
	req.Empty(filename)
}

func TestImageWriterNonExistPath(t *testing.T) {

	req := require.New(t)

	src, err := os.Open("../testdata/image.jpg")
	req.NoError(err)
	w := newWriter()
	defer src.Close()

	filename, err := w.Write(src, "../testdata/1")

	req.Error(err)
	req.Empty(filename)
}

func testFormat(t *testing.T, format string) {

	srcFile := "../testdata/image." + format
	dstPath := "../testdata"

	req := require.New(t)

	src, err := os.Open(srcFile)
	req.NoError(err)
	w := newWriter()
	defer src.Close()

	filename, err := w.Write(src, "../testdata")

	req.NoError(err)
	req.NotEmpty(filename)
	req.True(strings.HasSuffix(filename, format))

	dstFile := path.Join(dstPath, filename)

	srcInfo, err := os.Stat(srcFile)
	require.NoError(t, err)
	dstInfo, err := os.Stat(dstFile)
	require.NoError(t, err)
	require.Equal(t, srcInfo.Size(), dstInfo.Size())
	os.Remove(dstFile)
}

func TestImageWriter(t *testing.T) {
	testFormat(t, "png")
	testFormat(t, "jpg")
}
