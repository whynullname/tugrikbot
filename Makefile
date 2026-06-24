.PHONY: migrate-up migrate-down up down

up:
	docker compose up -d

down:
	docker compose down

migrate-up:
	go run ./cmd/migrate up

migrate-down:
	go run ./cmd/migrate down