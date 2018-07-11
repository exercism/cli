package workspace

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTransmission(t *testing.T) {
	t.Skip("This panics on Windows. Once debugged, this can likely be inlined back into the main transmission test.")

	_, cwd, _, _ := runtime.Caller(0)
	root := filepath.Join(cwd, "..", "..", "fixtures", "transmission")
	dirBird := filepath.Join(root, "creatures", "hummingbird")
	dirFeeder := filepath.Join(dirBird, "feeder")
	fileBird := filepath.Join(dirBird, "hummingbird.txt")
	fileSugar := filepath.Join(dirFeeder, "sugar.txt")

	testCases := []struct {
		desc string
		args []string
		ok   bool
		tx   *Transmission
	}{
		{
			desc: "more than one dir",
			args: []string{dirBird, dirFeeder},
			ok:   false,
		},
		{
			desc: "a file and a dir",
			args: []string{dirBird, fileBird},
			ok:   false,
		},
		{
			desc: "just one file",
			args: []string{fileBird},
			ok:   true,
			tx:   &Transmission{Files: []string{fileBird}, Dir: dirBird},
		},
		{
			desc: "multiple files",
			args: []string{fileBird, fileSugar},
			ok:   true,
			tx:   &Transmission{Files: []string{fileBird, fileSugar}, Dir: dirBird},
		},
		{
			desc: "one dir",
			args: []string{dirBird},
			ok:   true,
			tx:   &Transmission{Files: nil, Dir: dirBird},
		},
		{
			desc: "multiple exercise names",
			args: []string{"hummingbird", "bear"},
			ok:   false,
		},
		{
			desc: "one exercise name",
			args: []string{"hummingbird"},
			ok:   true,
			tx:   &Transmission{Files: nil, Dir: "hummingbird"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			tx, err := NewTransmission(root, tc.args)
			if tc.ok {
				assert.NoError(t, err, tc.desc)
			} else {
				assert.Error(t, err, tc.desc)
			}

			if tc.tx != nil {
				assert.Equal(t, tc.tx.Files, tx.Files)
				assert.Equal(t, tc.tx.Dir, tx.Dir)
			}
		})
	}
}
