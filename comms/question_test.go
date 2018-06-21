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
		{"removes trailing \\r in addition to trailing \\", "hello\r\n", "Fine.", "hello"},
		{"removes trailing white spaces", "hello  \n", "Fine.", "hello"},
		{"falls back to default value", "  \n", "Default", "Default"},
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
