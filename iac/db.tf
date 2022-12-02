data "kubectl_path_documents" "db" {
  pattern = "./postgres/*.yaml"
  vars    = { docker_image = "postgres:15.1-alpine3.16" }
}

resource "kubectl_manifest" "db" {
  for_each  = toset(sort(data.kubectl_path_documents.db.documents))
  yaml_body = each.value
}