package processing

import (
	"mime/multipart"
	"path"
)

// FormProcessor processes image from form field
type FormProcessor struct {
	file  *multipart.FileHeader
	tools Tools
}

// NewFormProcessor store file data and create processor for it
func NewFormProcessor(aFile *multipart.FileHeader, aTools Tools) (*FormProcessor, error) {
	return &FormProcessor{
		file:  aFile,
		tools: aTools,
	}, nil
}

// Process handle data processing
func (p *FormProcessor) Process() (string, error) {
	src, err := p.file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	fileName, err := p.tools.SourceWriter().Write(src, p.tools.StorePath())
	if err != nil {
		return "", err
	}

	imagePath := path.Join(p.tools.StorePath(), fileName)
	previewPath := path.Join(p.tools.StorePath(), previewFilename(fileName))

	if err = p.tools.PreviewWriter().CreatePreview(imagePath, previewPath); err != nil {
		return "", err
	}

	return fileName, nil
}
