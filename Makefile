.PHONY: build test install clean fmt vet lint setup-hooks

BINARY_NAME=plane-cli
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS=-ldflags "-X github.com/rohithmahesh3/plane-cli/cmd.version=$(VERSION) -X github.com/rohithmahesh3/plane-cli/cmd.commit=$(COMMIT) -X github.com/rohithmahesh3/plane-cli/cmd.date=$(DATE)"

build:
	go build $(LDFLAGS) -o $(BINARY_NAME) main.go

test:
	go test -v ./...

install: build
	cp $(BINARY_NAME) $(GOPATH)/bin/

clean:
	rm -f $(BINARY_NAME)
	go clean

fmt:
	go fmt ./...

vet:
	go vet ./...

lint:
	golangci-lint run

deps:
	go mod download
	go mod tidy

run: build
	./$(BINARY_NAME)

# Cross compilation
build-all:
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-amd64 main.go
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-arm64 main.go
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-amd64 main.go
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-arm64 main.go
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-windows-amd64.exe main.go

# Pre-commit hooks setup
setup-hooks:
	@echo "Setting up pre-commit hooks..."
	@command -v pre-commit >/dev/null 2>&1 || { echo "pre-commit not found. Installing..."; pip install pre-commit; }
	pre-commit install
	@echo "Pre-commit hooks installed successfully!"

# Run all checks (fmt, vet, lint, test)
check: fmt vet lint test
