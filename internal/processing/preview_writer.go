package processing

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"os"

	"golang.org/x/image/draw"
)

// PreviewWriter store image preview for image from source path
type PreviewWriter interface {
	CreatePreview(aSrcPath, aDstPath string) error
}

type previewWriter struct {
	sizeLimit int
}

// NewPreviewWriter create preview writer
func newPreviewWriter(aSizeLimit int) PreviewWriter {
	return &previewWriter{
		sizeLimit: aSizeLimit,
	}
}

func max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

func (pw *previewWriter) CreatePreview(aSrcPath, aDstPath string) error {
	srcImageFile, err := os.Open(aSrcPath)
	if err != nil {
		return err
	}
	defer srcImageFile.Close()
	src, format, err := image.Decode(srcImageFile)

	if err != nil {
		return err
	}

	bounds := src.Bounds()
	w, h := bounds.Dx(), bounds.Dy()

	if w <= 0 || h <= 0 {
		return fmt.Errorf("invalid image size")
	}

	dstImageFile, err := os.Create(aDstPath)
	if err != nil {
		return err
	}

	if max(w, h) <= pw.sizeLimit {
		srcImageFile.Seek(0, 0)
		_, err := io.Copy(dstImageFile, srcImageFile)
		if err != nil {
			dstImageFile.Close()
			os.Remove(aDstPath)
			return err
		}
		return nil
	}

	newW, newH := w*pw.sizeLimit/h, pw.sizeLimit
	if w > h {
		newW, newH = pw.sizeLimit, h*pw.sizeLimit/w
	}

	dst := image.NewRGBA(image.Rect(0, 0, newW, newH))

	// Slower than `draw.ApproxBiLinear.Scale()` but better quality.
	draw.BiLinear.Scale(dst, dst.Bounds(), src, src.Bounds(), draw.Src, nil)

	switch format {
	case "jpeg":
		err = jpeg.Encode(dstImageFile, dst, nil)
	case "png":
		err = png.Encode(dstImageFile, dst)
	default:
		err = fmt.Errorf("unsupported format")
	}

	if err != nil {
		dstImageFile.Close()
		os.Remove(aDstPath)
		return err
	}

	dstImageFile.Close()
	return nil
}
