package operator

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/test_tools"
)

func TestK8sOperator(t *testing.T) {
	testCases := map[string]struct {
		operatorVersion   string
		overrideFile      string
		expectedResources []string
	}{
		"defaults": {
			operatorVersion:   "current",
			overrideFile:      "",
			expectedResources: defaultExpectedResources,
		},
		"override-manifests": {
			operatorVersion: "will-be-ignored-because-we-explicitly-use-an-override",
			overrideFile: asFile(t, join(
				"apiVersion: v1",
				"kind: Namespace",
				"metadata:",
				"  name: some-namespace",
				"---",
				"apiVersion: apiextensions.k8s.io/v1",
				"kind: CustomResourceDefinition",
				"metadata:",
				"  name: some-crd",
				"spec: {}",
				"---",
				"apiVersion: v1",
				"kind: ConfigMap",
				"metadata:",
				"  name: some-cm",
				"  namespace: some-namespace",
				"data: {}",
			)),
			expectedResources: []string{
				`kubectl_manifest.operator_resources["v1::Namespace::_cluster_scoped_::some-namespace"]`,
				`kubectl_manifest.operator_resources["v1::ConfigMap::some-namespace::some-cm"]`,
				`kubectl_manifest.operator_resources["apiextensions.k8s.io/v1::CustomResourceDefinition::_cluster_scoped_::some-crd"]`,
			},
		},
	}

	for tn, tc := range testCases {
		t.Run(tn, func(t *testing.T) {
			plan := test_tools.Plan(t, map[string]any{
				"operator_version":   tc.operatorVersion,
				"override_manifests": tc.overrideFile,
			})
			test_tools.AssertProviderConfigExists(t, plan, "kubectl")
			test_tools.AssertWithPlan(t, plan, tc.expectedResources)
		})
	}
}

func TestK8sOperator_versions(t *testing.T) {
	testCases := map[string]bool{
		"no-such-version": false,
	}

	entries, err := os.ReadDir("versions")
	assert.NoError(t, err)
	for _, e := range entries {
		if e.Type() == fs.ModeSymlink || e.IsDir() {
			testCases[e.Name()] = true
		}
	}

	for version, expectSuccess := range testCases {
		t.Run(version, func(t *testing.T) {
			test_tools.PlanAndAssertError(t, map[string]any{
				"operator_version":   version,
				"override_manifests": "",
			}, !expectSuccess)
		})
	}
}

func asFile(t *testing.T, content string) string {
	t.Helper()

	filePath := filepath.Join(t.TempDir(), t.Name()+".yaml")
	assert.NoError(t, os.WriteFile(filePath, []byte(content), 0640))

	return filePath
}

func join(parts ...string) string {
	return strings.Join(parts, "\n")
}

var defaultExpectedResources = []string{
	`kubectl_manifest.operator_resources["apiextensions.k8s.io/v1::CustomResourceDefinition::_cluster_scoped_::runners.apps.gitlab.com"]`,
	`kubectl_manifest.operator_resources["apps/v1::Deployment::gitlab-runner-system::gitlab-runner-gitlab-runnercontroller-manager"]`,
	`kubectl_manifest.operator_resources["rbac.authorization.k8s.io/v1::ClusterRole::_cluster_scoped_::gitlab-runner-proxy-role"]`,
	`kubectl_manifest.operator_resources["rbac.authorization.k8s.io/v1::Role::gitlab-runner-system::gitlab-runner-app-role"]`,
	`kubectl_manifest.operator_resources["rbac.authorization.k8s.io/v1::ClusterRole::_cluster_scoped_::gitlab-runner-manager-role"]`,
	`kubectl_manifest.operator_resources["rbac.authorization.k8s.io/v1::ClusterRole::_cluster_scoped_::gitlab-runner-metrics-reader"]`,
	`kubectl_manifest.operator_resources["rbac.authorization.k8s.io/v1::ClusterRoleBinding::_cluster_scoped_::gitlab-runner-manager-rolebinding"]`,
	`kubectl_manifest.operator_resources["v1::Service::gitlab-runner-system::gitlab-runner-controller-manager-metrics-service"]`,
	`kubectl_manifest.operator_resources["rbac.authorization.k8s.io/v1::Role::gitlab-runner-system::gitlab-runner-leader-election-role"]`,
	`kubectl_manifest.operator_resources["v1::ServiceAccount::gitlab-runner-system::gitlab-runner-sa"]`,
	`kubectl_manifest.operator_resources["rbac.authorization.k8s.io/v1::ClusterRoleBinding::_cluster_scoped_::gitlab-runner-proxy-rolebinding"]`,
	`kubectl_manifest.operator_resources["rbac.authorization.k8s.io/v1::RoleBinding::gitlab-runner-system::gitlab-runner-app-rolebinding"]`,
	`kubectl_manifest.operator_resources["rbac.authorization.k8s.io/v1::RoleBinding::gitlab-runner-system::gitlab-runner-leader-election-rolebinding"]`,
	`kubectl_manifest.operator_resources["v1::Namespace::_cluster_scoped_::gitlab-runner-system"]`,
}
