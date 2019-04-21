package processing

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestJSONProcessor(t *testing.T) {
	req := require.New(t)

	src, err := os.Open("../testdata/img.json")
	req.NoError(err)

	mock := toolsMock{}

	p, err := NewJSONProcessor(src, &mock)
	req.NoError(err)
	req.True(len(p.data) > 0)

	// Success story
	filename, err := p.Process()

	require.NoError(t, err)
	require.Equal(t, "filename.png", filename)

	// Error story
	mock.w.Error = fmt.Errorf("Some process error")
	filename, err = p.Process()

	req.Error(err)
	req.Empty(filename)
}

func TestJSONProcessorWronJSON(t *testing.T) {
	req := require.New(t)

	buff := bytes.NewBufferString("{\"test\": 100}")

	// Error story
	p, err := NewJSONProcessor(buff, &toolsMock{})
	req.Error(err)
	req.Nil(p)
}
