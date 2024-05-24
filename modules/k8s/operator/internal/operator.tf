resource "kubectl_manifest" "operator_resources" {
  for_each = {
    for n, m in local.operator_manifests : n => yamlencode(m)
  }
  yaml_body = each.value
  wait      = true
}
