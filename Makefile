BINARY     := agency-cli
MODULE     := $(shell go list -m)
VERSION    := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS    := -ldflags "-X $(MODULE)/cmd.Version=$(VERSION)"

BUILD_DIR  := ./bin

GOBIN                 := $(shell go env GOBIN)
GOLANGCI_LINT_VERSION := v2.11.3
LEFTHOOK_VERSION      := v1.13.6

export PATH := $(GOBIN):$(PATH)

.PHONY: all build run install clean test lint fmt vet tidy check tools hooks setup help

all: check build ## Run checks and build

# ── Setup ─────────────────────────────────────────────────────────────────────

setup: tools hooks ## Install tools and git hooks

tools: ## Install required development tools
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)
	go install github.com/evilmartians/lefthook@$(LEFTHOOK_VERSION)

hooks: ## Install lefthook git hooks
	lefthook install

# ── Build ─────────────────────────────────────────────────────────────────────

build: ## Build the binary into ./bin/
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY) .

run: ## Run the application (usage: make run ARGS="list --category engineering")
	go run . $(ARGS)

install: ## Install the binary to $GOPATH/bin
	go install $(LDFLAGS) .

clean: ## Remove build artifacts
	@rm -rf $(BUILD_DIR)

# ── Quality ───────────────────────────────────────────────────────────────────

test: ## Run all tests
	go test ./...

test-verbose: ## Run all tests with verbose output
	go test -v ./...

lint: ## Run golangci-lint
	golangci-lint run ./...

fmt: ## Format all Go source files
	gofmt -w .

vet: ## Run go vet
	go vet ./...

fix: ## Run go fix
	go fix ./...

tidy: ## Tidy and verify go modules
	go mod tidy
	go mod verify

check: vet lint ## Run vet and lint

# ── Help ──────────────────────────────────────────────────────────────────────

help: ## Show this help message
	@grep -E '^[a-zA-Z_-]+:.*##' $(MAKEFILE_LIST) \
		| awk 'BEGIN {FS = ":.*##"}; {printf "  \033[36m%-16s\033[0m %s\n", $$1, $$2}'
