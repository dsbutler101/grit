package cli

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError_Error(t *testing.T) {
	const exitCode = 1

	tt := map[string]struct {
		err      error
		expected string
	}{
		"nil error": {
			err:      nil,
			expected: fmt.Sprintf("unknown error (exit code %d)", exitCode),
		},
		"not nil error": {
			err:      assert.AnError,
			expected: assert.AnError.Error(),
		},
	}

	for tn, tc := range tt {
		t.Run(tn, func(t *testing.T) {
			assert.NotPanics(t, func() {
				assert.Equal(t, tc.expected, NewError(exitCode, tc.err).Error())
			})
		})
	}
}
