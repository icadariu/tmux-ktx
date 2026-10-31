MAIN    ?= .
BIN     ?= tmux-ktx
GOFLAGS ?=

BRANCH    := $(shell git rev-parse --abbrev-ref HEAD 2>/dev/null || echo dev)
ifeq ($(BRANCH),main)
VERSION   := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
else
VERSION   := $(shell git describe --tags --always 2>/dev/null || echo dev)-dev
endif
COMMIT    := $(shell git rev-parse --short HEAD 2>/dev/null || echo none)
BUILDTIME := $(shell LC_ALL=C date -u +%d-%b-%Y_%H:%M)

LDFLAGS ?= -ldflags '-s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.buildTime=$(BUILDTIME)'

.DEFAULT_GOAL := help
.PHONY: help build install test fmt vet clean version demo release

help: ## Show this help
	@awk 'BEGIN {FS = ":.*?## "; printf "Usage: make \033[36m<target>\033[0m\n\nTargets:\n"} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the binary into ./$(BIN)
	go build $(GOFLAGS) $(LDFLAGS) -o $(BIN) $(MAIN)

install: ## Install latest into $$GOPATH/bin
	go install $(GOFLAGS) $(LDFLAGS) $(MAIN)

version: ## Print the version that would be embedded
	@echo "$(BIN) $(VERSION) (built $(BUILDTIME), commit $(COMMIT))"

test: ## Run all Go tests
	go test ./...

fmt: ## Format Go source files with gofmt -s
	gofmt -s -w .

vet: ## Run go vet on all packages
	go vet ./...

demo: build ## Record docs/demo.cast + docs/demo.gif (needs asciinema >= 3, agg, zsh, kubectl)
	./docs/demo.sh docs/demo.cast
	agg -q --theme monokai --idle-time-limit 2 --last-frame-duration 4 docs/demo.cast docs/demo.gif

release: ## Tag and push a signed release: make release TAG=vX.Y.Z
	@echo "$(TAG)" | grep -Eq '^v[0-9]+\.[0-9]+\.[0-9]+(-[0-9A-Za-z.-]+)?$$' || { echo "usage: make release TAG=vX.Y.Z"; exit 1; }
	@test "$$(git rev-parse --abbrev-ref HEAD)" = main || { echo "release from main only"; exit 1; }
	@test -z "$$(git status --porcelain)" || { echo "working tree is not clean"; exit 1; }
	@git fetch -q origin main --tags
	@test "$$(git rev-parse HEAD)" = "$$(git rev-parse origin/main)" || { echo "main is not in sync with origin/main"; exit 1; }
	@! git rev-parse -q --verify "refs/tags/$(TAG)" >/dev/null || { echo "tag $(TAG) already exists"; exit 1; }
	go test -count=1 ./...
	git tag -s $(TAG) -m "$(TAG)"
	git push origin $(TAG)
	@echo "pushed $(TAG); follow the build with: gh run watch \$$(gh run list -w release -L1 --json databaseId -q '.[0].databaseId')"

clean: ## Remove built binary
	rm -f $(BIN)
