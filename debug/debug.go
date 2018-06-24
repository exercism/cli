package debug

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
)

var (
	// Verbose determines if debugging output is displayed to the user
	Verbose bool
	output  io.Writer = os.Stderr
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
	req.Body = ioutil.NopCloser(body)

	dump, err := httputil.DumpRequest(req, req.ContentLength > 0)
	if err != nil {
		log.Fatal(err)
	}

	Println("\n========================= BEGIN DumpRequest =========================")
	Println(string(dump))
	Println("========================= END DumpRequest =========================")
	Println("")

	req.Body = ioutil.NopCloser(&bodyCopy)
}

// DumpResponse dumps out the provided http.Response
func DumpResponse(res *http.Response) {
	if !Verbose {
		return
	}

	var bodyCopy bytes.Buffer
	body := io.TeeReader(res.Body, &bodyCopy)
	res.Body = ioutil.NopCloser(body)

	dump, err := httputil.DumpResponse(res, res.ContentLength > 0)
	if err != nil {
		log.Fatal(err)
	}

	Println("\n========================= BEGIN DumpResponse =========================")
	Println(string(dump))
	Println("========================= END DumpResponse =========================")
	Println("")

	res.Body = ioutil.NopCloser(body)
}
