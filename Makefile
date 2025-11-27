# pgEdge Anonymizer Makefile

BINARY_NAME=pgedge-anonymizer
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "1.0.0-alpha1")
BUILD_TIME=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS=-ldflags "-X github.com/pgedge/pgedge-anonymizer/internal/version.Version=$(VERSION) -X github.com/pgedge/pgedge-anonymizer/internal/version.BuildTime=$(BUILD_TIME)"

.PHONY: all build test lint clean fmt vet install

all: fmt vet lint test build

build:
	go build $(LDFLAGS) -o bin/$(BINARY_NAME) ./cmd/pgedge-anonymizer

install:
	go install $(LDFLAGS) ./cmd/pgedge-anonymizer

test:
	go test -v -race -cover ./...

lint:
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not installed, skipping..."; \
	fi

fmt:
	gofmt -s -w .

vet:
	go vet ./...

clean:
	rm -rf bin/
	go clean

# Run tests with coverage report
coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Download dependencies
deps:
	go mod download
	go mod tidy

# Build for multiple platforms
build-all: build-linux build-darwin build-windows

build-linux:
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-linux-amd64 ./cmd/pgedge-anonymizer

build-darwin:
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-darwin-amd64 ./cmd/pgedge-anonymizer
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-darwin-arm64 ./cmd/pgedge-anonymizer

build-windows:
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-windows-amd64.exe ./cmd/pgedge-anonymizer
