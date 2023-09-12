resource "helm_release" "gitlab-runner" {
  name       = "gitlab-runner"
  repository = "https://charts.gitlab.io"
  chart      = "gitlab-runner"
  set {
    name  = "gitlabUrl"
    value = var.gitlab_url
  }
  set {
    name  = "rbac.create"
    value = "true"
  }
  set {
    name  = "runnerToken"
    value = gitlab_user_runner.primary.token
  }
}
