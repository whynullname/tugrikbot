.PHONY: migrate-up migrate-down up down run stop prod

up:
	docker compose up -d

down:
	docker compose down

migrate-up:
	go run ./cmd/migrate up

migrate-down:
	go run ./cmd/migrate down

run:
	docker compose up -d --wait db
	go run ./cmd/migrate up
	go run ./cmd

stop:
	docker compose down

prod:
	docker compose up --build