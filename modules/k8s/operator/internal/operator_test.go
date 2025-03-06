package operator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/test_tools"
)

const defaultKubectlManifestTemplate = `kubectl_manifest.operator_resources["%s"]`

func TestK8sOperator(t *testing.T) {
	testCases := map[string]struct {
		operatorVersion   string
		overrideManifests string
		expectedResources []string
	}{
		"defaults": {
			operatorVersion:   "v1.31.0",
			overrideManifests: "",
			expectedResources: mapResources(defaultExpectedResources),
		},
		"override-manifests": {
			operatorVersion: "will-be-ignored-because-we-explicitly-use-an-override",
			overrideManifests: asFile(t, join(
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
			expectedResources: mapResources([]string{
				"Namespace:some-namespace",
				"CustomResourceDefinition:some-crd",
				"ConfigMap:some-cm",
			}),
		},
	}

	for tn, tc := range testCases {
		t.Run(tn, func(t *testing.T) {
			plan := test_tools.Plan(t, map[string]any{
				"operator_version":   tc.operatorVersion,
				"override_manifests": tc.overrideManifests,
			})
			test_tools.AssertProviderConfigExists(t, plan, "kubectl")
			test_tools.AssertWithPlan(t, plan, tc.expectedResources)
		})
	}
}

func asFile(t *testing.T, content string) string {
	t.Helper()

	filePath := filepath.Join(t.TempDir(), t.Name()+".yaml")
	assert.NoError(t, os.WriteFile(filePath, []byte(content), 0640))

	return "file://" + filePath
}

func join(parts ...string) string {
	return strings.Join(parts, "\n")
}

func mapResources(resources []string) []string {
	mapped := make([]string, 0, len(resources))
	for _, r := range resources {
		mapped = append(mapped, fmt.Sprintf(defaultKubectlManifestTemplate, r))
	}

	return mapped
}

var defaultExpectedResources = []string{
	"ClusterRole:gitlab-runner-manager-role",
	"ClusterRole:gitlab-runner-metrics-reader",
	"ClusterRole:gitlab-runner-proxy-role",
	"ClusterRoleBinding:gitlab-runner-manager-rolebinding",
	"ClusterRoleBinding:gitlab-runner-proxy-rolebinding",
	"CustomResourceDefinition:runners.apps.gitlab.com",
	"Deployment:gitlab-runner-gitlab-runnercontroller-manager",
	"Namespace:gitlab-runner-system",
	"Role:gitlab-runner-app-role",
	"Role:gitlab-runner-leader-election-role",
	"RoleBinding:gitlab-runner-app-rolebinding",
	"RoleBinding:gitlab-runner-leader-election-rolebinding",
	"Service:gitlab-runner-controller-manager-metrics-service",
	"ServiceAccount:gitlab-runner-sa",
}
