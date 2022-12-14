resource "helm_release" "kubernetes_dashboard" {
  name             = "dashboard"
  chart            = "kubernetes-dashboard"
  repository       = "https://kubernetes.github.io/dashboard"
  namespace        = "lab-dashboard"
  create_namespace = true
  version          = "6.0.0"
  timeout          = 300

  set {
    name  = "service.type"
    value = "LoadBalancer"
  }

  set {
    name  = "protocolHttp"
    value = "true"
  }

  set {
    name  = "service.externalPort"
    value = 80
  }

  set {
    name  = "replicaCount"
    value = 1
  }

  set {
    name  = "rbac.clusterReadOnlyRole"
    value = "true"
  }

  depends_on = [
    google_container_cluster.main,
    google_container_node_pool.general
  ]
}

resource "helm_release" "ingress-nginx" {
  name             = "webproxy"
  chart            = "ingress-nginx"
  repository       = "https://kubernetes.github.io/ingress-nginx"
  namespace        = "lab-app"
  create_namespace = true
  version          = "4.4.0"
  timeout          = 300

  values = [file("${path.module}/templates/ingress-nginx-values.yaml")]

  depends_on = [
    google_container_cluster.main,
    google_container_node_pool.general
  ]
}

resource "helm_release" "prometheus" {
  name             = "reader"
  chart            = "prometheus"
  repository       = "https://prometheus-community.github.io/helm-charts"
  namespace        = "lab-monitoring"
  create_namespace = true
  version          = "19.0.1"
  timeout          = 300

  set {
    name  = "podSecurityPolicy.enabled"
    value = true
  }

  set {
    name  = "server.persistentVolume.enabled"
    value = false
  }

  set {
    name = "server\\.resources"
    value = yamlencode({
      limits = {
        cpu    = "200m"
        memory = "50Mi"
      }
      requests = {
        cpu    = "100m"
        memory = "30Mi"
      }
    })
  }

  depends_on = [
    google_container_cluster.main,
    google_container_node_pool.general
  ]
}

resource "helm_release" "grafana" {
  name             = "monitor"
  chart            = "grafana"
  repository       = "https://grafana.github.io/helm-charts"
  namespace        = kubernetes_namespace.grafana.metadata[0].name
  create_namespace = true
  version          = "6.44.11"
  timeout          = 300

  values = [
    templatefile("${path.module}/templates/grafana-values.yaml", {
      admin_existing_secret = kubernetes_secret.grafana.metadata[0].name
      admin_user_key        = "admin-user"
      admin_password_key    = "admin-password"
      prometheus_svc        = "${helm_release.prometheus.name}-prometheus-server"
      replicas              = 1
    })
  ]
}

resource "kubernetes_namespace" "grafana" {
  metadata { name = "lab-monitoring" }
}

resource "kubernetes_secret" "grafana" {
  metadata {
    name      = "grafana"
    namespace = kubernetes_namespace.grafana.metadata[0].name
  }

  data = {
    admin-user     = "admin"
    admin-password = random_password.grafana.result
  }

  depends_on = [
    google_container_cluster.main,
    google_container_node_pool.general,
    kubernetes_namespace.grafana
  ]
}

resource "random_password" "grafana" { length = 24 }

resource "helm_release" "elasticsearch" {
  name             = "exporter"
  chart            = "elasticsearch"
  repository       = "https://helm.elastic.co"
  namespace        = "lab-logging"
  create_namespace = true
  version          = "7.17.3"
  timeout          = 300

  set {
    name  = "replicas"
    value = "1"
  }

  depends_on = [
    google_container_cluster.main,
    google_container_node_pool.general
  ]
}

resource "helm_release" "kibana" {
  name             = "indexer"
  chart            = "kibana"
  repository       = "https://helm.elastic.co"
  namespace        = "lab-logging"
  create_namespace = true
  version          = "7.17.3"
  timeout          = 300

  depends_on = [
    google_container_cluster.main,
    google_container_node_pool.general,
    helm_release.prometheus,
    helm_release.elasticsearch,
    helm_release.dapr
  ]
}

resource "helm_release" "dapr" {
  name             = "runtime"
  chart            = "dapr"
  repository       = "https://dapr.github.io/helm-charts/"
  namespace        = "lab-logging"
  create_namespace = true
  version          = "1.9.4"
  timeout          = 300

  set {
    name  = "global.logAsJson"
    value = true
  }

  depends_on = [
    google_container_cluster.main,
    google_container_node_pool.general
  ]
}

resource "helm_release" "cert-manager" {
  name             = "certificate"
  chart            = "cert-manager"
  repository       = "https://charts.jetstack.io"
  namespace        = "lab-observability"
  create_namespace = true
  wait             = true
  version          = "v1.10.1"
  timeout          = 300

  set {
    name  = "installCRDs"
    value = true
  }

  depends_on = [
    google_container_cluster.main,
    google_container_node_pool.general
  ]
}