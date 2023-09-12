output "runner_token" {
  value       = gitlab_user_runner.primary.token
  description = "The token used by GitLab Runner to get jobs"
}

