data "kubectl_path_documents" "challengeapp" {
  pattern = "./challengeapp/*.yaml"
  vars    = { docker_image = "romuloslv/challengeapp:1.0" }
}

resource "kubectl_manifest" "challengeapp" {
  for_each  = toset(sort(data.kubectl_path_documents.challengeapp.documents))
  yaml_body = each.value

  depends_on = [
    helm_release.kubernetes_dashboard,
    helm_release.ingress-nginx,
    helm_release.prometheus,
    helm_release.grafana,
    helm_release.elasticsearch,
    helm_release.kibana,
    helm_release.dapr
  ]
}