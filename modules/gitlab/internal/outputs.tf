output "runner_token" {
  value     = gitlab_user_runner.primary.token
  sensitive = true
}

output "url" {
  value = var.url
}
