package user

import (
	"errors"
	"testing"

	"github.com/exercism/cli/api"
	"github.com/exercism/cli/config"
)

func TestItemReport(t *testing.T) {
	testCases := []struct {
		max      int
		expected string
	}{
		{10, "Go (Clock) /tmp/go/clock\n"},
		{15, "Go (Clock)      /tmp/go/clock\n"},
		{25, "Go (Clock)                /tmp/go/clock\n"},
	}

	for _, tc := range testCases {
		problem1 := &api.Problem{
			TrackID:  "go",
			Slug:     "clock",
			Language: "Go",
			Name:     "Clock",
		}

		hw := NewHomework([]*api.Problem{problem1}, &config.Config{Dir: "/tmp"})
		if len(hw.Items) == 0 {
			t.Fatal(errors.New("failed to initialize homework correctly"))
		}
		item := hw.Items[0]
		actual := item.Report(hw.template, tc.max)
		if tc.expected != actual {
			t.Errorf("Expected:\n'%s', Got:\n'%s'\n", tc.expected, actual)
		}
	}
}
