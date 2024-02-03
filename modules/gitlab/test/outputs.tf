output "runner_token" {
  description = "The token used by GitLab Runner to get jobs."
  value       = module.gitlab.runner_token

  sensitive = true
}

output "url" {
  description = "The URL where the runner was registered."
  value       = module.gitlab.url
}
