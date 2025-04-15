package name

import (
	"strconv"
	"testing"

	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/test_tools"
)

func TestName(t *testing.T) {
	testCases := map[string]struct {
		vars    map[string]any
		wantErr bool
	}{
		"short label value is fine": {
			vars: map[string]any{
				"labels": map[string]string{
					"managed": "gritty-mcgritface",
				},
			},
		},
		"label with numbers, underscores, and dashes is fine": {
			vars: map[string]any{
				"labels": map[string]string{
					"managed": "gritty-mcgritface_the_gritter",
				},
			},
		},
		"long value gives an error": {
			vars: map[string]any{
				"labels": map[string]string{
					"managed": "gritty-mcgritface-the-great-and-powerful-gritmeister-extraordinaire", // 67 characters
				},
			},
			wantErr: true,
		},
		"exactly the limit is okay": {
			vars: map[string]any{
				"labels": map[string]string{
					"managed": "gritty-mcgritface-the-great-and-powerful-gritmeister-1234567890", // 63 characters
				},
			},
		},
		"error if first character of key is number": {
			vars: map[string]any{
				"labels": map[string]string{
					"1managed": "gritty",
				},
			},
			wantErr: true,
		},
		"error if first character of key is -": {
			vars: map[string]any{
				"labels": map[string]string{
					"-managed": "gritty",
				},
			},
			wantErr: true,
		},
		"error if first character of key is _": {
			vars: map[string]any{
				"labels": map[string]string{
					"_managed": "gritty",
				},
			},
			wantErr: true,
		},
		"error if more than 64 labels": {
			vars: map[string]any{
				"labels": lotsOfKeys(65),
			},
			wantErr: true,
		},
		"ok if exactly 64 labels": {
			vars: map[string]any{
				"labels": lotsOfKeys(64),
			},
		},
	}

	for tn, tc := range testCases {
		t.Run(tn, func(t *testing.T) {
			test_tools.PlanAndAssertError(t, tc.vars, tc.wantErr)
		})
	}
}

func lotsOfKeys(size int) map[string]string {
	m := make(map[string]string, 0)
	for i := range size {
		key := "key" + strconv.Itoa(i)
		m[key] = ""
	}
	return m
}
