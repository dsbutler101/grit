package isoneof

import (
	"testing"

	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/test_tools"
)

func TestIsOneOf(t *testing.T) {
	testCases := map[string]struct {
		vars    map[string]any
		wantErr bool
	}{
		"value is allowed": {
			vars: map[string]any{
				"value":   "current",
				"allowed": []any{"current", "v1.25"},
				"prefix":  "runner",
			},
		},
		"value is not allowed": {
			vars: map[string]any{
				"value":   "v1.28",
				"allowed": []any{"current", "v1.25"},
				"prefix":  "runner",
			},
			wantErr: true,
		},
		"validation is disabled": {
			vars: map[string]any{
				"value":   "v1.28",
				"allowed": []any{"current", "v1.25"},
				"prefix":  "runner",
				"disable": true,
			},
		},
	}

	for tn, tc := range testCases {
		t.Run(tn, func(t *testing.T) {
			test_tools.PlanAndAssertError(t, tc.vars, tc.wantErr)
		})
	}
}
