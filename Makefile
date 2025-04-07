.PHONY: test coverage build run docker-up docker-down docker-build lint

GO_CMD=go
APP_NAME=chatheon

build:
	$(GO_CMD) build -o bin/$(APP_NAME) ./cmd/$(APP_NAME)

run: build
	./bin/$(APP_NAME)

test:
	$(GO_CMD) test -race -cover ./...

coverage:
	$(GO_CMD) test -race -coverprofile=coverage.out ./...
	$(GO_CMD) tool cover -html=coverage.out

docker-up:
	docker-compose up --build

docker-down:
	docker-compose down

lint:
	golangci-lint run --timeout 5m
