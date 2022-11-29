terraform {
  required_version = ">= 1.0.0"

  backend "gcs" {
    bucket = "poc-from-gke-tf-state"
    prefix = "state"
  }

  required_providers {
    google     = { source = "hashicorp/google" }
    kubernetes = { source = "hashicorp/kubernetes" }
    helm       = { source = "hashicorp/helm" }
    kubectl    = { source = "gavinbunney/kubectl" }
    random     = { source = "hashicorp/random" }
  }
}

provider "helm" {
  kubernetes {
    host                   = "https://${google_container_cluster.main.endpoint}"
    client_certificate     = base64decode(google_container_cluster.main.master_auth.0.client_certificate)
    client_key             = base64decode(google_container_cluster.main.master_auth.0.client_key)
    cluster_ca_certificate = base64decode(google_container_cluster.main.master_auth.0.cluster_ca_certificate)
  }
}

provider "kubernetes" {
  host                   = "https://${google_container_cluster.main.endpoint}"
  client_certificate     = base64decode(google_container_cluster.main.master_auth.0.client_certificate)
  client_key             = base64decode(google_container_cluster.main.master_auth.0.client_key)
  cluster_ca_certificate = base64decode(google_container_cluster.main.master_auth.0.cluster_ca_certificate)

}

provider "kubectl" {
  apply_retry_count      = 15
  host                   = "https://${google_container_cluster.main.endpoint}"
  client_certificate     = base64decode(google_container_cluster.main.master_auth.0.client_certificate)
  client_key             = base64decode(google_container_cluster.main.master_auth.0.client_key)
  cluster_ca_certificate = base64decode(google_container_cluster.main.master_auth.0.cluster_ca_certificate)
}