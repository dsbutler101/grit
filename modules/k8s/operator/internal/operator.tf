locals {
  supported_versions = {
    for file in sort(fileset(path.module, "versions/**/manifests.yaml")) :
    basename(dirname(file)) => "${path.module}/${file}"
  }

  manifests_file = (
    var.override_manifests == ""
    ? local.supported_versions[var.operator_version]
    : var.override_manifests
  )

  operator_manifests = {
    # - split the multi-doc yaml
    # - parse it
    # - pull out some data (eg. namespace, ...)
    # - return map of name => parsed yaml
    for resource in [
      for yaml in [
        for doc in split("\n---\n", file(local.manifests_file)) : yamldecode(doc)
      ] :
      {
        apiVersion = yaml.apiVersion
        kind       = yaml.kind
        namespace  = lookup(yaml.metadata, "namespace", "_cluster_scoped_")
        name       = yaml.metadata.name
        full       = yaml
      }
    ] :
    "${resource.apiVersion}::${resource.kind}::${resource.namespace}::${resource.name}" => resource.full
  }
}

resource "kubectl_manifest" "operator_resources" {
  for_each = {
    for n, m in local.operator_manifests : n => yamlencode(m)
  }
  yaml_body = each.value
  wait      = true
}
