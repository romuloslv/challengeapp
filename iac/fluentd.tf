data "kubectl_path_documents" "fluentd" {
  pattern = "./fluentd/*.yaml"
  vars    = { docker_image = "fluent/fluentd-kubernetes-daemonset:v1.9.2-debian-elasticsearch7-1.0" }
}

resource "kubectl_manifest" "fluentd" {
  for_each  = toset(sort(data.kubectl_path_documents.fluentd.documents))
  yaml_body = each.value

  depends_on = [
    helm_release.kubernetes_dashboard,
    helm_release.ingress-nginx,
    helm_release.prometheus,
    helm_release.grafana,
    helm_release.elasticsearch,
    helm_release.kibana,
    helm_release.dapr,
    kubectl_manifest.challengeapp
  ]
}