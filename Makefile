# Define variables
DOCKER_COMPOSE_FILE=docker-compose.yml
DOCKER_SERVICE=db
APP_NAME=bank-api

# Define Go version
GO_VERSION=1.22

# Ensure the correct Go version is used
ensure-go-version:
	@echo "Ensuring Go version $(GO_VERSION) is used..."
	@if [ "$(shell go version | grep -o 'go[0-9].[0-9]\{1,2\}')" != "go$(GO_VERSION)" ]; then \
		echo "Go version mismatch. Please use Go $(GO_VERSION)."; \
		exit 1; \
	fi

# Clean up Docker containers
clean:
	@echo "Stopping and removing Docker containers..."
	docker-compose -f $(DOCKER_COMPOSE_FILE) down

# Start fresh Docker containers
start:
	@echo "Starting Docker containers..."
	docker-compose -f $(DOCKER_COMPOSE_FILE) up -d

# Build the Go application
build: ensure-go-version
	@echo "Building the Go application..."
	go build -o $(APP_NAME) main.go

# Run the Go application
run:
	@echo "Running the Go application..."
	./$(APP_NAME)

# Complete process: clean, start, build, and run
all: clean start build run

.PHONY: clean start build run all ensure-go-version