.PHONY: build test lint run

build:
	mkdir -p bin
	go build -o bin/hexlet-path-size ./cmd/hexlet-path-size/main.go

test:
	go test -v ./...

lint:
	golangci-lint run

run: build
	./bin/hexlet-path-size $(ARGS)