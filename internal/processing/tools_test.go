package processing

import (
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestPrevieFileName(t *testing.T) {
	require.True(t, strings.HasPrefix(previewFilename("xxxx"), "preview_"))
}
