data "gitlab_project" "grit" {
  path_with_namespace = "gitlab-org/ci-cd/runner-tools/grit"
}

data "gitlab_project" "grit-e2e" {
  path_with_namespace = "gitlab-org/ci-cd/runner-tools/grit-e2e"
}
