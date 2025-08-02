APP_NAME=limestone-app

build:
	docker build -t $(APP_NAME) .

run:
	docker-compose up

down:
	docker-compose down

rebuild:
	docker-compose down
	docker-compose build
	docker-compose up

proto:
	buf generate

clean:
	docker system prune -f
