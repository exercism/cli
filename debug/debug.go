package debug

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
)

var (
	// Verbose determines if debugging output is displayed to the user
	Verbose bool
	output  io.Writer = os.Stderr
	// UnmaskAPIKey determines if the API key should de displayed during a dump
	UnmaskAPIKey bool
)

// Println conditionally outputs a message to Stderr
func Println(args ...interface{}) {
	if Verbose {
		fmt.Fprintln(output, args...)
	}
}

// Printf conditionally outputs a formatted message to Stderr
func Printf(format string, args ...interface{}) {
	if Verbose {
		fmt.Fprintf(output, format, args...)
	}
}

// DumpRequest dumps out the provided http.Request
func DumpRequest(req *http.Request) {
	if !Verbose {
		return
	}

	var bodyCopy bytes.Buffer
	body := io.TeeReader(req.Body, &bodyCopy)
	req.Body = io.NopCloser(body)

	authHeader := req.Header.Get("Authorization")

	if authParts := strings.Split(authHeader, " "); len(authParts) > 1 && !UnmaskAPIKey {
		if token := authParts[1]; token != "" {
			req.Header.Set("Authorization", "Bearer "+Redact(token))
		}
	}

	dump, err := httputil.DumpRequest(req, req.ContentLength > 0)
	if err != nil {
		log.Fatal(err)
	}

	Println("\n========================= BEGIN DumpRequest =========================")
	Println(string(dump))
	Println("========================= END DumpRequest =========================")
	Println("")

	req.Header.Set("Authorization", authHeader)
	req.Body = io.NopCloser(&bodyCopy)
}

// DumpResponse dumps out the provided http.Response
func DumpResponse(res *http.Response) {
	if !Verbose {
		return
	}

	var bodyCopy bytes.Buffer
	body := io.TeeReader(res.Body, &bodyCopy)
	res.Body = io.NopCloser(body)

	dump, err := httputil.DumpResponse(res, res.ContentLength > 0)
	if err != nil {
		log.Fatal(err)
	}

	Println("\n========================= BEGIN DumpResponse =========================")
	Println(string(dump))
	Println("========================= END DumpResponse =========================")
	Println("")

	res.Body = io.NopCloser(body)
}

// Redact masks the given token by replacing part of the string with *
func Redact(token string) string {
	str := token[4 : len(token)-3]
	redaction := strings.Repeat("*", len(str))
	return string(token[:4]) + redaction + string(token[len(token)-3:])
}
