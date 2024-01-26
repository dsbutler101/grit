output "runner_token" {
  value = gitlab_user_runner.primary.token
}

output "url" {
  value = var.url
}
