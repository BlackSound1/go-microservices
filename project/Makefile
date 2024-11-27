SHELL := /bin/bash
.ONESHELL:
.DEFAULT_GOAL := help
OS := $(strip $(shell uname -o))

FRONT_END_BINARY := frontend-app
BROKER_BINARY := broker-app
AUTH_BINARY := auth-app
LOGGER_BINARY := logger-app


.PHONY: up
up:   ## Start all containers in the background without forcing build
	@echo "Starting containers..."
	docker-compose up -d
	@echo "Containers started!"


.PHONY: up-build
up-build: build-broker build-auth build-logger  ## Stops docker-compose if running, builds all projects, starts docker compose
	@echo "Stopping docker images, if running..."
	docker-compose down
	@echo "Building, when required, and starting docker images..."
	docker-compose up --build -d
	@echo "Images built and started!"


.PHONY: down
down:   ## Stop docker compose
	@echo "Stopping docker compose..."
	docker-compose down
	@echo "Docker compose stopped!"


.PHONY: build-broker
build-broker:   ## Build the broker binary as a linux executable
	@echo "Building broker binary..."
	cd ../broker-service && env GOOS=linux CGO_ENABLED=0 go build -o ${BROKER_BINARY} ./cmd/api
	@echo "Done!"


.PHONY: build-auth
build-auth:   ## Build the authentication binary as a linux executable
	@echo "Building auth binary..."
	cd ../auth-service && env GOOS=linux CGO_ENABLED=0 go build -o ${AUTH_BINARY} ./cmd/api
	@echo "Done!"


.PHONY: build-logger
build-logger:   ## Build the logger binary as a linux executable
	@echo "Building logger binary..."
	cd ../logger-service && env GOOS=linux CGO_ENABLED=0 go build -o ${LOGGER_BINARY} ./cmd/api
	@echo "Done!"


.PHONY: build-front
build-front:   ## Build the front end binary as a windows executable
	@echo "Building front end binary..."
	cd ../front-end && env CGO_ENABLED=0 go build -o ${FRONT_END_BINARY} ./cmd/web
	@echo "Done!"


.PHONY: start
start: build-front   ## Start the front end
	@echo "Starting front end..."
	cd ../front-end && ./${FRONT_END_BINARY} &


.PHONY: stop
stop:   ## Stop the front end
	@echo "Stopping front end..."
	@if test "$(OS)" = "GNU/Linux"; then \
		pkill -SIGTERM -f "./${FRONT_END_BINARY}"
	@elif test "$(OS)" = "Msys"; then \
		taskkill -IM "${FRONT_END_BINARY}" -F 
	@fi
	@echo "Front end stopped!"


.PHONY: help
help:   ## Show this help
	@echo -e "\nCommands:\n"
	@egrep '^[a-zA-Z_-]+:.*?## .*' Makefile | sort |
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
