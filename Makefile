.PHONY: test test-integration lint build run sec

# Unit tests
test:
	go test ./... -v -short

# Integration tests
test-integration:
	go test ./... -v -tags=integration

# Linting
lint:
	golangci-lint run

# Security (vulnerability scan)
sec:
	gosec ./...

# Build the binary
build:
	go build -o dockpilot ./cmd/orchestrator

# Run dashboard
run:
	go run ./cmd/orchestrator/main.go dashboard 