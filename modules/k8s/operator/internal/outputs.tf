output "namespace" {
  value = one([
    for id, res in kubectl_manifest.operator_resources :
    res.name if res.kind == "Namespace"
  ])
}

output "supported_operator_versions" {
  value = local.supported_versions_info
}
