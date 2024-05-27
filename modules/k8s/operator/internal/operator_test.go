package operator

import (
	"testing"

	"gitlab.com/gitlab-org/ci-cd/runner-tools/grit/test_tools"
)

func TestK8sOperator(t *testing.T) {
	plan := test_tools.Plan(t, nil)

	test_tools.AssertProviderConfigExists(t, plan, "kubectl")

	test_tools.AssertWithPlan(t, plan, []string{
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
	})
}
