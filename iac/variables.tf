locals { region = "us-east4-a" }

variable "kubernetes_name" {
  type        = string
  description = "Enter your GKE cluster name"
}