MAIN    ?= .
BIN     ?= tmux-ktx

.DEFAULT_GOAL := help

include .make/version.mk .make/release.mk

.PHONY: help test fmt vet clean demo

help: ## Show this help
	@awk 'BEGIN {FS = ":.*?## "; printf "Usage: make \033[36m<target>\033[0m\n\nTargets:\n"} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

test: ## Run all Go tests
	go test ./...

fmt: ## Format Go source files with gofmt -s
	gofmt -s -w .

vet: ## Run go vet on all packages
	go vet ./...

demo: build ## Record docs/demo.cast + docs/demo.gif (needs asciinema >= 3, agg, zsh, kubectl)
	./docs/demo.sh docs/demo.cast
	agg -q --theme monokai --idle-time-limit 2 --last-frame-duration 4 docs/demo.cast docs/demo.gif

clean: ## Remove built binary
	rm -f $(BIN)
