package comms

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type thing struct {
	name   string
	rating int
}

func (t thing) String() string {
	return fmt.Sprintf("%s (+%d)", t.name, t.rating)
}

var (
	things = []thing{
		{name: "water", rating: 10},
		{name: "food", rating: 3},
		{name: "music", rating: 0},
	}
)

func TestSelectionDisplay(t *testing.T) {
	// We have to manually add each thing to the options collection.
	var sel Selection
	for _, thing := range things {
		sel.Items = append(sel.Items, thing)
	}

	display := "  [1] water (+10)\n  [2] food (+3)\n  [3] music (+0)\n"
	assert.Equal(t, display, sel.Display())
}

func TestSelectionGet(t *testing.T) {
	var sel Selection
	for _, thing := range things {
		sel.Items = append(sel.Items, thing)
	}

	_, err := sel.Get(0)
	assert.Error(t, err)

	o, err := sel.Get(1)
	assert.NoError(t, err)
	// We need to do a type assertion to access
	// any non-stringer stuff.
	t1 := o.(thing)
	assert.Equal(t, "water", t1.name)

	o, err = sel.Get(2)
	assert.NoError(t, err)
	t2 := o.(thing)
	assert.Equal(t, "food", t2.name)

	o, err = sel.Get(3)
	assert.NoError(t, err)
	t3 := o.(thing)
	assert.Equal(t, "music", t3.name)

	_, err = sel.Get(4)
	assert.Error(t, err)
}

func TestSelectionRead(t *testing.T) {
	var sel Selection
	n, err := sel.Read(strings.NewReader("5"))
	assert.NoError(t, err)
	assert.Equal(t, 5, n)

	_, err = sel.Read(strings.NewReader("abc"))
	assert.Error(t, err)
}

func TestSelectionPick(t *testing.T) {
	tests := []struct {
		desc      string
		selection Selection
		things    []thing
		expected  string
	}{
		{
			desc: "autoselect the only one",
			selection: Selection{
				// it never hits the error,
				// because it doesn't actually do
				// the prompt and read response.
				Reader: strings.NewReader("BOOM!"),
			},
			things: []thing{
				{"hugs", 100},
			},
			expected: "hugs",
		},
		{
			desc: "it picks the one corresponding to the selection",
			selection: Selection{
				Reader: strings.NewReader("2"),
			},
			things: []thing{
				{"food", 10},
				{"water", 3},
				{"music", 0},
			},
			expected: "water",
		},
	}

	for _, test := range tests {
		test.selection.Writer = ioutil.Discard
		for _, th := range test.things {
			test.selection.Items = append(test.selection.Items, th)
		}

		item, err := test.selection.Pick("which one? %s")
		assert.NoError(t, err)
		th, ok := item.(thing)
		assert.True(t, ok)
		assert.Equal(t, test.expected, th.name)
	}
}
