package processing

import (
	"fmt"
	"io"
	mathrand "math/rand"
	"net/http"
	"os"
	"path"
	"sync"
	"time"

	"github.com/oklog/ulid"
)

const (
	preLoadSize = 512
)

// ImageWriter store received image to file
// Create unique filename and check input data to be supported images
type ImageWriter interface {
	Write(aSrc io.Reader, aStorePath string) (string, error)
}

type imageWriter struct {
	entropy *mathrand.Rand
	lock    sync.Mutex
	outDir  string
}

func newWriter() ImageWriter {
	return &imageWriter{
		entropy: mathrand.New(mathrand.NewSource(time.Now().UnixNano())),
	}
}

func (iw *imageWriter) nextID() string {
	iw.lock.Lock()
	defer iw.lock.Unlock()

	return ulid.MustNew(ulid.Timestamp(time.Now()), iw.entropy).String()
}

// Write checks that aSrc contains image and then store data to file.
// If success - return store filename
func (iw *imageWriter) Write(aSrc io.Reader, aStorePath string) (string, error) {
	signData := make([]byte, preLoadSize)

	read, err := aSrc.Read(signData)
	if err != nil {
		return "", err
	}

	if read < preLoadSize {
		signData = signData[:read]
	}

	mimeType := http.DetectContentType(signData)

	fileName := iw.nextID()

	switch mimeType {
	case "image/jpeg":
		fileName = fmt.Sprintf("%s.jpg", fileName)
	case "image/png":
		fileName = fmt.Sprintf("%s.png", fileName)
	default:
		return "", fmt.Errorf("unsupported mime type")
	}

	storePath := path.Join(aStorePath, fileName)

	dst, err := os.Create(storePath)
	if err != nil {
		return "", err
	}

	// Write first 512 bytes to file
	if _, err = dst.Write(signData); err != nil {
		dst.Close()
		os.Remove(storePath)
		return "", err
	}

	// write rest data
	if _, err = io.Copy(dst, aSrc); err != nil {
		dst.Close()
		os.Remove(storePath)
		return "", err
	}

	dst.Close()
	return fileName, nil
}
