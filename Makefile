.PHONY: run build docker-up docker-down docker-restart

# Run the Go application
run:
	go run cmd/main.go

# Build the Go application
build:
	go build -o bin/server cmd/main.go

# Start PostgreSQL container
docker-up:
	docker compose up -d

# Stop PostgreSQL container
docker-down:
	docker compose down

# Restart PostgreSQL container
docker-restart:
	docker compose down && docker compose up -d

# View PostgreSQL logs
docker-logs:
	docker compose logs -f postgres

# Tidy up Go modules
tidy:
	go mod tidy
