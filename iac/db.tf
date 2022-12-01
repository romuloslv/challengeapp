data "kubectl_path_documents" "db" {
  pattern = "./postgres/*.yaml"
}

resource "kubectl_manifest" "db" {
  for_each  = toset(sort(data.kubectl_path_documents.db.documents))
  yaml_body = each.value
}