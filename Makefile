include .env
export

.PHONY: run build test lint fmt clean tidy

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

.PHONY: ent-generate
ent-generate:
	go generate ./internal/infra/ent/...

.PHONY: atlas-diff atlas-apply
atlas-diff:
	atlas migrate diff --env local

atlas-apply:
	atlas migrate apply --env local
