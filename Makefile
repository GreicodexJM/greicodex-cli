# Makefile for AmableSagitta

BINARY_NAME=grei

.PHONY: all deps clean build test coverage run

all: build

deps:
	@echo "==> Installing dependencies..."
	go mod tidy

clean:
	@echo "==> Cleaning..."
	go clean
	rm -f $(BINARY_NAME)-cli

build: deps
	@echo "==> Building binary..."
	go build -o $(BINARY_NAME)-cli ./cmd/$(BINARY_NAME)/

test:
	@echo "==> Running tests..."
	go test -coverprofile=coverage.out ./...

coverage: test
	@echo "==> Checking coverage..."
	go tool cover -func=coverage.out

run: build
	@echo "==> Running application..."
	./$(BINARY_NAME)-cli
