package api

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"golang.org/x/net/html/charset"
	"golang.org/x/text/transform"
)

const (
	mimeType = "text/plain"
)

var (
	errNoFiles = errors.New("no files submitted")
	utf8BOM    = []byte{0xef, 0xbb, 0xbf}
)

var msgSubmitCalledFromWrongDir = `Unable to identify track and file.

It seems like you've tried to submit a solution file located outside of your
configured exercises directory.

Current directory:	{{ .Current }}
Configured directory:	{{ .Configured }}

Try re-running "exercism fetch". Then move your solution file to the correct
exercise directory for the problem you're working on. It should be somewhere
inside {{ .Configured }}

For example, to submit the JavaScript "hello-world.js" problem, run
"exercism submit hello-world.js" from this directory:

{{ .Configured }}{{ .Separator }}javascript{{ .Separator }}hello-world

You can see where exercism is looking for your files with "exercism debug".

`

var msgGenericPathError = `Bad path to exercise file.

You're trying to submit a solution file from inside your exercises directory,
but it looks like the directory structure is something that exercism doesn't
recognize as a valid file path.

First, make a copy of your solution file and save it outside of
{{ .Configured }}

Then, run "exercism fetch". Move your solution file back to the correct
exercise directory for the problem you're working on. It should be somewhere
inside {{ .Configured }}

If you are having trouble, you can file a GitHub issue at (https://github.com/exercism/exercism.io/issues)

`

// Iteration represents a version of a particular exercise.
// This gets submitted to the API.
type Iteration struct {
	Key      string            `json:"key"`
	Code     string            `json:"code"`
	Dir      string            `json:"dir"`
	TrackID  string            `json:"language"`
	Problem  string            `json:"problem"`
	Solution map[string]string `json:"solution"`
	Comment  string            `json:"comment,omitempty"`
}

// NewIteration prepares an iteration of a problem in a track for submission to the API.
// It takes a dir (from the global config) and a list of files which it will read from disk.
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
			// User has run exercism submit in the wrong directory.
			return nil, newIterationError(msgSubmitCalledFromWrongDir, iter.Dir)
		}
	}

	// Identify the language track and problem slug.
	path := filenames[0][len(dir):]

	segments := strings.Split(path, string(filepath.Separator))
	if len(segments) < 4 {
		// Submit was called from inside exercism directory, but the path
		// is still bad. Has the user modified their path in some way?
		return nil, newIterationError(msgGenericPathError, iter.Dir)
	}
	iter.TrackID = segments[1]
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

// RelativePath returns the iteration's relative path.
func (iter *Iteration) RelativePath() string {
	return filepath.Join(iter.Dir, iter.TrackID, iter.Problem) + string(filepath.Separator)
}

// isValidFilepath checks a files's absolute filepath and returns true if it is
// within the configured exercise directory.
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

	encoding, _, certain := charset.DetermineEncoding(b, mimeType)
	if !certain {
		// We don't want to use an uncertain encoding.
		// In particular, doing that may mangle UTF-8 files
		// that have only ASCII in their first 1024 bytes.
		// See https://github.com/exercism/cli/issues/309.
		// So if we're unsure, use UTF-8 (no transformation).
		s := string(b)
		return &s, nil
	}
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

// newIterationError executes an error message template to create a detailed
// message for the end user. An error type is returned.
func newIterationError(msgTemplate, configured string) error {
	buffer := bytes.NewBufferString("")
	t, err := template.New("iterErr").Parse(msgTemplate)
	if err != nil {
		return err
	}

	current, err := os.Getwd()
	if err != nil {
		return err
	}

	var pathData = struct {
		Current    string
		Configured string
		Separator  string
	}{
		current,
		configured,
		string(filepath.Separator),
	}

	t.Execute(buffer, pathData)
	msg := buffer.String()
	return errors.New(msg)
}
