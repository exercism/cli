package cmd

import (
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
	co := newCapturedOutput()
	co.override()
	defer co.reset()

	testCases := []struct {
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

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			fc := &fakeCLI{UpToDate: tc.upToDate}

			err := updateCLI(fc)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, fc.UpgradeCalled)
		})
	}
}
