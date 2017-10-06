package comms

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// Selection wraps a list of items.
// It is used for interactive communication.
type Selection struct {
	Items  []fmt.Stringer
	Reader io.Reader
	Writer io.Writer
}

// NewSelection prepares an empty collection for interactive input.
func NewSelection() Selection {
	return Selection{
		Reader: os.Stdin,
		Writer: os.Stdout,
	}
}

// Pick lets a user interactively select an option from a list.
func (sel Selection) Pick(prompt string) (fmt.Stringer, error) {
	// If there's just one, then we're done here.
	if len(sel.Items) == 1 {
		return sel.Items[0], nil
	}

	fmt.Fprintf(sel.Writer, prompt, sel.Display())

	n, err := sel.Read(sel.Reader)
	if err != nil {
		return nil, err
	}

	o, err := sel.Get(n)
	if err != nil {
		return nil, err
	}
	return o, nil
}

// Display shows a numbered list of the solutions to choose from.
// The list starts at 1, since that seems better in a user interface.
func (sel Selection) Display() string {
	str := ""
	for i, item := range sel.Items {
		str += fmt.Sprintf("  [%d] %s\n", i+1, item)
	}
	return str
}

// Read reads the user's selection and converts it to a number.
func (sel Selection) Read(r io.Reader) (int, error) {
	reader := bufio.NewReader(r)
	text, _ := reader.ReadString('\n')
	n, err := strconv.Atoi(strings.TrimSpace(text))
	if err != nil {
		return 0, err
	}
	return n, nil
}

// Get returns the solution corresponding to the number.
// The list starts at 1, since that seems better in a user interface.
func (sel Selection) Get(n int) (fmt.Stringer, error) {
	if n <= 0 || n > len(sel.Items) {
		return nil, errors.New("we don't have that one")
	}
	return sel.Items[n-1], nil
}
