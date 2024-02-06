package scale_parameters

import (
	"testing"

	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/test_tools"
)

func TestScaleParameters(t *testing.T) {
	testCases := map[string]struct {
		vars    map[string]any
		wantErr bool
	}{
		"shell doesn't need scale parameters": {
			vars: map[string]any{
				"executor":              "shell",
				"scale_min":             -1,
				"scale_max":             -1,
				"idle_percentage":       -1,
				"capacity_per_instance": -1,
			},
		},
		"instance with scale parameters passes": {
			vars: map[string]any{
				"executor":              "instance",
				"scale_min":             0,
				"scale_max":             10,
				"idle_percentage":       10,
				"capacity_per_instance": 1,
			},
		},
		"instance without scale parameters fails": {
			vars: map[string]any{
				"executor":              "instance",
				"scale_min":             -1,
				"scale_max":             -1,
				"idle_percentage":       -1,
				"capacity_per_instance": -1,
			},
			wantErr: true,
		},
		"docker-autoscaler with scale parameters passes": {
			vars: map[string]any{
				"executor":              "docker-autoscaler",
				"scale_min":             0,
				"scale_max":             10,
				"idle_percentage":       10,
				"capacity_per_instance": 1,
			},
		},
		"docker-autoscaler without scale parameters fails": {
			vars: map[string]any{
				"executor":              "docker-autoscaler",
				"scale_min":             -1,
				"scale_max":             -1,
				"idle_percentage":       -1,
				"capacity_per_instance": -1,
			},
			wantErr: true,
		},
	}

	for tn, tc := range testCases {
		t.Run(tn, func(t *testing.T) {
			test_tools.PlanAndAssertError(t, tc.vars, tc.wantErr)
		})
	}
}
