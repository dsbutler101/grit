package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/test_tools"
)

func TestNodeExporterContainer(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		moduleVars map[string]any
		wantOutput map[string]any
		wantErr    bool
	}{
		"empty inputs should error with required values": {
			wantErr: true,
		},
		"with tag": {
			moduleVars: map[string]any{
				"image_tag": "v1.2.3",
			},
			wantOutput: map[string]any{
				"container_config": map[string]any{
					"name":    "node-exporter",
					"image":   "quay.io/prometheus/node-exporter:v1.2.3",
					"network": "host",
					"pid":     "host",
					"volumes": []any{"/:/host:ro,rslave"},
					"command": "    --web.listen-address=0.0.0.0:9100 \\\n    --path.rootfs=/host\n",
					"service_options": []any{
						map[string]any{
							"ExecStartPost":  "/sbin/iptables -A INPUT -p tcp -m tcp --dport 9100 -j ACCEPT",
							"TimeoutStopSec": "30",
						},
					},
				},
			},
		},
		"with default overrides": {
			moduleVars: map[string]any{
				"image_tag":    "latest",
				"service_name": "custom-name",
				"port":         1234,
				"registry":     "custom-registry.com",
				"image_path":   "custom/image-path",
			},
			wantOutput: map[string]any{
				"container_config": map[string]any{
					"name":    "custom-name",
					"image":   "custom-registry.com/custom/image-path:latest",
					"network": "host",
					"pid":     "host",
					"volumes": []any{"/:/host:ro,rslave"},
					"command": "    --web.listen-address=0.0.0.0:1234 \\\n    --path.rootfs=/host\n",
					"service_options": []any{
						map[string]any{
							"ExecStartPost":  "/sbin/iptables -A INPUT -p tcp -m tcp --dport 1234 -j ACCEPT",
							"TimeoutStopSec": "30",
						},
					},
				},
			},
		},
	}

	for tn, tc := range testCases {
		t.Run(tn, func(t *testing.T) {
			op, err := test_tools.ApplyE(t, tc.moduleVars)
			if tc.wantErr {
				assert.Error(t, err)
				assert.Nil(t, op)
				return
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tc.wantOutput, op)
		})
	}
}
