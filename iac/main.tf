resource "google_container_cluster" "main" {
  name     = var.kubernetes_name
  location = local.region

  node_pool { name = "builtin" }
  lifecycle { ignore_changes = [node_pool] }
}

resource "google_container_node_pool" "general" {
  name               = "general"
  cluster            = google_container_cluster.main.id
  initial_node_count = 3

  node_config {
    preemptible  = false
    machine_type = "e2-standard-2"
  }

  management {
    auto_repair  = true
    auto_upgrade = true
  }
}