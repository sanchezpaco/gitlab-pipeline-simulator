BINARY_NAME = glps 
GO_FILES = $(shell find . -type f -name '*.go')
VERSION = $(shell git describe --tags --always --dirty)

.PHONY: all build test lint clean run

all: build

build: $(GO_FILES)
	@echo "Building $(BINARY_NAME) v$(VERSION)"
	@go build -o $(BINARY_NAME) -ldflags "-X main.version=$(VERSION)"

test:
	@echo "Running tests..."
	@go test -v -coverprofile=coverage.out ./...