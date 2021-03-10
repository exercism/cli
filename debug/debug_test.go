package debug

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVerboseEnabled(t *testing.T) {
	b := &bytes.Buffer{}
	output = b
	Verbose = true

	Println("World")
	if b.String() != "World\n" {
		t.Error("expected 'World' got", b.String())
	}
}

func TestVerboseDisabled(t *testing.T) {
	b := &bytes.Buffer{}
	output = b
	Verbose = false

	Println("World")
	if b.String() != "" {
		t.Error("expected '' got", b.String())
	}
}

func TestDumpRequest(t *testing.T) {
	testCases := []struct {
		desc    string
		auth    string
		verbose bool
		unmask  bool
	}{
		{
			desc:    "Do not attempt to dump request if 'Verbose' is set to false",
			auth:    "",
			verbose: false,
			unmask:  false,
		},
		{
			desc:    "Dump request without authorization header",
			auth:    "", //not set
			verbose: true,
			unmask:  false,
		},
		{
			desc:    "Dump request with malformed 'Authorization' header",
			auth:    "malformed",
			verbose: true,
			unmask:  true,
		},
		{
			desc:    "Dump request with properly formed 'Authorization' header",
			auth:    "Bearer abc12-345abcde1234-5abc12",
			verbose: true,
			unmask:  false,
		},
	}

	b := &bytes.Buffer{}
	output = b
	for _, tc := range testCases {
		Verbose = tc.verbose
		UnmaskAPIKey = tc.unmask
		r, _ := http.NewRequest("GET", "https://api.example.com/bogus", nil)
		if tc.auth != "" {
			r.Header.Set("Authorization", tc.auth)
		}

		DumpRequest(r)
		if tc.verbose {
			assert.Regexp(t, "GET /bogus", b.String(), tc.desc)
			assert.Equal(t, tc.auth, r.Header.Get("Authorization"), tc.desc)
			if tc.unmask {
				assert.Regexp(t, "Authorization: "+tc.auth, b.String(), tc.desc)
			}
		} else {
			assert.NotRegexp(t, "GET /bogus", b.String(), tc.desc)
		}
	}
}

func TestDumpResponse(t *testing.T) {
	b := &bytes.Buffer{}
	output = b
	Verbose = true
	r := &http.Response{
		StatusCode: 200,
		ProtoMajor: 1,
		ProtoMinor: 1,
	}

	DumpResponse(r)
	assert.Regexp(t, "HTTP/1.1 200 OK", b.String())
}

func TestRedact(t *testing.T) {
	fakeToken := "1a11111aaaa111aa1a11111a11111aa1"
	expected := "1a11*************************aa1"

	assert.Equal(t, expected, Redact(fakeToken))
}
