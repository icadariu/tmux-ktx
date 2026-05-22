CMD_PKG ?= .
BIN     ?= tmux-ktx
GOFLAGS ?=
LDFLAGS ?= -s -w

.DEFAULT_GOAL := help
.PHONY: help build install test fmt vet clean

help: ## Show this help
	@awk 'BEGIN {FS = ":.*?## "; printf "Usage: make \033[36m<target>\033[0m\n\nTargets:\n"} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the binary into ./$(BIN)
	go build $(GOFLAGS) -ldflags '$(LDFLAGS)' -o $(BIN) $(CMD_PKG)

install: ## Install the binary into $$GOPATH/bin
	go install $(GOFLAGS) -ldflags '$(LDFLAGS)' $(CMD_PKG)

test: ## Run all Go tests
	go test ./...

fmt: ## Format Go source files with gofmt -s
	gofmt -s -w .

vet: ## Run go vet on all packages
	go vet ./...

clean: ## Remove built binary
	rm -f $(BIN)
