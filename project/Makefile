SHELL := /bin/bash
.ONESHELL:
.DEFAULT_GOAL := help
OS := $(strip $(shell uname -o))

FRONT_END_BINARY := front-end-app
BROKER_BINARY := broker-app
AUTH_BINARY := auth-app
LOGGER_BINARY := logger-app
MAIL_BINARY := mail-app
LISTENER_BINARY := listener-app


.PHONY: up
up:   ## Start all containers in the background without forcing build
	@echo "Starting containers..."
	docker-compose up -d
	@echo "Containers started!"


.PHONY: up-build
up-build: build-broker build-auth build-logger build-mail build-listener build-front ## Stops docker-compose if running, builds all projects, starts docker compose
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


.PHONY: build-mail
build-mail:   ## Build the mail binary as a linux executable
	@echo "Building mail binary..."
	cd ../mail-service && env GOOS=linux CGO_ENABLED=0 go build -o ${MAIL_BINARY} ./cmd/api
	@echo "Done!"


.PHONY: build-listener
build-listener:   ## Build the listener binary as a linux executable
	@echo "Building listener binary..."
	cd ../listener-service && env GOOS=linux CGO_ENABLED=0 go build -o ${LISTENER_BINARY} .
	@echo "Done!"


.PHONY: build-front
build-front:   ## Build the front end binary as a linux executable
	@echo "Building front end binary..."
	cd ../front-end-service && env GOOS=linux CGO_ENABLED=0 go build -o ${FRONT_END_BINARY} ./cmd/web
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


TAG ?= latest
SERVICE ?= "logger"


.PHONY: build-push
build-push:   ## Build and push the given service (with the given tag) to Docker Hub
	@echo "Building and pushing Docker image ${SERVICE}-service with tag ${TAG}..."
	docker build -f ../${SERVICE}-service/${SERVICE}-service.dockerfile -t blacksound1/${SERVICE}-service:${TAG} ../${SERVICE}-service
	docker push blacksound1/${SERVICE}-service:${TAG}
	@echo "Docker image built and pushed!"


.PHONY: swarm-init
swarm-init:   ## Initialize the Docker Swarm
	@echo "Initializing Docker Swarm..."
	docker swarm init
	@echo "Docker Swarm initialized!"


.PHONY: swarm-deploy
swarm-deploy:   ## Deploy the Docker Swarm
	@echo "Deploying Docker Swarm..."
	docker stack deploy -c swarm.yml myapp
	@echo "Docker Swarm deployed!"


.PHONY: swarm-undeploy
swarm-undeploy:   ## Undeploy the Docker Swarm
	@echo "Undeploying Docker Swarm..."
	docker stack rm myapp
	@echo "Docker Swarm undeployed!"


.PHONY: swarm-leave
swarm-leave:   ## Leave the Docker Swarm
	@echo "Leaving Docker Swarm..."
	docker swarm leave --force
	@echo "Docker Swarm left!"


SCALE ?= 1

.PHONY: swarm-scale
swarm-scale:   ## Scale the a service in the Docker Swarm by a given number
	@echo "Scaling service ${SERVICE}..."
	docker service scale myapp_${SERVICE}-service=${SCALE}
	@echo "Service scaled!"


.PHONY: swarm-update
swarm-update:   ## Update and image in the Docker Swarm
	@echo "Updating service ${SERVICE}..."
	docker service update --image blacksound1/${SERVICE}-service:${TAG} myapp_${SERVICE}-service
	@echo "Service updated!"


.PHONY: help
help:   ## Show this help
	@echo -e "\nCommands:\n"
	@egrep '^[a-zA-Z_-]+:.*?## .*' Makefile | sort |
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
