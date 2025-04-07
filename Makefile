.PHONY: test coverage lint

test:
	go test -race -cover ./...

coverage:
	go test -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
	rm coverage.out

lint:
	golangci-lint run --timeout 5m
