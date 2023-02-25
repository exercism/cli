package workspace

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizedDocumentPath(t *testing.T) {
	root, err := os.MkdirTemp("", "docpath")
	assert.NoError(t, err)
	defer os.RemoveAll(root)

	err = os.MkdirAll(filepath.Join(root, "subdirectory"), os.FileMode(0755))
	assert.NoError(t, err)

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
		err = os.WriteFile(tc.filepath, []byte("a file"), os.FileMode(0600))
		assert.NoError(t, err)

		doc, err := NewDocument(root, tc.filepath)
		assert.NoError(t, err)

		assert.Equal(t, doc.Path(), tc.path)
	}
}
