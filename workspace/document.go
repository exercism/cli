package workspace

import (
	"os"
	"path/filepath"
	"strings"
)

// Document is a file in a directory.
type Document struct {
	Root     string
	Filepath string
}

// NewDocument creates a document from a filepath.
// The root is typically the root of the exercise, and
// path is the relative path to the file within the root directory.
func NewDocument(root, path string) Document {
	return Document{
		Root:     root,
		Filepath: path,
	}
}

// Path is the normalized path.
// It uses forward slashes regardless of the operating system.
func (doc Document) Path() string {
	path := strings.Replace(doc.Filepath, doc.Root, "", 1)
	path = strings.TrimLeft(path, string(os.PathSeparator))
	return filepath.ToSlash(path)
}
