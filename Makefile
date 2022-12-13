MAKEFLAGS += --silent

.PHONY: run
run:
	@echo "Running local using docker run"
	@./cmd/app/docker-run.sh

.PHONY:	prod
prod:
	@echo "Running local prod using docker-compose"
	@docker compose -f stack.yaml down -v
	@docker compose -f stack.yaml up --build -d

.PHONY: dev
dev:
	@echo "Running tests local using docker-compose"
	@docker compose -f stack.yaml down -v
	@docker compose -f stack.yaml -f stack.dev.yaml up -d

.PHONY: tests
tests:
	@sleep 15 && cd api/accounts && go test -v

.PHONY:	stop
stop:
	@docker compose -f stack.yaml down -v 2> /dev/null
	@docker rm -f api db pg 2> /dev/null

.PHONY: alltf
alltf: terraform-login terraform-validation terraform-apply-cluster terraform-apply-pkgs

.PHONY: terraform-login
terraform-login:
	@gcloud auth application-default login
	@gcloud auth application-default set-quota-project $(project_name)
	@gcloud auth login
	@gcloud config set project $(project_name)

.PHONY: terraform-validation
terraform-validation:
	@terraform -chdir=iac init && terraform -chdir=iac validate && terraform -chdir=iac fmt

.PHONY: terraform-apply-cluster
terraform-apply-cluster:
	@terraform -chdir=iac apply -var kubernetes_name=$(cluster_name) \
			  -target=google_container_cluster.main \
			  -target=google_container_node_pool.general \
			  -compact-warnings -auto-approve

.PHONY: terraform-apply-pkgs
terraform-apply-pkgs:
	@gcloud container clusters get-credentials $(cluster_name) \
					--zone southamerica-east1-a \
					--project $(project_name)
	@terraform -chdir=iac apply -var kubernetes_name=$(cluster_name) -auto-approve

.PHONY: terraform-destroy
terraform-destroy:
	@terraform -chdir=iac destroy -var kubernetes_name=$(cluster_name) -auto-approve