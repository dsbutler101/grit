package terraform

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlags_Validate(t *testing.T) {
	tests := map[string]struct {
		flags       Flags
		assertError func(t *testing.T, err error)
	}{
		"both target flags not set": {
			assertError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, ErrNoTargetFlagSet)
			},
		},
		"both target flags set": {
			flags: Flags{
				Target:          "foo",
				TargetStateFile: "bar",
			},
			assertError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, ErrBothTargetFlagsSet)
			},
		},
		"Target flag set": {
			flags: Flags{
				Target: "foo",
			},
		},
		"TargetStateFlag flag set": {
			flags: Flags{
				TargetStateFile: "foo",
			},
		},
	}

	for tn, tt := range tests {
		t.Run(tn, func(t *testing.T) {
			err := tt.flags.Validate()

			if tt.assertError != nil {
				tt.assertError(t, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}
