.PHONY: build test coverage lint lint-fix run

BINARY := bin/hexlet-path-size
ifeq ($(OS),Windows_NT)
BINARY := $(BINARY).exe
endif

build:
	go build -o $(BINARY) ./cmd/hexlet-path-size

test:
	go test -v ./...

coverage:
	go test -coverpkg=./... -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

lint:
	golangci-lint run

lint-fix:
	golangci-lint fmt
	golangci-lint run --fix

run: build
	./$(BINARY) $(ARGS)
