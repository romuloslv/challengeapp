# Build
.PHONY: build
build:
	@echo "Building local image..."
	@docker build -t api:1.0 .

# Stack
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