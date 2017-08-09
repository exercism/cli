package workspace

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// Solutions is a collection of solutions to interactively choose from.
type Solutions []*Solution

// Pick lets a user interactively select a solution from a list.
func (sx Solutions) Pick(prompt string) (*Solution, error) {
	// If there's just one, then we're done here.
	if len(sx) == 1 {
		return sx[0], nil
	}

	fmt.Printf(prompt, sx.Display())

	n, err := sx.ReadSelection(os.Stdin)
	if err != nil {
		return &Solution{}, err
	}

	s, err := sx.Get(n)
	if err != nil {
		return &Solution{}, err
	}
	return s, nil
}

// Display shows a numbered list of the solutions to choose from.
// The list starts at 1, since that seems better in a user interface.
func (sx Solutions) Display() string {
	str := ""
	for i, s := range sx {
		str += fmt.Sprintf("  [%d] %s\n", i+1, s)
	}
	return str
}

// ReadSelection reads the user's selection and converts it to a number.
func (sx Solutions) ReadSelection(r io.Reader) (int, error) {
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
func (sx Solutions) Get(n int) (*Solution, error) {
	if n <= 0 || n > len(sx) {
		return &Solution{}, errors.New("can't do that")
	}
	return sx[n-1], nil
}
