package processing

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"path"
)

// JSONProcessor process image from json base64 field
type JSONProcessor struct {
	data  []byte
	tools Tools
}

// NewJSONProcessor load data from JSON and create processor for it
func NewJSONProcessor(aReader io.Reader, aTools Tools) (*JSONProcessor, error) {

	values := map[string]string{}

	if err := json.NewDecoder(aReader).Decode(&values); err != nil {
		return nil, err
	}

	b64Value, ok := values["image"]
	if !ok {
		return nil, fmt.Errorf("image key not found")
	}

	data, err := base64.StdEncoding.DecodeString(b64Value)

	if err != nil {
		return nil, err
	}

	return &JSONProcessor{
		data:  data,
		tools: aTools,
	}, nil
}

// Process handle data processing
func (p *JSONProcessor) Process() (string, error) {

	buff := bytes.NewBuffer(p.data)

	fileName, err := p.tools.SourceWriter().Write(buff, p.tools.StorePath())
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
