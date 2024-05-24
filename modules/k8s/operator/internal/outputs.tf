output "namespace" {
  value = one([
    for id, res in kubectl_manifest.operator_resources :
    res.name if res.kind == "Namespace"
  ])
}
