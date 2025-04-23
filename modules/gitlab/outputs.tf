output "enabled" {
  value = tobool(true)
}

output "runner_token" {
  description = "The token used by GitLab Runner to get jobs."
  value       = tostring(gitlab_user_runner.primary.token)
  sensitive   = true
}

output "url" {
  description = "The URL where the runner was registered."
  value       = tostring(var.url)
}
