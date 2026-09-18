BINDIR ?= ./bin
PREFIX ?= $(HOME)/.local/bin
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS = -X main.version=$(VERSION)

.PHONY: build test install uninstall

build:
	mkdir -p $(BINDIR)
	go build -ldflags "$(LDFLAGS)" -o $(BINDIR)/canlc ./compiler

test:
	go test ./...

install: build
	mkdir -p $(PREFIX)
	cp $(BINDIR)/canlc $(PREFIX)/canlc
	@echo "installed canlc to $(PREFIX)/canlc"

uninstall:
	rm -f $(PREFIX)/canlc
