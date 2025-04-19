.PHONY: test coverage build run docker-up docker-down docker-build lint

GO_CMD=go
APP_NAME=chatheon
MOCK_GEN=mockery

build:
	$(GO_CMD) build -o bin/$(APP_NAME) ./cmd/$(APP_NAME)

mocks:
	$(MOCK_GEN)

run: build
	./bin/$(APP_NAME)

test:
	$(GO_CMD) test -race -cover ./...

coverage:
	$(GO_CMD) test -race -coverprofile=coverage.out ./...
	$(GO_CMD) tool cover -html=coverage.out

docker-build:
	docker build -t $(APP_NAME) .

docker-up:
	docker-compose up --build

docker-down:
	docker-compose down

lint:
	golangci-lint run --timeout 5m

format:
	goimports -w .
