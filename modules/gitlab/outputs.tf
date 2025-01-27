output "runner_token" {
  description = "The token used by GitLab Runner to get jobs."
  value       = gitlab_user_runner.primary.token

  sensitive = true
}

output "url" {
  description = "The URL where the runner was registered."
  value       = var.url
}
