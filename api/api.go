package api

import (
	"bytes"
	"os"

	"golang.org/x/net/html/charset"
	"golang.org/x/text/transform"
)

const (
	mimeType = "text/plain"
)

var (
	utf8BOM = []byte{0xef, 0xbb, 0xbf}
)

func readFileAsUTF8String(filename string) (*string, error) {
	b, err := os.ReadFile(filename)
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
