package comms

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// Question provides an interactive session.
type Question struct {
	Reader       io.Reader
	Writer       io.Writer
	Prompt       string
	DefaultValue string
}

// Read reads the user's input.
func (q Question) Read(r io.Reader) (string, error) {
	reader := bufio.NewReader(r)
	s, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	s = strings.Trim(s, "\n")
	if s == "" {
		return q.DefaultValue, nil
	}
	return s, nil
}

// Ask displays the prompt, then records the response.
func (q *Question) Ask() (string, error) {
	fmt.Fprintf(q.Writer, q.Prompt)
	return q.Read(q.Reader)
}
