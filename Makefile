BINDIR ?= ./bin
PREFIX ?= $(HOME)/.local/bin

.PHONY: build test install uninstall

build:
	mkdir -p $(BINDIR)
	go build -o $(BINDIR)/ailc ./compiler

test:
	go test ./...

install: build
	mkdir -p $(PREFIX)
	cp $(BINDIR)/ailc $(PREFIX)/ailc
	@echo "installed ailc to $(PREFIX)/ailc"

uninstall:
	rm -f $(PREFIX)/ailc
