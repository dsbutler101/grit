---
stage: Verify
group: Runner
info: To determine the technical writer assigned to the Stage/Group associated with this page, see https://handbook.gitlab.com/handbook/product/ux/technical-writing/#assignments
---

# Scenario: Google - Kubernetes/GKE - GitLab Runner Operator

This scenario template deploys GitLab Runner to a GKE cluster by leveraging the
GitLab Runner Operator.

## Prerequisites

To use this scenario, you must meet the following prerequisites:

### Terraform and Google Cloud setup

To use this scenario, you must have:

- [Terraform prerequisites](../../../index.md#prerequisites)
- [Google Cloud prerequisites](../../../index.md#google-cloud-integration-prerequisites)
- Terraform 1.7 or later to use Terraform features and syntax specific to this scenario

### API services

The following API services must be enabled in the Google Cloud project:

- `compute.googleapis.com`
- `container.googleapis.com`

Ask someone with `owner` access to your Google Cloud project to run the following
command:

```shell
gcloud services enable compute.googleapis.com container.googleapis.com
```

### Google Cloud permissions for Terraform execution

The actors that execute the Terraform code must have the
following [permissions](https://cloud.google.com/kms/docs/reference/permissions-and-roles) in Google Cloud:

<details>
<summary>Required permissions</summary>

- `compute.disks.create`
- `compute.firewalls.create`
- `compute.firewalls.delete`
- `compute.firewalls.get`
- `compute.instanceGroupManagers.create`
- `compute.instanceGroupManagers.delete`
- `compute.instanceGroupManagers.get`
- `compute.instanceGroups.create`
- `compute.instanceGroups.delete`
- `compute.instanceTemplates.create`
- `compute.instanceTemplates.delete`
- `compute.instanceTemplates.get`
- `compute.instanceTemplates.useReadOnly`
- `compute.instances.create`
- `compute.instances.delete`
- `compute.instances.get`
- `compute.instances.setLabels`
- `compute.instances.setMetadata`
- `compute.instances.setServiceAccount`
- `compute.instances.setTags`
- `compute.networks.create`
- `compute.networks.delete`
- `compute.networks.get`
- `compute.networks.updatePolicy`
- `compute.regionOperations.get`
- `compute.subnetworks.create`
- `compute.subnetworks.delete`
- `compute.subnetworks.get`
- `compute.subnetworks.useExternalIp`
- `compute.subnetworks.use`
- `compute.zones.get`
- `container.apiServices.create`
- `container.apiServices.delete`
- `container.apiServices.getStatus`
- `container.apiServices.get`
- `container.apiServices.list`
- `container.apiServices.updateStatus`
- `container.apiServices.update`
- `container.auditSinks.create`
- `container.auditSinks.delete`
- `container.auditSinks.get`
- `container.auditSinks.list`
- `container.auditSinks.update`
- `container.backendConfigs.create`
- `container.backendConfigs.delete`
- `container.backendConfigs.get`
- `container.backendConfigs.list`
- `container.backendConfigs.update`
- `container.bindings.create`
- `container.bindings.delete`
- `container.bindings.get`
- `container.bindings.list`
- `container.bindings.update`
- `container.certificateSigningRequests.approve`
- `container.certificateSigningRequests.create`
- `container.certificateSigningRequests.delete`
- `container.certificateSigningRequests.getStatus`
- `container.certificateSigningRequests.get`
- `container.certificateSigningRequests.list`
- `container.certificateSigningRequests.updateStatus`
- `container.certificateSigningRequests.update`
- `container.clusterRoleBindings.create`
- `container.clusterRoleBindings.delete`
- `container.clusterRoleBindings.get`
- `container.clusterRoleBindings.list`
- `container.clusterRoleBindings.update`
- `container.clusterRoles.bind`
- `container.clusterRoles.create`
- `container.clusterRoles.delete`
- `container.clusterRoles.escalate`
- `container.clusterRoles.get`
- `container.clusterRoles.list`
- `container.clusterRoles.update`
- `container.clusters.connect`
- `container.clusters.createTagBinding`
- `container.clusters.create`
- `container.clusters.deleteTagBinding`
- `container.clusters.delete`
- `container.clusters.getCredentials`
- `container.clusters.get`
- `container.clusters.impersonate`
- `container.clusters.listEffectiveTags`
- `container.clusters.listTagBindings`
- `container.clusters.list`
- `container.clusters.update`
- `container.componentStatuses.get`
- `container.componentStatuses.list`
- `container.configMaps.create`
- `container.configMaps.delete`
- `container.configMaps.get`
- `container.configMaps.list`
- `container.configMaps.update`
- `container.controllerRevisions.create`
- `container.controllerRevisions.delete`
- `container.controllerRevisions.get`
- `container.controllerRevisions.list`
- `container.controllerRevisions.update`
- `container.cronJobs.create`
- `container.cronJobs.delete`
- `container.cronJobs.getStatus`
- `container.cronJobs.get`
- `container.cronJobs.list`
- `container.cronJobs.updateStatus`
- `container.cronJobs.update`
- `container.csiDrivers.create`
- `container.csiDrivers.delete`
- `container.csiDrivers.get`
- `container.csiDrivers.list`
- `container.csiDrivers.update`
- `container.csiNodeInfos.create`
- `container.csiNodeInfos.delete`
- `container.csiNodeInfos.get`
- `container.csiNodeInfos.list`
- `container.csiNodeInfos.update`
- `container.csiNodes.create`
- `container.csiNodes.delete`
- `container.csiNodes.get`
- `container.csiNodes.list`
- `container.csiNodes.update`
- `container.customResourceDefinitions.create`
- `container.customResourceDefinitions.delete`
- `container.customResourceDefinitions.getStatus`
- `container.customResourceDefinitions.get`
- `container.customResourceDefinitions.list`
- `container.customResourceDefinitions.updateStatus`
- `container.customResourceDefinitions.update`
- `container.daemonSets.create`
- `container.daemonSets.delete`
- `container.daemonSets.getStatus`
- `container.daemonSets.get`
- `container.daemonSets.list`
- `container.daemonSets.updateStatus`
- `container.daemonSets.update`
- `container.deployments.create`
- `container.deployments.delete`
- `container.deployments.getScale`
- `container.deployments.getStatus`
- `container.deployments.get`
- `container.deployments.list`
- `container.deployments.rollback`
- `container.deployments.updateScale`
- `container.deployments.updateStatus`
- `container.deployments.update`
- `container.endpointSlices.create`
- `container.endpointSlices.delete`
- `container.endpointSlices.get`
- `container.endpointSlices.list`
- `container.endpointSlices.update`
- `container.endpoints.create`
- `container.endpoints.delete`
- `container.endpoints.get`
- `container.endpoints.list`
- `container.endpoints.update`
- `container.events.create`
- `container.events.delete`
- `container.events.get`
- `container.events.list`
- `container.events.update`
- `container.frontendConfigs.create`
- `container.frontendConfigs.delete`
- `container.frontendConfigs.get`
- `container.frontendConfigs.list`
- `container.frontendConfigs.update`
- `container.horizontalPodAutoscalers.create`
- `container.horizontalPodAutoscalers.delete`
- `container.horizontalPodAutoscalers.getStatus`
- `container.horizontalPodAutoscalers.get`
- `container.horizontalPodAutoscalers.list`
- `container.horizontalPodAutoscalers.updateStatus`
- `container.horizontalPodAutoscalers.update`
- `container.hostServiceAgent.use`
- `container.ingresses.create`
- `container.ingresses.delete`
- `container.ingresses.getStatus`
- `container.ingresses.get`
- `container.ingresses.list`
- `container.ingresses.updateStatus`
- `container.ingresses.update`
- `container.initializerConfigurations.create`
- `container.initializerConfigurations.delete`
- `container.initializerConfigurations.get`
- `container.initializerConfigurations.list`
- `container.initializerConfigurations.update`
- `container.jobs.create`
- `container.jobs.delete`
- `container.jobs.getStatus`
- `container.jobs.get`
- `container.jobs.list`
- `container.jobs.updateStatus`
- `container.jobs.update`
- `container.leases.create`
- `container.leases.delete`
- `container.leases.get`
- `container.leases.list`
- `container.leases.update`
- `container.limitRanges.create`
- `container.limitRanges.delete`
- `container.limitRanges.get`
- `container.limitRanges.list`
- `container.limitRanges.update`
- `container.localSubjectAccessReviews.create`
- `container.localSubjectAccessReviews.list`
- `container.managedCertificates.create`
- `container.managedCertificates.delete`
- `container.managedCertificates.get`
- `container.managedCertificates.list`
- `container.managedCertificates.update`
- `container.mutatingWebhookConfigurations.create`
- `container.mutatingWebhookConfigurations.delete`
- `container.mutatingWebhookConfigurations.get`
- `container.mutatingWebhookConfigurations.list`
- `container.mutatingWebhookConfigurations.update`
- `container.namespaces.create`
- `container.namespaces.delete`
- `container.namespaces.finalize`
- `container.namespaces.getStatus`
- `container.namespaces.get`
- `container.namespaces.list`
- `container.namespaces.updateStatus`
- `container.namespaces.update`
- `container.networkPolicies.create`
- `container.networkPolicies.delete`
- `container.networkPolicies.get`
- `container.networkPolicies.list`
- `container.networkPolicies.update`
- `container.nodes.create`
- `container.nodes.delete`
- `container.nodes.getStatus`
- `container.nodes.get`
- `container.nodes.list`
- `container.nodes.proxy`
- `container.nodes.updateStatus`
- `container.nodes.update`
- `container.operations.get`
- `container.operations.list`
- `container.persistentVolumeClaims.create`
- `container.persistentVolumeClaims.delete`
- `container.persistentVolumeClaims.getStatus`
- `container.persistentVolumeClaims.get`
- `container.persistentVolumeClaims.list`
- `container.persistentVolumeClaims.updateStatus`
- `container.persistentVolumeClaims.update`
- `container.persistentVolumes.create`
- `container.persistentVolumes.delete`
- `container.persistentVolumes.getStatus`
- `container.persistentVolumes.get`
- `container.persistentVolumes.list`
- `container.persistentVolumes.updateStatus`
- `container.persistentVolumes.update`
- `container.petSets.create`
- `container.petSets.delete`
- `container.petSets.get`
- `container.petSets.list`
- `container.petSets.updateStatus`
- `container.petSets.update`
- `container.podDisruptionBudgets.create`
- `container.podDisruptionBudgets.delete`
- `container.podDisruptionBudgets.getStatus`
- `container.podDisruptionBudgets.get`
- `container.podDisruptionBudgets.list`
- `container.podDisruptionBudgets.updateStatus`
- `container.podDisruptionBudgets.update`
- `container.podPresets.create`
- `container.podPresets.delete`
- `container.podPresets.get`
- `container.podPresets.list`
- `container.podPresets.update`
- `container.podSecurityPolicies.create`
- `container.podSecurityPolicies.delete`
- `container.podSecurityPolicies.get`
- `container.podSecurityPolicies.list`
- `container.podSecurityPolicies.update`
- `container.podSecurityPolicies.use`
- `container.podTemplates.create`
- `container.podTemplates.delete`
- `container.podTemplates.get`
- `container.podTemplates.list`
- `container.podTemplates.update`
- `container.pods.attach`
- `container.pods.create`
- `container.pods.delete`
- `container.pods.evict`
- `container.pods.exec`
- `container.pods.getLogs`
- `container.pods.getStatus`
- `container.pods.get`
- `container.pods.initialize`
- `container.pods.list`
- `container.pods.portForward`
- `container.pods.proxy`
- `container.pods.updateStatus`
- `container.pods.update`
- `container.priorityClasses.create`
- `container.priorityClasses.delete`
- `container.priorityClasses.get`
- `container.priorityClasses.list`
- `container.priorityClasses.update`
- `container.replicaSets.create`
- `container.replicaSets.delete`
- `container.replicaSets.getScale`
- `container.replicaSets.getStatus`
- `container.replicaSets.get`
- `container.replicaSets.list`
- `container.replicaSets.updateScale`
- `container.replicaSets.updateStatus`
- `container.replicaSets.update`
- `container.replicationControllers.create`
- `container.replicationControllers.delete`
- `container.replicationControllers.getScale`
- `container.replicationControllers.getStatus`
- `container.replicationControllers.get`
- `container.replicationControllers.list`
- `container.replicationControllers.updateScale`
- `container.replicationControllers.updateStatus`
- `container.replicationControllers.update`
- `container.resourceQuotas.create`
- `container.resourceQuotas.delete`
- `container.resourceQuotas.getStatus`
- `container.resourceQuotas.get`
- `container.resourceQuotas.list`
- `container.resourceQuotas.updateStatus`
- `container.resourceQuotas.update`
- `container.roleBindings.create`
- `container.roleBindings.delete`
- `container.roleBindings.get`
- `container.roleBindings.list`
- `container.roleBindings.update`
- `container.roles.bind`
- `container.roles.create`
- `container.roles.delete`
- `container.roles.escalate`
- `container.roles.get`
- `container.roles.list`
- `container.roles.update`
- `container.runtimeClasses.create`
- `container.runtimeClasses.delete`
- `container.runtimeClasses.get`
- `container.runtimeClasses.list`
- `container.runtimeClasses.update`
- `container.scheduledJobs.create`
- `container.scheduledJobs.delete`
- `container.scheduledJobs.get`
- `container.scheduledJobs.list`
- `container.scheduledJobs.updateStatus`
- `container.scheduledJobs.update`
- `container.secrets.create`
- `container.secrets.delete`
- `container.secrets.get`
- `container.secrets.list`
- `container.secrets.update`
- `container.selfSubjectAccessReviews.create`
- `container.selfSubjectAccessReviews.list`
- `container.selfSubjectRulesReviews.create`
- `container.serviceAccounts.createToken`
- `container.serviceAccounts.create`
- `container.serviceAccounts.delete`
- `container.serviceAccounts.get`
- `container.serviceAccounts.list`
- `container.serviceAccounts.update`
- `container.services.create`
- `container.services.delete`
- `container.services.getStatus`
- `container.services.get`
- `container.services.list`
- `container.services.proxy`
- `container.services.updateStatus`
- `container.services.update`
- `container.statefulSets.create`
- `container.statefulSets.delete`
- `container.statefulSets.getScale`
- `container.statefulSets.getStatus`
- `container.statefulSets.get`
- `container.statefulSets.list`
- `container.statefulSets.updateScale`
- `container.statefulSets.updateStatus`
- `container.statefulSets.update`
- `container.storageClasses.create`
- `container.storageClasses.delete`
- `container.storageClasses.get`
- `container.storageClasses.list`
- `container.storageClasses.update`
- `container.storageStates.create`
- `container.storageStates.delete`
- `container.storageStates.getStatus`
- `container.storageStates.get`
- `container.storageStates.list`
- `container.storageStates.updateStatus`
- `container.storageStates.update`
- `container.storageVersionMigrations.create`
- `container.storageVersionMigrations.delete`
- `container.storageVersionMigrations.getStatus`
- `container.storageVersionMigrations.get`
- `container.storageVersionMigrations.list`
- `container.storageVersionMigrations.updateStatus`
- `container.storageVersionMigrations.update`
- `container.subjectAccessReviews.create`
- `container.subjectAccessReviews.list`
- `container.thirdPartyObjects.create`
- `container.thirdPartyObjects.delete`
- `container.thirdPartyObjects.get`
- `container.thirdPartyObjects.list`
- `container.thirdPartyObjects.update`
- `container.thirdPartyResources.create`
- `container.thirdPartyResources.delete`
- `container.thirdPartyResources.get`
- `container.thirdPartyResources.list`
- `container.thirdPartyResources.update`
- `container.tokenReviews.create`
- `container.updateInfos.create`
- `container.updateInfos.delete`
- `container.updateInfos.get`
- `container.updateInfos.list`
- `container.updateInfos.update`
- `container.validatingWebhookConfigurations.create`
- `container.validatingWebhookConfigurations.delete`
- `container.validatingWebhookConfigurations.get`
- `container.validatingWebhookConfigurations.list`
- `container.validatingWebhookConfigurations.update`
- `container.volumeAttachments.create`
- `container.volumeAttachments.delete`
- `container.volumeAttachments.getStatus`
- `container.volumeAttachments.get`
- `container.volumeAttachments.list`
- `container.volumeAttachments.updateStatus`
- `container.volumeAttachments.update`
- `container.volumeSnapshotClasses.create`
- `container.volumeSnapshotClasses.delete`
- `container.volumeSnapshotClasses.get`
- `container.volumeSnapshotClasses.list`
- `container.volumeSnapshotClasses.update`
- `container.volumeSnapshotContents.create`
- `container.volumeSnapshotContents.delete`
- `container.volumeSnapshotContents.getStatus`
- `container.volumeSnapshotContents.get`
- `container.volumeSnapshotContents.list`
- `container.volumeSnapshotContents.updateStatus`
- `container.volumeSnapshotContents.update`
- `container.volumeSnapshots.create`
- `container.volumeSnapshots.delete`
- `container.volumeSnapshots.getStatus`
- `container.volumeSnapshots.get`
- `container.volumeSnapshots.list`
- `container.volumeSnapshots.updateStatus`
- `container.volumeSnapshots.update`
- `recommender.containerDiagnosisInsights.get`
- `recommender.containerDiagnosisInsights.list`
- `recommender.containerDiagnosisInsights.update`
- `recommender.containerDiagnosisRecommendations.get`
- `recommender.containerDiagnosisRecommendations.list`
- `recommender.containerDiagnosisRecommendations.update`
- `recommender.locations.get`
- `recommender.locations.list`
- `recommender.networkAnalyzerGkeConnectivityInsights.get`
- `recommender.networkAnalyzerGkeConnectivityInsights.list`
- `recommender.networkAnalyzerGkeConnectivityInsights.update`
- `recommender.networkAnalyzerGkeIpAddressInsights.get`
- `recommender.networkAnalyzerGkeIpAddressInsights.list`
- `recommender.networkAnalyzerGkeIpAddressInsights.update`
- `resourcemanager.projects.get`
- `resourcemanager.projects.list`

</details>

You can also create a [custom role](https://cloud.google.com/iam/docs/creating-custom-roles)
with these permissions. You can then assign this role to the user or service account
responsible for provisioning the GRIT Terraform configuration.

Ask someone with `owner` access to your Google Cloud project to run the following
commands:

<details>

```shell
cat > grit-provisioner-role.json <<EOF
{
  "title": "GRITProvisioner",
  "description": "A role with minimum list of permissions required for GRIT provisioning",
  "includedPermissions": [
    "compute.disks.create",
    "compute.firewalls.create",
    "compute.firewalls.delete",
    "compute.firewalls.get",
    "compute.instanceGroupManagers.create",
    "compute.instanceGroupManagers.delete",
    "compute.instanceGroupManagers.get",
    "compute.instanceGroups.create",
    "compute.instanceGroups.delete",
    "compute.instanceTemplates.create",
    "compute.instanceTemplates.delete",
    "compute.instanceTemplates.get",
    "compute.instanceTemplates.useReadOnly",
    "compute.instances.create",
    "compute.instances.delete",
    "compute.instances.get",
    "compute.instances.setLabels",
    "compute.instances.setMetadata",
    "compute.instances.setServiceAccount",
    "compute.instances.setTags",
    "compute.networks.create",
    "compute.networks.delete",
    "compute.networks.get",
    "compute.networks.updatePolicy",
    "compute.regionOperations.get",
    "compute.subnetworks.create",
    "compute.subnetworks.delete",
    "compute.subnetworks.get",
    "compute.subnetworks.useExternalIp",
    "compute.subnetworks.use",
    "compute.zones.get",
    "container.apiServices.create",
    "container.apiServices.delete",
    "container.apiServices.getStatus",
    "container.apiServices.get",
    "container.apiServices.list",
    "container.apiServices.updateStatus",
    "container.apiServices.update",
    "container.auditSinks.create",
    "container.auditSinks.delete",
    "container.auditSinks.get",
    "container.auditSinks.list",
    "container.auditSinks.update",
    "container.backendConfigs.create",
    "container.backendConfigs.delete",
    "container.backendConfigs.get",
    "container.backendConfigs.list",
    "container.backendConfigs.update",
    "container.bindings.create",
    "container.bindings.delete",
    "container.bindings.get",
    "container.bindings.list",
    "container.bindings.update",
    "container.certificateSigningRequests.approve",
    "container.certificateSigningRequests.create",
    "container.certificateSigningRequests.delete",
    "container.certificateSigningRequests.getStatus",
    "container.certificateSigningRequests.get",
    "container.certificateSigningRequests.list",
    "container.certificateSigningRequests.updateStatus",
    "container.certificateSigningRequests.update",
    "container.clusterRoleBindings.create",
    "container.clusterRoleBindings.delete",
    "container.clusterRoleBindings.get",
    "container.clusterRoleBindings.list",
    "container.clusterRoleBindings.update",
    "container.clusterRoles.bind",
    "container.clusterRoles.create",
    "container.clusterRoles.delete",
    "container.clusterRoles.escalate",
    "container.clusterRoles.get",
    "container.clusterRoles.list",
    "container.clusterRoles.update",
    "container.clusters.connect",
    "container.clusters.createTagBinding",
    "container.clusters.create",
    "container.clusters.deleteTagBinding",
    "container.clusters.delete",
    "container.clusters.getCredentials",
    "container.clusters.get",
    "container.clusters.impersonate",
    "container.clusters.listEffectiveTags",
    "container.clusters.listTagBindings",
    "container.clusters.list",
    "container.clusters.update",
    "container.componentStatuses.get",
    "container.componentStatuses.list",
    "container.configMaps.create",
    "container.configMaps.delete",
    "container.configMaps.get",
    "container.configMaps.list",
    "container.configMaps.update",
    "container.controllerRevisions.create",
    "container.controllerRevisions.delete",
    "container.controllerRevisions.get",
    "container.controllerRevisions.list",
    "container.controllerRevisions.update",
    "container.cronJobs.create",
    "container.cronJobs.delete",
    "container.cronJobs.getStatus",
    "container.cronJobs.get",
    "container.cronJobs.list",
    "container.cronJobs.updateStatus",
    "container.cronJobs.update",
    "container.csiDrivers.create",
    "container.csiDrivers.delete",
    "container.csiDrivers.get",
    "container.csiDrivers.list",
    "container.csiDrivers.update",
    "container.csiNodeInfos.create",
    "container.csiNodeInfos.delete",
    "container.csiNodeInfos.get",
    "container.csiNodeInfos.list",
    "container.csiNodeInfos.update",
    "container.csiNodes.create",
    "container.csiNodes.delete",
    "container.csiNodes.get",
    "container.csiNodes.list",
    "container.csiNodes.update",
    "container.customResourceDefinitions.create",
    "container.customResourceDefinitions.delete",
    "container.customResourceDefinitions.getStatus",
    "container.customResourceDefinitions.get",
    "container.customResourceDefinitions.list",
    "container.customResourceDefinitions.updateStatus",
    "container.customResourceDefinitions.update",
    "container.daemonSets.create",
    "container.daemonSets.delete",
    "container.daemonSets.getStatus",
    "container.daemonSets.get",
    "container.daemonSets.list",
    "container.daemonSets.updateStatus",
    "container.daemonSets.update",
    "container.deployments.create",
    "container.deployments.delete",
    "container.deployments.getScale",
    "container.deployments.getStatus",
    "container.deployments.get",
    "container.deployments.list",
    "container.deployments.rollback",
    "container.deployments.updateScale",
    "container.deployments.updateStatus",
    "container.deployments.update",
    "container.endpointSlices.create",
    "container.endpointSlices.delete",
    "container.endpointSlices.get",
    "container.endpointSlices.list",
    "container.endpointSlices.update",
    "container.endpoints.create",
    "container.endpoints.delete",
    "container.endpoints.get",
    "container.endpoints.list",
    "container.endpoints.update",
    "container.events.create",
    "container.events.delete",
    "container.events.get",
    "container.events.list",
    "container.events.update",
    "container.frontendConfigs.create",
    "container.frontendConfigs.delete",
    "container.frontendConfigs.get",
    "container.frontendConfigs.list",
    "container.frontendConfigs.update",
    "container.horizontalPodAutoscalers.create",
    "container.horizontalPodAutoscalers.delete",
    "container.horizontalPodAutoscalers.getStatus",
    "container.horizontalPodAutoscalers.get",
    "container.horizontalPodAutoscalers.list",
    "container.horizontalPodAutoscalers.updateStatus",
    "container.horizontalPodAutoscalers.update",
    "container.hostServiceAgent.use",
    "container.ingresses.create",
    "container.ingresses.delete",
    "container.ingresses.getStatus",
    "container.ingresses.get",
    "container.ingresses.list",
    "container.ingresses.updateStatus",
    "container.ingresses.update",
    "container.initializerConfigurations.create",
    "container.initializerConfigurations.delete",
    "container.initializerConfigurations.get",
    "container.initializerConfigurations.list",
    "container.initializerConfigurations.update",
    "container.jobs.create",
    "container.jobs.delete",
    "container.jobs.getStatus",
    "container.jobs.get",
    "container.jobs.list",
    "container.jobs.updateStatus",
    "container.jobs.update",
    "container.leases.create",
    "container.leases.delete",
    "container.leases.get",
    "container.leases.list",
    "container.leases.update",
    "container.limitRanges.create",
    "container.limitRanges.delete",
    "container.limitRanges.get",
    "container.limitRanges.list",
    "container.limitRanges.update",
    "container.localSubjectAccessReviews.create",
    "container.localSubjectAccessReviews.list",
    "container.managedCertificates.create",
    "container.managedCertificates.delete",
    "container.managedCertificates.get",
    "container.managedCertificates.list",
    "container.managedCertificates.update",
    "container.mutatingWebhookConfigurations.create",
    "container.mutatingWebhookConfigurations.delete",
    "container.mutatingWebhookConfigurations.get",
    "container.mutatingWebhookConfigurations.list",
    "container.mutatingWebhookConfigurations.update",
    "container.namespaces.create",
    "container.namespaces.delete",
    "container.namespaces.finalize",
    "container.namespaces.getStatus",
    "container.namespaces.get",
    "container.namespaces.list",
    "container.namespaces.updateStatus",
    "container.namespaces.update",
    "container.networkPolicies.create",
    "container.networkPolicies.delete",
    "container.networkPolicies.get",
    "container.networkPolicies.list",
    "container.networkPolicies.update",
    "container.nodes.create",
    "container.nodes.delete",
    "container.nodes.getStatus",
    "container.nodes.get",
    "container.nodes.list",
    "container.nodes.proxy",
    "container.nodes.updateStatus",
    "container.nodes.update",
    "container.operations.get",
    "container.operations.list",
    "container.persistentVolumeClaims.create",
    "container.persistentVolumeClaims.delete",
    "container.persistentVolumeClaims.getStatus",
    "container.persistentVolumeClaims.get",
    "container.persistentVolumeClaims.list",
    "container.persistentVolumeClaims.updateStatus",
    "container.persistentVolumeClaims.update",
    "container.persistentVolumes.create",
    "container.persistentVolumes.delete",
    "container.persistentVolumes.getStatus",
    "container.persistentVolumes.get",
    "container.persistentVolumes.list",
    "container.persistentVolumes.updateStatus",
    "container.persistentVolumes.update",
    "container.petSets.create",
    "container.petSets.delete",
    "container.petSets.get",
    "container.petSets.list",
    "container.petSets.updateStatus",
    "container.petSets.update",
    "container.podDisruptionBudgets.create",
    "container.podDisruptionBudgets.delete",
    "container.podDisruptionBudgets.getStatus",
    "container.podDisruptionBudgets.get",
    "container.podDisruptionBudgets.list",
    "container.podDisruptionBudgets.updateStatus",
    "container.podDisruptionBudgets.update",
    "container.podPresets.create",
    "container.podPresets.delete",
    "container.podPresets.get",
    "container.podPresets.list",
    "container.podPresets.update",
    "container.podSecurityPolicies.create",
    "container.podSecurityPolicies.delete",
    "container.podSecurityPolicies.get",
    "container.podSecurityPolicies.list",
    "container.podSecurityPolicies.update",
    "container.podSecurityPolicies.use",
    "container.podTemplates.create",
    "container.podTemplates.delete",
    "container.podTemplates.get",
    "container.podTemplates.list",
    "container.podTemplates.update",
    "container.pods.attach",
    "container.pods.create",
    "container.pods.delete",
    "container.pods.evict",
    "container.pods.exec",
    "container.pods.getLogs",
    "container.pods.getStatus",
    "container.pods.get",
    "container.pods.initialize",
    "container.pods.list",
    "container.pods.portForward",
    "container.pods.proxy",
    "container.pods.updateStatus",
    "container.pods.update",
    "container.priorityClasses.create",
    "container.priorityClasses.delete",
    "container.priorityClasses.get",
    "container.priorityClasses.list",
    "container.priorityClasses.update",
    "container.replicaSets.create",
    "container.replicaSets.delete",
    "container.replicaSets.getScale",
    "container.replicaSets.getStatus",
    "container.replicaSets.get",
    "container.replicaSets.list",
    "container.replicaSets.updateScale",
    "container.replicaSets.updateStatus",
    "container.replicaSets.update",
    "container.replicationControllers.create",
    "container.replicationControllers.delete",
    "container.replicationControllers.getScale",
    "container.replicationControllers.getStatus",
    "container.replicationControllers.get",
    "container.replicationControllers.list",
    "container.replicationControllers.updateScale",
    "container.replicationControllers.updateStatus",
    "container.replicationControllers.update",
    "container.resourceQuotas.create",
    "container.resourceQuotas.delete",
    "container.resourceQuotas.getStatus",
    "container.resourceQuotas.get",
    "container.resourceQuotas.list",
    "container.resourceQuotas.updateStatus",
    "container.resourceQuotas.update",
    "container.roleBindings.create",
    "container.roleBindings.delete",
    "container.roleBindings.get",
    "container.roleBindings.list",
    "container.roleBindings.update",
    "container.roles.bind",
    "container.roles.create",
    "container.roles.delete",
    "container.roles.escalate",
    "container.roles.get",
    "container.roles.list",
    "container.roles.update",
    "container.runtimeClasses.create",
    "container.runtimeClasses.delete",
    "container.runtimeClasses.get",
    "container.runtimeClasses.list",
    "container.runtimeClasses.update",
    "container.scheduledJobs.create",
    "container.scheduledJobs.delete",
    "container.scheduledJobs.get",
    "container.scheduledJobs.list",
    "container.scheduledJobs.updateStatus",
    "container.scheduledJobs.update",
    "container.secrets.create",
    "container.secrets.delete",
    "container.secrets.get",
    "container.secrets.list",
    "container.secrets.update",
    "container.selfSubjectAccessReviews.create",
    "container.selfSubjectAccessReviews.list",
    "container.selfSubjectRulesReviews.create",
    "container.serviceAccounts.createToken",
    "container.serviceAccounts.create",
    "container.serviceAccounts.delete",
    "container.serviceAccounts.get",
    "container.serviceAccounts.list",
    "container.serviceAccounts.update",
    "container.services.create",
    "container.services.delete",
    "container.services.getStatus",
    "container.services.get",
    "container.services.list",
    "container.services.proxy",
    "container.services.updateStatus",
    "container.services.update",
    "container.statefulSets.create",
    "container.statefulSets.delete",
    "container.statefulSets.getScale",
    "container.statefulSets.getStatus",
    "container.statefulSets.get",
    "container.statefulSets.list",
    "container.statefulSets.updateScale",
    "container.statefulSets.updateStatus",
    "container.statefulSets.update",
    "container.storageClasses.create",
    "container.storageClasses.delete",
    "container.storageClasses.get",
    "container.storageClasses.list",
    "container.storageClasses.update",
    "container.storageStates.create",
    "container.storageStates.delete",
    "container.storageStates.getStatus",
    "container.storageStates.get",
    "container.storageStates.list",
    "container.storageStates.updateStatus",
    "container.storageStates.update",
    "container.storageVersionMigrations.create",
    "container.storageVersionMigrations.delete",
    "container.storageVersionMigrations.getStatus",
    "container.storageVersionMigrations.get",
    "container.storageVersionMigrations.list",
    "container.storageVersionMigrations.updateStatus",
    "container.storageVersionMigrations.update",
    "container.subjectAccessReviews.create",
    "container.subjectAccessReviews.list",
    "container.thirdPartyObjects.create",
    "container.thirdPartyObjects.delete",
    "container.thirdPartyObjects.get",
    "container.thirdPartyObjects.list",
    "container.thirdPartyObjects.update",
    "container.thirdPartyResources.create",
    "container.thirdPartyResources.delete",
    "container.thirdPartyResources.get",
    "container.thirdPartyResources.list",
    "container.thirdPartyResources.update",
    "container.tokenReviews.create",
    "container.updateInfos.create",
    "container.updateInfos.delete",
    "container.updateInfos.get",
    "container.updateInfos.list",
    "container.updateInfos.update",
    "container.validatingWebhookConfigurations.create",
    "container.validatingWebhookConfigurations.delete",
    "container.validatingWebhookConfigurations.get",
    "container.validatingWebhookConfigurations.list",
    "container.validatingWebhookConfigurations.update",
    "container.volumeAttachments.create",
    "container.volumeAttachments.delete",
    "container.volumeAttachments.getStatus",
    "container.volumeAttachments.get",
    "container.volumeAttachments.list",
    "container.volumeAttachments.updateStatus",
    "container.volumeAttachments.update",
    "container.volumeSnapshotClasses.create",
    "container.volumeSnapshotClasses.delete",
    "container.volumeSnapshotClasses.get",
    "container.volumeSnapshotClasses.list",
    "container.volumeSnapshotClasses.update",
    "container.volumeSnapshotContents.create",
    "container.volumeSnapshotContents.delete",
    "container.volumeSnapshotContents.getStatus",
    "container.volumeSnapshotContents.get",
    "container.volumeSnapshotContents.list",
    "container.volumeSnapshotContents.updateStatus",
    "container.volumeSnapshotContents.update",
    "container.volumeSnapshots.create",
    "container.volumeSnapshots.delete",
    "container.volumeSnapshots.getStatus",
    "container.volumeSnapshots.get",
    "container.volumeSnapshots.list",
    "container.volumeSnapshots.updateStatus",
    "container.volumeSnapshots.update",
    "recommender.containerDiagnosisInsights.get",
    "recommender.containerDiagnosisInsights.list",
    "recommender.containerDiagnosisInsights.update",
    "recommender.containerDiagnosisRecommendations.get",
    "recommender.containerDiagnosisRecommendations.list",
    "recommender.containerDiagnosisRecommendations.update",
    "recommender.locations.get",
    "recommender.locations.list",
    "recommender.networkAnalyzerGkeConnectivityInsights.get",
    "recommender.networkAnalyzerGkeConnectivityInsights.list",
    "recommender.networkAnalyzerGkeConnectivityInsights.update",
    "recommender.networkAnalyzerGkeIpAddressInsights.get",
    "recommender.networkAnalyzerGkeIpAddressInsights.list",
    "recommender.networkAnalyzerGkeIpAddressInsights.update",
    "resourcemanager.projects.get",
    "resourcemanager.projects.list",
  ],
  "stage": "BETA"
}
EOF

gcloud iam roles create GRITProvisioner --project=[projectID] --file=./grit-provisioner-role.json
```

</details>

where `[projectID]` is the ID of your Google Cloud project.

## Variables

You can use variables to control the behavior of the scenario.

Variables can be:

- **Required**: Variables must be provided when you define the module and do not have
  a default value.

- **Not required with a default value**: Variables are required for the scenario to work properly, but you
  can use the provided default values to experiment with the scenario.

- **Not required with no default value**: Variables are optional and don't need to be provided
  unless a specific use case requires changes in the default configuration.

- **Simple**: Variables use simple types as `string`, `number` or `boolean`.

- **Complex**: Variables are either lists, maps, or objects, or combination of these types.

<!-- begin: input vars -->
| Name                     | Type                                                     | Required? | Default value            | Description                                                                                |
|--------------------------|----------------------------------------------------------|-----------|--------------------------|--------------------------------------------------------------------------------------------|
| `name`                   | `string`                                                 | yes       |                          | Name of the deployment. Must be unique in scope of a Google Cloud project.                 |
| `labels`                 | `map(string)`                                            | no        |                          | Arbitrary list of `key=value` pairs that are added as labels to resources created by GRIT. |
| `google_region`          | `string`                                                 | yes       |                          | Google Cloud region chosen for the deployment.                                             |
| `google_zone`            | `string`                                                 | yes       |                          | Google Cloud zone chosen for the deployment.                                               |
| `gitlab_pat`             | `string`                                                 | yes       |                          | GitLab Personal Access token, which allows creating a Runner Registration token.           |
| `gitlab_project_id`      | `string`                                                 | yes       |                          | The GitLab project ID to register the Runner to.                                           |
| `node_count`             | `string`                                                 | no        | 1                        | The number of nodes to deploy the GKE cluster with.                                        |
| `runner_description`     | `string`                                                 | no        | "default GitLab Runner"  | The description of the deployed Runner, visible on the GitLab Runner configuration page.   |
| `subnet_cidr`            | `string`                                                 | no        | 10.0.0.0/10              | The CIDR for the subnetwork the GKE cluster will be deployed on.                           |
<!-- end: input vars -->

## Usage

### Terraform code

To use this scenario, we expect a couple of terraform providers to be
configured, so that they are implicitly inherited by the scenario (and its
modules) and it can make use of them:

- [`google`](https://registry.terraform.io/providers/hashicorp/google/latest/docs)
- [`gitlab`](https://registry.terraform.io/providers/gitlabhq/gitlab/latest/docs)
- [`kubectl`](https://registry.terraform.io/providers/alekc/kubectl/latest/docs)

The `gitlab` provider is authenticated by a Personal Access
Token (PAT), which we pass in with the terraform variable `gitlab_pat`. This is
a sensitive value, is marked as such in the `variable` block, and externalized.

The `kubectl` provider needs to authenticate against a
Kubernetes cluster, in our case the GKE cluster set up by the scenario. We only
know the coordinates and credentials of the cluster once it is deployed, thus
our scenario has output variables giving you access to the cluster's hostname
(`cluster_host`), its CA certificate (`cluster_ca_certificate`), and the Access
Token (`cluster_access_token`). These need to be wired into the
`kubectl` provider.

```terraform
# version setup
terraform {
  required_version = "~> 1.7"

  required_providers {
    kubectl = {
      source  = "alekc/kubectl"
      version = "~> 2.0"
    }
    google = {
      source  = "hashicorp/google"
      version = ">= 5.30.0"
    }
    gitlab = {
      source  = "gitlabhq/gitlab"
      version = ">= 17.0.0"
    }
  }
}

# input variables
variable "gitlab_pat" {
  description = "The personal access token for GitLab instance, to create the runner registration token"
  type        = string
  sensitive   = true
}

# provider setup
provider "google" {
  project = "my-google-project"
}

provider "gitlab" {
  token = var.gitlab_pat
}

provider "kubectl" {
  host                   = module.gke_runner.cluster_host
  cluster_ca_certificate = module.gke_runner.cluster_ca_certificate
  token                  = module.gke_runner.cluster_access_token
  load_config_file       = false
}

# the scenario
module "gke_runner" {
  source = "git::https://gitlab.com/gitlab-org/ci-cd/runner-tools/grit.git//scenarios/google/gke/operator"

  name   = "some-name"
  labels = {
    unit  = "some-unit"
    owner = "me"
  }

  google_region      = "europe-north1"
  google_zone        = "europe-north1-c"
  subnet_cidr        = "10.0.0.0/16"
  gitlab_pat         = var.gitlab_pat
  gitlab_project_id  = "123121213"
  runner_description = "my new GRIT runner"
}
```

### Terraform execution

Plan and deploy the example as [documented](../../../terraform.md).

### Inspecting created resources

This scenario created some resources on Google Cloud, most notably the GKE
cluster, and some on the GKE cluster itself. On the cluster you should see two
deployments in the new namespace `gitlab-runner-system`:

```shell
kubectl -n gitlab-runner-system get deployments
```

```terminal
NAME                                            READY   UP-TO-DATE   AVAILABLE   AGE
gitlab-runner-gitlab-runnercontroller-manager   1/1     1            1           16h
some-name-runner                                1/1     1            1           16h
```

Every time this runner is now used to run some task, you will see an ephemeral
`Pod` come into existence in this namespace, and vanish again as soon as its
work is done.
