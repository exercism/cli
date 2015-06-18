package api

import (
	"bytes"
	"errors"
	"io/ioutil"
	"path/filepath"
	"strings"

	"golang.org/x/net/html/charset"
	"golang.org/x/text/transform"
)

const (
	mimeType = "text/plain"
)

var (
	errUnidentifiable = errors.New("unable to identify track and problem")
	errNoFiles        = errors.New("no files submitted")
	utf8BOM           = []byte{0xef, 0xbb, 0xbf}
)

// Iteration represents a version of a particular exercise.
// This gets submitted to the API.
type Iteration struct {
	Key      string            `json:"key"`
	Code     string            `json:"code"`
	Dir      string            `json:"dir"`
	Language string            `json:"language"`
	Problem  string            `json:"problem"`
	Solution map[string]string `json:"solution"`
}

// NewIteration prepares an iteration of a problem in a track for submission to the API.
// It takes a dir and a list of files which it will read from disk.
// All paths are assumed to be absolute paths with symlinks resolved.
func NewIteration(dir string, filenames []string) (*Iteration, error) {
	if len(filenames) == 0 {
		return nil, errNoFiles
	}

	iter := &Iteration{
		Dir:      dir,
		Solution: map[string]string{},
	}

	// All the files should be within the exercism path.
	for _, filename := range filenames {
		if !iter.isValidFilepath(filename) {
			return nil, errUnidentifiable
		}
	}

	// Identify language track and problem slug.
	path := filenames[0][len(dir):]
	segments := strings.Split(path, string(filepath.Separator))
	if len(segments) < 4 {
		return nil, errUnidentifiable
	}
	iter.Language = segments[1]
	iter.Problem = segments[2]

	for _, filename := range filenames {
		fileContents, err := readFileAsUTF8String(filename)
		if err != nil {
			return nil, err
		}

		path := filename[len(iter.RelativePath()):]
		iter.Solution[path] = *fileContents
	}
	return iter, nil
}

// RelativePath returns the iterations relative path
// iter.Dir/iter.Language/iter.Problem/
func (iter *Iteration) RelativePath() string {
	return filepath.Join(iter.Dir, iter.Language, iter.Problem) + string(filepath.Separator)
}

func (iter *Iteration) isValidFilepath(path string) bool {
	if iter == nil {
		return false
	}
	return strings.HasPrefix(strings.ToLower(path), strings.ToLower(iter.Dir))
}

func readFileAsUTF8String(filename string) (*string, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	encoding, _, _ := charset.DetermineEncoding(b, mimeType)
	decoder := encoding.NewDecoder()
	decodedBytes, _, err := transform.Bytes(decoder, b)
	if err != nil {
		return nil, err
	}

	// Drop the UTF-8 BOM that may have been added. This isn't necessary, and
	// it's going to be written into another UTF-8 buffer anyway once it's JSON
	// serialized.
	//
	// The standard recommends omitting the BOM. See
	// http://www.unicode.org/versions/Unicode5.0.0/ch02.pdf
	decodedBytes = bytes.TrimPrefix(decodedBytes, utf8BOM)

	s := string(decodedBytes)
	return &s, nil
}
