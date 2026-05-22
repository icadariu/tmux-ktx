CMD_PKG ?= .
BIN     ?= tmux-ktx
GOFLAGS ?=
LDFLAGS ?= -s -w

.PHONY: build install test fmt vet clean

build:
	go build $(GOFLAGS) -ldflags '$(LDFLAGS)' -o $(BIN) $(CMD_PKG)

install:
	go install $(GOFLAGS) -ldflags '$(LDFLAGS)' $(CMD_PKG)

test:
	go test ./...

fmt:
	gofmt -s -w .

vet:
	go vet ./...

clean:
	rm -f $(BIN)
