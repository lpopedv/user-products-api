
COMPOSE=docker compose
SERVICE_APP=user_products
SERVICE_DB=user_products_postgres_dev
VOLUME_DB=postgres_data

# Docker commands
.PHONY: up up-detached restart log-app log-db down reset-db

up:
	$(COMPOSE) up

up-detached:
	$(COMPOSE) up -d

restart:
	$(COMPOSE) restart

log-app:
	docker logs -f $(SERVICE_APP)

log-db:
	docker logs -f $(SERVICE_DB)

down:
	$(COMPOSE) down

reset-db:
	$(COMPOSE) down
	@if docker volume ls -q | grep -q '^$(VOLUME_DB)$$'; then \
		docker volume rm $(VOLUME_DB); \
	fi
	$(COMPOSE) up -d

# Go commands
.PHONY: run-user-products test

run-user-products:
	go run ./main.go

test:
	go test ./...
