# Makefile — Go project helper
APP ?= startheme
PKG ?= ./...
GO ?= go
CGO_ENABLED ?= 0
OUTDIR ?= bin

# Default build tags (empty by default)
BUILD_TAGS ?=

LDFLAGS ?= -s -w
TIME := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)
GOVERSION := $(shell $(GO) version | awk '{print $$3}')

# ldflags injection for version info
LDFLAGS_VERSION := -X 'main.buildTime=$(TIME)' -X 'main.commit=$(GIT_COMMIT)' -X 'main.goVersion=$(GOVERSION)'

.PHONY: all build run install test fmt vet tidy lint clean deps cross

all: build

# Build for current platform
build:
	@mkdir -p $(OUTDIR)
	CGO_ENABLED=$(CGO_ENABLED) $(GO) build -trimpath -tags="$(BUILD_TAGS)" -ldflags "$(LDFLAGS) $(LDFLAGS_VERSION)" -o $(OUTDIR)/$(APP) ./src/main.go

# Run locally (uses the package main in ./cmd/$(APP))
run:
	$(GO) run -tags="$(BUILD_TAGS)" ./cmd/$(APP)

# Install to $GOBIN or GOPATH/bin
install:
	CGO_ENABLED=$(CGO_ENABLED) $(GO) install -tags="$(BUILD_TAGS)" -ldflags "$(LDFLAGS) $(LDFLAGS_VERSION)" ./bin/$(APP)

# Test everything
test:
	$(GO) test ./... -v

# Format code
fmt:
	$(GO) fmt ./...

# Vet code
vet:
	$(GO) vet ./...

# Tidy modules
tidy:
	$(GO) mod tidy

# Download modules (useful for CI)
deps:
	$(GO) mod download

# Lint (if you have golangci-lint installed)
lint:
	@golangci-lint run || true

# Cross-compile helper
# Example: make cross GOOS=linux GOARCH=arm64
cross:
ifndef GOOS
	$(error GOOS is not set. e.g. make cross GOOS=linux GOARCH=amd64)
endif
ifndef GOARCH
	$(error GOARCH is not set. e.g. make cross GOOS=linux GOARCH=amd64)
endif
	@mkdir -p dist
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=0 $(GO) build -trimpath -tags="$(BUILD_TAGS)" -ldflags "$(LDFLAGS) $(LDFLAGS_VERSION)" -o dist/$(APP)-$(GOOS)-$(GOARCH) ./cmd/$(APP)

# Quick dev loop: format, vet, test, build
dev: fmt vet test build

# Clean artifacts
clean:
	@rm -rf $(OUTDIR) dist

# Show make help
help:
	@printf "Makefile targets:\n"
	@printf "  make build       Build binary (./bin/$(APP))\n"
	@printf "  make run         Run from source (go run)\n"
	@printf "  make install     Install to GOBIN\n"
	@printf "  make test        Run tests\n"
	@printf "  make fmt         go fmt ./...\n"
	@printf "  make vet         go vet ./...\n"
	@printf "  make tidy        go mod tidy\n"
	@printf "  make deps        go mod download\n"
	@printf "  make lint        Run golangci-lint (if installed)\n"
	@printf "  make cross GOOS=<os> GOARCH=<arch>  Cross-build (outputs in dist/)\n"
	@printf "  make dev         fmt, vet, test, build\n"
	@printf "  make clean       Remove build artifacts\n"

