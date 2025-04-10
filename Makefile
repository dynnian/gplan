# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=gplan
PACKAGE=codeberg.org/dynnian/gplan

# Build directory
BUILD_DIR=build

# Git information
GIT_COMMIT=$(shell git rev-parse HEAD)
BUILD_DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Linker flags
LDFLAGS=-ldflags "-X '$(PACKAGE)/internal/version.Get().Version=$(VERSION)' -X '$(PACKAGE)/internal/version.Get().BuildDate=$(BUILD_DATE)' -X '$(PACKAGE)/internal/version.Get().GitCommit=$(GIT_COMMIT)'"

# Platforms to build for
PLATFORMS=windows/amd64 windows/arm64 darwin/amd64 darwin/arm64 linux/amd64 linux/arm64

# GOPATH
GOPATH := $(shell go env GOPATH)

# Tools
GOLANGCI_LINT := $(shell command -v golangci-lint 2> /dev/null)
GOFUMPT := $(shell command -v gofumpt 2> /dev/null)
GOIMPORTS := $(shell command -v goimports 2> /dev/null)
GOLINES := $(shell command -v golines 2> /dev/null)

.PHONY: all build clean deps lint format build-all version install uninstall help

all: deps lint format build

build:
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) .

clean:
	$(GOCLEAN)
	@rm -rf $(BUILD_DIR)

deps:
	@echo "Checking and updating dependencies..."
	@go mod tidy
	@if [ -z "$$(git status --porcelain go.mod go.sum)" ]; then \
		echo "No missing dependencies. All modules are up to date."; \
	else \
		echo "Dependencies updated. Please review changes in go.mod and go.sum."; \
	fi
ifndef GOLANGCI_LINT
	@echo "Installing golangci-lint..."
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(GOPATH)/bin
endif
ifndef GOFUMPT
	@echo "Installing gofumpt..."
	@go install mvdan.cc/gofumpt@latest
endif
ifndef GOIMPORTS
	@echo "Installing goimports..."
	@go install golang.org/x/tools/cmd/goimports@latest
endif
ifndef GOLINES
	@echo "Installing golines..."
	@go install github.com/segmentio/golines@latest
endif

lint:
	@echo "Running linter..."
	@golangci-lint run --fix -c .golangci.yml ./...


format:
	@echo "Formatting code..."
	@gofumpt -l -w .
	@golines -l -m 120 -t 4 -w .
	@golines -w .
	echo "Code formatted."; \

build-all:
	mkdir -p $(BUILD_DIR)
	$(foreach PLATFORM,$(PLATFORMS),\
		$(eval GOOS=$(word 1,$(subst /, ,$(PLATFORM))))\
		$(eval GOARCH=$(word 2,$(subst /, ,$(PLATFORM))))\
		$(eval EXTENSION=$(if $(filter $(GOOS),windows),.exe,))\
		GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=0 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-$(GOOS)-$(GOARCH)$(EXTENSION) .;\
	)

# Version information
version:
	@echo "Git Commit: $(GIT_COMMIT)"
	@echo "Build Date: $(BUILD_DATE)"

# Install the application
install:
	$(GOBUILD) -o $(GOPATH)/bin/$(BINARY_NAME) .

# Uninstall the application
uninstall:
	@rm $(GOPATH)/bin/$(BINARY_NAME)

# Installation help
help:
	@echo "Available commands:"
	@echo "  make              - Run deps, lint, format, test, and build"
	@echo "  make build        - Build for the current platform"
	@echo "  make clean        - Remove build artifacts"
	@echo "  make deps         - Download dependencies and install tools"
	@echo "  make lint         - Run golangci-lint for linting"
	@echo "  make format       - Format code using gofumpt, goimports, and golines"
	@echo "  make build-all    - Build for all specified platforms"
	@echo "  make version      - Display the current git commit and build date"
	@echo "  make install      - Install the application to GOPATH/bin"
	@echo "  make help         - Display this help information"
