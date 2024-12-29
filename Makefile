include .env

APP_NAME=crypto-service
SWAGGER_DIR=./docs/swagger
LOGS_DIR=./logs

.PHONY: all
all: clean swagger build run test

.PHONY: swagger
swagger:
	swag init -q -g ./internal/transport/http/http.go -o $(SWAGGER_DIR)

.PHONY: build
build: clean swagger
	go build ./cmd/$(APP_NAME)

.PHONY: run
run: build
	./$(APP_NAME)

.PHONY: test
test:
	go test -v ./...

.PHONY: clean 
clean:
	rm $(APP_NAME) -f
	rm $(LOGS_DIR) -rf

.PHONY: docker-build
docker-build:
	docker compose build

.PHONY: docker-up
docker-up:
	docker compose up -d

.PHONY: docker-down
docker-down:
	docker compose down

