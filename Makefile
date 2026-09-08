.PHONY: start deps-up deps-down

start:
	go run internal/cmd/main.go

deps-up:
	docker compose up -d

deps-down:
	docker compose down