.PHONY: run migrate-up migrate-down build

run:
	go run ./cmd/api/main.go

migrate-up:
	migrate -path ./migrations -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path ./migrations -database "$(DATABASE_URL)" down

build:
	go build -o bin/api ./cmd/api/main.go
