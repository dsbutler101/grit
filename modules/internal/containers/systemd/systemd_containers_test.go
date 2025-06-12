package test

import (
	"bytes"
	"strings"
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/test_tools"
)

const serviceTemplate = `[Unit]
Description=Start {{.Name}} container
After=docker.service docker.socket
Wants=docker.service docker.socket

[Service]
ExecStart=/usr/bin/docker run \
            --rm \
            --name={{.Name}} \
            {{- range .Ports }}
            -p {{.}} \
            {{- end }}
            {{- range .Volumes }}
            -v {{.}} \
            {{- end }}
            {{- range .Env}}
            -e {{.}} \
            {{- end }}
            {{- if .Entrypoint }}
            --entrypoint {{.Entrypoint}} \
            {{- end }}
            {{- if .Network }}
            --network {{.Network}} \
            {{- end }}
            {{- if .Pid }}
            --pid {{.Pid}} \
            {{- end }}
            {{.Image}}{{- if .Command }} \
            {{.Command}}
            {{- end }}

ExecStop=/usr/bin/docker stop {{.Name}}
ExecStopPost=/usr/bin/docker rm {{.Name}}
{{- range .ServiceOptions }}
{{- range $k, $v := . }}
{{$k}}={{$v}}
{{- end }}
{{- end }}
`

type moduleVars struct {
	Containers  []Container `json:"containers"`
	Owner       string      `json:"owner,omitempty"`
	Permissions string      `json:"permissions,omitempty"`
	ServicePath string      `json:"service_path,omitempty"`
}

type Container struct {
	Name           string              `json:"name"`
	Ports          []string            `json:"ports"`
	Volumes        []string            `json:"volumes"`
	Env            []string            `json:"env"`
	Entrypoint     string              `json:"entrypoint,omitempty"`
	Network        string              `json:"network,omitempty"`
	Pid            string              `json:"pid,omitempty"`
	Image          string              `json:"image"`
	Command        string              `json:"command,omitempty"`
	ServiceOptions []map[string]string `json:"service_options"`
}

func expectedModuleOutput(t *testing.T, mv moduleVars) map[string]any {
	tplt := template.Must(template.New("service").Parse(serviceTemplate))

	cs := mv.Containers
	write_files := make([]any, len(cs))
	file_names := make([]string, len(cs))
	for i, c := range cs {
		// indent multi-line commands to replicate tf indent function
		c.Command = strings.ReplaceAll(c.Command, `\\\n`, "\\\n            ")

		var tpl bytes.Buffer
		err := tplt.Execute(&tpl, c)
		require.NoError(t, err)
		fn := c.Name + ".service"
		file_names[i] = fn

		path := "/etc/systemd/system"
		if mv.ServicePath != "" {
			path = mv.ServicePath
		}
		fn = path + "/" + fn

		owner := "root:root"
		if mv.Owner != "" {
			owner = mv.Owner
		}

		permissions := "0644"
		if mv.Permissions != "" {
			permissions = mv.Permissions
		}

		write_files[i] = map[string]any{
			"path":        fn,
			"owner":       owner,
			"permissions": permissions,
			"content":     tpl.String(),
		}
	}

	runcmd := ""
	if len(cs) > 0 {
		runcmd = "systemctl daemon-reload && systemctl enable --now " + strings.Join(file_names, " ")
	}

	return map[string]any{
		"write_files": write_files,
		"run_command": runcmd,
	}
}

func TestSystemdContainers(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		moduleVars moduleVars
		wantErr    bool
	}{
		"empty containers list is valid": {
			moduleVars{
				Containers: []Container{},
			},
			false,
		},
		"duplicate container names are not allowed": {
			moduleVars{
				Containers: []Container{
					{Name: "foo", Image: "foo:latest"},
					{Name: "bar", Image: "bar:latest"},
					{Name: "foo", Image: "baz:latest"},
				},
			},
			true,
		},
		"container names must not be longer than 63 characters": {
			moduleVars{
				Containers: []Container{
					{Name: "this-is-a-very-long-container-name-that-exceeds-the-max-allowed-length", Image: "foo:latest"},
				}},
			true,
		},
		"container names must not contain invalid characters": {
			moduleVars{
				Containers: []Container{
					{Name: "container#name", Image: "foo:latest"},
				}},
			true,
		},
		"container names must start with alphanumeric": {
			moduleVars{
				Containers: []Container{
					{Name: "_must_start_with_alphanumeric", Image: "bar:latest"},
				}},
			true,
		},
		"some valid container names": {
			moduleVars{
				Containers: []Container{
					{Name: "can_have_underscores", Image: "foo:latest"},
					{Name: "can.have.dots", Image: "foo:latest"},
					{Name: "can-have-dashes", Image: "foo:latest"},
					{Name: "justalpha", Image: "foo:latest"},
					{Name: "mix01234num", Image: "foo:latest"},
					{Name: "1234numstart", Image: "foo:latest"},
				}},
			false,
		},
		"with some ports": {
			moduleVars{
				Containers: []Container{
					{Name: "web", Image: "nginx:latest", Ports: []string{"80:80/tcp", "443:443/tcp"}},
				}},
			false,
		},
		"with some volumes": {
			moduleVars{
				Containers: []Container{
					{Name: "db", Image: "postgres:latest", Volumes: []string{"/data:/var/lib/postgresql/data"}},
				}},
			false,
		},
		"with some environment variables": {
			moduleVars{
				Containers: []Container{
					{Name: "web", Image: "nginx:latest", Env: []string{"FOO=bar", "BAZ=qux"}},
				}},
			false,
		},
		"with entrypoint": {
			moduleVars{
				Containers: []Container{
					{Name: "web", Image: "nginx:latest", Entrypoint: "/usr/sbin/nginx"},
				}},
			false,
		},
		"with custom network": {
			moduleVars{
				Containers: []Container{
					{Name: "web", Image: "nginx:latest", Network: "foo"},
				}},
			false,
		},
		"with host pid": {
			moduleVars{
				Containers: []Container{
					{Name: "web", Image: "nginx:latest", Pid: "host"},
				}},
			false,
		},
		"with command": {
			moduleVars{
				Containers: []Container{
					{Name: "web", Image: "nginx:latest", Command: "sleep infinity"},
				}},
			false,
		},
		"with multi-line command": {
			moduleVars{
				Containers: []Container{
					{
						Name: "gitlab-runner", Image: "gitlab-runner:latest",
						Command: `run \\\n--config /etc/gitlab-runner/config.toml \\\n--user=gitlab-runner \\\n--working-directory=/home/gitlab-runner`,
					},
				},
			},
			false,
		},
		"with service options": {
			moduleVars{
				Containers: []Container{
					{Name: "web", Image: "nginx:latest", ServiceOptions: []map[string]string{
						{
							"Restart":         "always",
							"TimeoutStartSec": "0",
						},
					}},
				}},
			false,
		},
		"with multiple ordered service options": {
			moduleVars{
				Containers: []Container{
					{Name: "gitlab-runner", Image: "gitlab-runner:latest", ServiceOptions: []map[string]string{
						{
							"ExecStart": "",
						},
						{
							"ExecStart": "/opt/foo/bar",
						},
						{
							"Restart":         "always",
							"TimeoutStartSec": "0",
						},
					}},
				}},
			false,
		},
		"with custom owners, permissions, and path": {
			moduleVars{
				Containers:  []Container{{Name: "web", Image: "nginx:latest"}},
				Owner:       "root:docker",
				Permissions: "0755",
				ServicePath: "/etc/systemd/system/custom",
			},
			false,
		},
	}

	for tn, tc := range testCases {
		t.Run(tn, func(t *testing.T) {
			vars := test_tools.ToModuleVars(t, tc.moduleVars)

			op, err := test_tools.ApplyE(t, vars)
			if tc.wantErr {
				assert.Error(t, err)
				assert.Nil(t, op)
				return
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, expectedModuleOutput(t, tc.moduleVars), op)
		})
	}
}
