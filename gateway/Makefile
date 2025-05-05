CONFIG_PATH ?= ./config/.env
SERVER_PORT ?= 8080

.PHONY: generate-api generate-models generate-swagger run-server stop-server build test-app lint

generate-api: generate-models generate-swagger

generate-models:
	oapi-codegen -config api/v1/oapi-models-codegen.yaml api/v1/openapi.yaml

generate-swagger:
	oapi-codegen -config api/v1/oapi-swagger-codegen.yaml api/v1/openapi.yaml


# Run server
run-server: stop-server
	@echo "Running the server..."
	@CONFIG_PATH=$(CONFIG_PATH) go run github.com/rshelekhov/msa-messenger/gateway/cmd &
	@sleep 5 # Wait for the server to start
	@while ! nc -z localhost $(SERVER_PORT); do \
		echo "Waiting for server to be ready..."; \
		sleep 1; \
	done
	@echo "Server is running with PID $$(lsof -t -i :$(SERVER_PORT))."

# Stop server
stop-server:
	@echo "Stopping the server..."
	@PID=$$(lsof -t -i :$(SERVER_PORT)); \
    	if [ -n "$$PID" ]; then \
    		kill $$PID; \
    		echo "Server stopped."; \
    	else \
    		echo "No server is running on port $(SERVER_PORT)."; \
    	fi

build:
	go build -v ./cmd/app

# Run tests
test-app: run-server
	@echo "Running tests..."
	@go test -v -timeout 60s -parallel=1 ./...
	@echo "Tests completed."

# Run linters
lint:
	@echo "Running linters..."
	golangci-lint run --fix
	@echo "Linters completed."

.DEFAULT_GOAL := build