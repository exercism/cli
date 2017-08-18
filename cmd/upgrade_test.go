package cmd

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

type fakeCLI struct {
	UpToDate      bool
	UpgradeCalled bool
}

func (fc *fakeCLI) IsUpToDate() (bool, error) {
	return fc.UpToDate, nil
}

func (fc *fakeCLI) Upgrade() error {
	fc.UpgradeCalled = true
	return nil
}

func TestUpgrade(t *testing.T) {
	oldOut := out
	out = ioutil.Discard
	defer func() { out = oldOut }()

	tests := []struct {
		desc     string
		upToDate bool
		expected bool
	}{
		{
			desc:     "upgrade should be called for an outdated CLI",
			upToDate: false,
			expected: true,
		},
		{
			desc:     "upgrade should not be called for an already updated CLI",
			upToDate: true,
			expected: false,
		},
	}

	for _, test := range tests {
		fc := &fakeCLI{UpToDate: test.upToDate}

		err := updateCLI(fc)
		assert.NoError(t, err)
		assert.Equal(t, test.expected, fc.UpgradeCalled, test.desc)
	}
}
