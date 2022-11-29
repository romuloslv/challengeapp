locals { region = "southamerica-east1-a" }

variable "kubernetes_name" {
  type        = string
  description = "Enter your GKE cluster name"
}

variable "dashboard_endpoint" {
  description = "Dashboard endpoint"
  type        = string
  default     = "TO CONNECT TO DASHBOARD: kubectl get svc -n lab-dashboard | awk '{print $4}'"
}

variable "prometheus_endpoint" {
  description = "Prometheus endpoint"
  type        = string
  default     = "TO CONNECT TO PROMETHEUS: kubectl port-forward $(kubectl get pods -l=app=prometheus -o name -n lab-monitoring | tail -n1) 9090 -n lab-monitoring"
}

variable "grafana_endpoint" {
  description = "Grafana endpoint"
  type        = string
  default     = "TO CONNECT TO GRAFANA: kubectl port-forward $(kubectl get pods -l=app.kubernetes.io/instance=monitor -o name -n lab-monitoring) 3000 -n lab-monitoring"
}

variable "kibana_endpoint" {
  description = "Kibana endpoint"
  type        = string
  default     = "TO CONNECT TO KIBANA: kubectl port-forward $(kubectl get pods -l=app=kibana -o name -n lab-logging) 5601 -n lab-logging"
}

variable "elasticsearch_endpoint" {
  description = "Elasticsearch endpoint"
  type        = string
  default     = "TO CONNECT TO ELASTICSEARCH: kubectl port-forward $(kubectl get pods -l=app=elasticsearch-master -o name -n lab-logging) 9200 -n lab-logging"
}