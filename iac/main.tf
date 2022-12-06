resource "google_container_cluster" "main" {
  name     = var.kubernetes_name
  location = local.region

  remove_default_node_pool = true
  initial_node_count       = 1
}

resource "google_container_node_pool" "general" {
  name               = "default"
  cluster            = google_container_cluster.main.id
  initial_node_count = 3

  node_config {
    preemptible  = true
    machine_type = "e2-standard-2"
  }

  management {
    auto_repair  = true
    auto_upgrade = true
  }
}