package comms

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQuestion(t *testing.T) {
	testCases := []struct {
		desc     string
		given    string
		fallback string
		expected string
	}{
		{"records interactive response", "hello\n", "", "hello"},
		{"responds with default if response is empty", "\n", "Fine.", "Fine."},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			q := &Question{
				Reader:       strings.NewReader(tc.given),
				Writer:       ioutil.Discard,
				Prompt:       "Say something: ",
				DefaultValue: tc.fallback,
			}

			answer, err := q.Ask()
			assert.NoError(t, err)
			assert.Equal(t, answer, tc.expected)
		})
	}
}
