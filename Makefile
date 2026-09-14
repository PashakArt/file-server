-include .env
export

.PHONY: start deps-up deps-down migrate migrate-down

start:
	go run internal/cmd/main.go

deps-up:
	docker compose up -d

deps-down:
	docker compose down

migrate:
	docker compose run --rm migration -path=/migrations -database "postgres://${DB_USER}:${DB_PASSWORD}@postgres:5432/${DB_NAME}?sslmode=disable" up

migrate-down:
	docker compose run --rm migration -path=/migrations -database "postgres://${DB_USER}:${DB_PASSWORD}@postgres:5432/${DB_NAME}?sslmode=disable" down 1