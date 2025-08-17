# Test commands
.PHONY: test test-unit test-integration test-coverage

test: test-unit test-integration

test-unit:
	@echo "Running unit tests..."
	go test -v ./internal/...

test-integration:
	@echo "Running integration tests..."
	go test -v ./tests/integration/...

test-coverage:
	@echo "Running tests with coverage..."
	go test -coverprofile=coverage.out ./internal/...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-watch:
	@echo "Running tests in watch mode..."
	find . -name "*.go" | entr -r make test-unit