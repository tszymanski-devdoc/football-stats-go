# Install swag CLI tool
install-swag:
	go install github.com/swaggo/swag/cmd/swag@latest

# Generate Swagger documentation
swagger:
	swag init -g cmd/api/main.go -o docs --parseDependency --parseInternal

# Run with auto-generated swagger docs
run: swagger
	go run cmd/api/main.go

# Build the application
build: swagger
	go build -o bin/api cmd/api/main.go

# Run tests
test:
	go test -v -race ./...

# Install dependencies
deps:
	go mod download
	go mod tidy

# Run database migrations (up)
migrate-up:
	go run cmd/migrate/main.go -direction up

# Rollback database migrations (down)
migrate-down:
	go run cmd/migrate/main.go -direction down

# Clean generated files
clean:
	rm -rf docs/
	rm -rf bin/

.PHONY: install-swag swagger run build test deps migrate-up migrate-down clean
