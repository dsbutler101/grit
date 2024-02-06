package support

import (
	"testing"

	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/test_tools"
)

func TestSupport(t *testing.T) {
	testCases := map[string]struct {
		vars    map[string]any
		wantErr bool
	}{
		"invalid min_support value": {
			vars: map[string]any{
				"min_support": "super_duper",
			},
			wantErr: true,
		},
		"experimental support okay for experimental requirements": {
			vars: map[string]any{
				"min_support": "experimental",
				"use_case":    "instance",
				"use_case_support": map[string]any{
					"instance": "experimental",
				},
			},
		},
		"experimental support not okay for beta requirements": {
			vars: map[string]any{
				"min_support": "beta",
				"use_case":    "instance",
				"use_case_support": map[string]any{
					"instance": "experimental",
				},
			},
			wantErr: true,
		},
		"experimental support not okay for ga requirements": {
			vars: map[string]any{
				"min_support": "ga",
				"use_case":    "instance",
				"use_case_support": map[string]any{
					"instance": "experimental",
				},
			},
			wantErr: true,
		},
		"beta support okay for experimental requirements": {
			vars: map[string]any{
				"min_support": "experimental",
				"use_case":    "instance",
				"use_case_support": map[string]any{
					"instance": "beta",
				},
			},
		},
		"beta support okay for beta requirements": {
			vars: map[string]any{
				"min_support": "beta",
				"use_case":    "instance",
				"use_case_support": map[string]any{
					"instance": "beta",
				},
			},
		},
		"beta support not okay for ga requirements": {
			vars: map[string]any{
				"min_support": "ga",
				"use_case":    "instance",
				"use_case_support": map[string]any{
					"instance": "beta",
				},
			},
			wantErr: true,
		},
		"ga support okay for experimental requirements": {
			vars: map[string]any{
				"min_support": "experimental",
				"use_case":    "instance",
				"use_case_support": map[string]any{
					"instance": "ga",
				},
			},
		},
		"ga support okay for beta requirements": {
			vars: map[string]any{
				"min_support": "beta",
				"use_case":    "instance",
				"use_case_support": map[string]any{
					"instance": "ga",
				},
			},
		},
		"ga support okay for ga requirements": {
			vars: map[string]any{
				"min_support": "ga",
				"use_case":    "instance",
				"use_case_support": map[string]any{
					"instance": "ga",
				},
			},
		},
	}

	for tn, tc := range testCases {
		t.Run(tn, func(t *testing.T) {
			test_tools.PlanAndAssertError(t, tc.vars, tc.wantErr)
		})
	}
}
