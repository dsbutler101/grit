output "runner_token" {
  value       = module.gitlab.runner_token
  description = "The token used by GitLab Runner to get jobs"

  sensitive = true
}

output "url" {
  value       = module.gitlab.url
  description = "TODO"
}
