BINDIR ?= ./bin
PREFIX ?= $(HOME)/.local/bin
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS = -X main.version=$(VERSION)

.PHONY: build test install uninstall

build:
	mkdir -p $(BINDIR)
	go build -ldflags "$(LDFLAGS)" -o $(BINDIR)/ailc ./compiler

test:
	go test ./...

install: build
	mkdir -p $(PREFIX)
	cp $(BINDIR)/ailc $(PREFIX)/ailc
	@echo "installed ailc to $(PREFIX)/ailc"

uninstall:
	rm -f $(PREFIX)/ailc
