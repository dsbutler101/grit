output "tags" {
  description = "List of tag names from the GitLab repository"
  value       = jsondecode(data.http.gitlab_tags.response_body)[*].name
}
