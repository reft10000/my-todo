.PHONY: run build test lint fmt clean

run:
	go run ./cmd/api

build:
	go build -o bin/api ./cmd/api

test:
	go test ./...

test-v:
	go test -v ./...

lint:
	golangci-lint run ./...

fmt:
	gofmt -w .

clean:
	rm -rf bin/

tidy:
	go mod tidy
