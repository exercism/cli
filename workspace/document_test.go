package workspace

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizedDocumentPath(t *testing.T) {
	root := filepath.Join("the", "root", "path", "the-track", "the-exercise")
	testCases := []struct {
		filepath string
		path     string
	}{
		{
			filepath: filepath.Join(root, "file.txt"),
			path:     "file.txt",
		},
		{
			filepath: filepath.Join(root, "subdirectory", "file.txt"),
			path:     "subdirectory/file.txt",
		},
	}

	for _, tc := range testCases {
		doc := NewDocument(root, tc.filepath)
		assert.Equal(t, tc.path, doc.Path())
	}
}
