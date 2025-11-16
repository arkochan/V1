.PHONY: build run dev lint test tidy

build:
	@go build -o bin/server ./cmd/api

run:
	@go run ./cmd/api

dev:
	@air


lint:
	@golangci-lint run

test:
	@go test ./...

tidy:
	@go mod tidy
