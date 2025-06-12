output "skew" {
  description = "The determined skew as 0, 1 or 2."
  value       = tonumber(local.skew)
}

output "runner_version" {
  description = "The determined runner version as X.Y.Z."
  value       = tostring(local.runner_version)
}