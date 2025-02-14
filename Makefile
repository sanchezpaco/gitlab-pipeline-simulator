BINARY_NAME = glps 
GO_FILES = $(shell find . -type f -name '*.go')
VERSION=$(shell cat .semver | tr -d '[:space:]')

.PHONY: all build test lint clean run

all: build

build: $(GO_FILES)
	@echo "Building $(BINARY_NAME) v$(VERSION)"
	@go build -o $(BINARY_NAME) -ldflags "-X main.version=$(VERSION)"

test:
	@echo "Running tests..."
	@go test -v -coverprofile=coverage.out ./...


.PHONY: tag
tag:
	@if [ ! -f .semver ]; then \
		echo "Error: .semver file not found!" >&2; \
		exit 1; \
	fi
	git tag -a "v$(VERSION)" -m "Release v$(VERSION)"
	git push origin "v$(VERSION)"