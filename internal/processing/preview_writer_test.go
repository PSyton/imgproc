package processing

import (
	"image"
	"os"
	"testing"
	// register
	_ "image/jpeg"
	// register
	_ "image/png"

	"github.com/stretchr/testify/require"
)

func TestResizeImageErrors(t *testing.T) {

	w := newPreviewWriter(100)

	var err error

	err = w.CreatePreview("../testdata/nonexist.jpg", "../testdata/p_nonexist.jpg")
	require.Error(t, err)

	err = w.CreatePreview("../testdata/img.json", "../testdata/p_nonexist.jpg")
	require.Error(t, err)

	w = newPreviewWriter(700)
	err = w.CreatePreview("../testdata/image.jpg", "../testdata_1/p_nonexist.jpg")

	require.Error(t, err)
}

func TestResizeImageSuccess(t *testing.T) {
	var err error

	src := "../testdata/image.jpg"
	dst := "../testdata/preview_image.jpg"

	w := newPreviewWriter(200)

	err = w.CreatePreview(src, dst)
	require.NoError(t, err)

	dstImageFile, err := os.Open(dst)
	require.NoError(t, err)
	img, _, err := image.Decode(dstImageFile)
	require.NoError(t, err)
	dstImageFile.Close()
	require.NoError(t, err)
	require.Equal(t, 200, max(img.Bounds().Dx(), img.Bounds().Dy()))
	os.Remove(dst)

	w = newPreviewWriter(640)
	err = w.CreatePreview(src, dst)
	require.NoError(t, err)
	srcInfo, err := os.Stat(src)
	require.NoError(t, err)
	dstInfo, err := os.Stat(dst)
	require.NoError(t, err)
	require.Equal(t, srcInfo.Size(), dstInfo.Size())
	os.Remove(dst)
}
