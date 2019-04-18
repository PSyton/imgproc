package processing

import (
	"fmt"
	"net/http"
	"path"
	"strings"
)

// URLProcessor handle images form urls
type URLProcessor struct {
	url   string
	tools Tools
}

// NewURLProcessor check url and create processor for it
func NewURLProcessor(aURL string, aTools Tools) (*URLProcessor, error) {
	if !strings.HasPrefix(aURL, "http") {
		return nil, fmt.Errorf("Invalid url")
	}

	return &URLProcessor{
		url:   aURL,
		tools: aTools,
	}, nil
}

// Process handle data processing
func (p *URLProcessor) Process() (string, error) {
	resp, err := http.Get(p.url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	fileName, err := p.tools.SourceWriter().Write(resp.Body, p.tools.StorePath())
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
