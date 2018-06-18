package comms

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQuestion(t *testing.T) {
	tests := []struct {
		desc     string
		given    string
		fallback string
		expected string
	}{
		{"records interactive response", "hello\n", "", "hello"},
		{"responds with default if response is empty", "\n", "Fine.", "Fine."},
		{"removes trailing \\r in addition to trailing \\", "hello\r\n", "Fine.", "hello"},
	}
	for _, test := range tests {
		q := &Question{
			Reader:       strings.NewReader(test.given),
			Writer:       ioutil.Discard,
			Prompt:       "Say something: ",
			DefaultValue: test.fallback,
		}

		answer, err := q.Ask()
		assert.NoError(t, err)
		assert.Equal(t, answer, test.expected, test.desc)
	}
}
