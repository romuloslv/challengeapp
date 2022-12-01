locals { region = "southamerica-east1-a" }

variable "kubernetes_name" {
  type        = string
  description = "Enter your GKE cluster name"
}