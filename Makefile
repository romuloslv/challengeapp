.PHONY: build
build:
	@echo "Building local image..."
	@docker build -t api:1.0 .

.PHONY:	stop
stop:
	@docker compose -f stack.yaml down -v

.PHONY:	prod
prod:
	@docker compose -f stack.yaml down -v
	@docker compose -f stack.yaml up --build

.PHONY: dev
dev:
	@docker compose -f stack.yaml down -v
	@docker compose -f stack.yaml -f stack.dev.yaml up

.PHONY: all-tf
all-tf: terraform-login terraform-validation terraform-apply-cluster terraform-apply-pkgs

.PHONY: terraform-login
terraform-login:
	gcloud auth application-default login && gcloud auth application-default set-quota-project $(project_name)
	gcloud auth login && gcloud config set project $(project_name)

.PHONY: terraform-validation
terraform-validation:
	terraform init && terraform validate && terraform fmt

.PHONY: terraform-apply-cluster
terraform-apply-cluster:
	terraform apply -var kubernetes_name=$(cluster_name) \
					-target=google_container_cluster.main \
					-target=google_container_node_pool.general \
					-compact-warnings -auto-approve

.PHONY: terraform-apply-pkgs
terraform-apply-pkgs:
	gcloud container clusters get-credentials $(cluster_name) \
					--zone southamerica-east1-a \
					--project $(project_name)
	terraform apply -var kubernetes_name=$(cluster_name) -auto-approve

.PHONY: terraform-destroy
terraform-destroy:
	terraform destroy -var kubernetes_name=$(cluster_name) -auto-approve