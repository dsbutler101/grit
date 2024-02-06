package name

import (
	"testing"

	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/test_tools"
)

func TestName(t *testing.T) {
	testCases := map[string]struct {
		vars    map[string]any
		wantErr bool
	}{
		"short name is fine": {
			vars: map[string]any{
				"name": "1234",
			},
		},
		"long name gives an error": {
			vars: map[string]any{
				"name": "0123456789012", // 13 characters
			},
			wantErr: true,
		},
		"exactly the limit is okay": {
			vars: map[string]any{
				"name": "012345678901", // 12 characters
			},
		},
	}

	for tn, tc := range testCases {
		t.Run(tn, func(t *testing.T) {
			test_tools.PlanAndAssertError(t, tc.vars, tc.wantErr)
		})
	}
}
