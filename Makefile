LIMESTONE_BINARY=limestoneApp
## up: starts all containers in the background without forcing build
up:
	@echo "Starting Docker images..."
	docker-compose up -d
	@echo "Docker images started!"
## up_build: starts containers with forced build (use when made changes to code)
up_build: build_limestone
	@echo "Stopping docker images (if running...)"
	docker-compose down
	@echo "Building (when required) and starting docker images..."
	docker-compose up --build -d
	@echo "Docker images built and started!"

## down: stop docker compose
down:
	@echo "Stopping docker compose..."
	docker-compose down
	@echo "Done!"

## build_limestone: builds the limestone binary as a linux executable
build_limestone:
	@echo "Building limestone binary..."
	env GOOS=linux CGO_ENABLED=0 go build -o ${LIMESTONE_BINARY} .
	@echo "Done!"