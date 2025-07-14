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
	export $(shell grep -v '^#' .env | xargs) ; \
	export DOCKER_HOST=$$DOCKER_HOST ; \
	go run ./cmd/orchestrator dashboard 

monitoring:
	export $(shell grep -v '^#' .env | xargs) ; \
	export DOCKER_HOST=$$DOCKER_HOST ; \
	go run ./cmd/orchestrator monitor all

status: 
	export $(shell grep -v '^#' .env | xargs) ; \
	export DOCKER_HOST=$$DOCKER_HOST ; \
	go run ./cmd/orchestrator status all 

start-all: 
	export $(shell grep -v '^#' .env | xargs) ; \
	export DOCKER_HOST=$$DOCKER_HOST ; \
	go run ./cmd/orchestrator start all 	
	