package processing

import (
	"fmt"
)

// Tools allow access to some shared tool and info
type Tools interface {
	SourceWriter() ImageWriter
	PreviewWriter() PreviewWriter
	StorePath() string
}

type imageTools struct {
	storePath     string
	imageWriter   ImageWriter
	previewWriter PreviewWriter
}

// NewTools create helpers
func NewTools(aStorePath string, aSizeLimit int) Tools {
	return &imageTools{
		storePath:     aStorePath,
		imageWriter:   newWriter(),
		previewWriter: newPreviewWriter(aSizeLimit),
	}
}

func (it *imageTools) SourceWriter() ImageWriter {
	return it.imageWriter
}

func (it *imageTools) PreviewWriter() PreviewWriter {
	return it.previewWriter
}

func (it *imageTools) StorePath() string {
	return it.storePath
}

func previewFilename(aFileName string) string {
	return fmt.Sprintf("preview_%s", aFileName)
}
