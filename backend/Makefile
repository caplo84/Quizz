.PHONY: run-dev run-staging run-prod build test docker-up docker-down

# Development
run-dev:
	@echo "Starting development server..."
	APP_ENV=development go run ./cmd/api/

# Staging
run-staging:
	@echo "Starting staging server..."
	APP_ENV=staging go run ./cmd/api/

# Production
run-prod:
	@echo "Starting production server..."
	APP_ENV=production go run ./cmd/api/

# Build for different environments
build-dev:
	APP_ENV=development go build -o bin/quiz-api-dev ./cmd/api/

build-staging:
	APP_ENV=staging go build -o bin/quiz-api-staging ./cmd/api/

build-prod:
	APP_ENV=production go build -ldflags="-s -w" -o bin/quiz-api-prod ./cmd/api/

# Test with different environments
test-dev:
	APP_ENV=development go test ./...

test-staging:
	APP_ENV=staging go test ./...

# Docker environments
docker-dev:
	docker-compose -f docker-compose.dev.yml up -d

docker-staging:
	docker-compose -f docker-compose.staging.yml up -d

docker-prod:
	docker-compose -f docker-compose.prod.yml up -d