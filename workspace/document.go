package workspace

import "path/filepath"

// Document is a file in a directory.
type Document struct {
	Root         string
	RelativePath string
}

// NewDocument creates a document from the filepath.
// The root is typically the root of the exercise, and
// path is the absolute path to the file.
func NewDocument(root, path string) (Document, error) {
	path, err := filepath.Rel(root, path)
	if err != nil {
		return Document{}, err
	}
	return Document{
		Root:         root,
		RelativePath: path,
	}, nil
}

// Filepath is the absolute path to the document on the filesystem.
func (doc Document) Filepath() string {
	return filepath.Join(doc.Root, doc.RelativePath)
}

// Path is the normalized path.
// It uses forward slashes regardless of the operating system.
func (doc Document) Path() string {
	return filepath.ToSlash(doc.RelativePath)
}
