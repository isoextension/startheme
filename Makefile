# Makefile — Go project helper
app ?= startheme
_PKG ?= ./...
MAKE ?= make
GO ?= go
CGO_ENABLED ?= 0
out_dir ?= bin
out = $(OUTDIR)/$(APP)

# Default build tags (empty by default)
BUILD_TAGS ?=

LDFLAGS ?= -s -w
_TIME := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
_GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)
_GOVERSION := $(shell $(GO) version | awk '{print $$3}')

# ldflags injection for version info
LDFLAGS_VERSION := -X 'main.buildTime=$(TIME)' -X 'main.commit=$(GIT_COMMIT)' -X 'main.goVersion=$(GOVERSION)'

.PHONY: all build run install test fmt vet tidy lint clean deps cross dev head

help:
	@echo "startheme Make"
	@echo "  all      Do everything and install Startheme."
	@echo "  help     Display this message."
	@echo "  build    Build Startheme from source."
	@echo "  run      Run Startheme."
	@echo "  install  Install Startheme."
	@echo "  test     Test Startheme."
	@echo "  fmt      Format source code."
	@echo "  vet      Audit code."
	@echo "  tidy     Update/make go.sum."
	@echo "  deps     Install dependencies."
	@echo "  cross    Cross compile helper. (the code still doesnt work for windows or macos, contribs are welcome)"
	@echo "  dev      Quick development helper."
	@echo "  head     Build from HEAD."

all: deps tidy build install clean

# Build for current platform
build:
	@mkdir -p $(OUTDIR)
	@CGO_ENABLED=$(CGO_ENABLED) $(GO) build -trimpath -tags="$(BUILD_TAGS)" -ldflags "$(LDFLAGS) $(LDFLAGS_VERSION)" -o $(OUTDIR)/$(APP) ./src/main.go

# Run locally (uses the package main in ./bin/startheme)
run:
	$(GO) run -tags="$(BUILD_TAGS)" ./bin/$(APP)

# Install 
install:
	sudo rm -f /usr/local/bin/startheme
	sudo mv ./bin/startheme /usr/local/bin/startheme

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
_binary_name = $(app)-$(os)-$(arch)
cross:
ifndef os
	$(error os is not set. e.g. make cross os=linux arch=amd64)
endif
ifndef arch
	$(error arch is not set. e.g. make cross os=linux arch=amd64)
endif
ifeq 
	@mkdir -p dist
	GOOS=$(os) GOARCH=$(arch) CGO_ENABLED=$(CGO_ENABLED) $(GO) build -trimpath -tags="$(BUILD_TAGS)" -ldflags "$(LDFLAGS) $(LDFLAGS_VERSION)" -o dist/$(app)-$(os)-$(arch) ./src/main.go

# Quick dev loop: format, vet, test, build
dev: fmt vet test build

# Clean artifacts
clean:
	@rm -rf $(out_dir) dist

head:
	@echo "Building from $(_GIT_COMMIT)"
	@git pull origin trunk
	@git push origin trunk
	@$(MAKE) cross os=$(os) arch=$(arch)



