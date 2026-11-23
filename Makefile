BIN        := gss
PKG        := ./src/cmd/server
VERSION    ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)

# Packagers override these; DESTDIR is the staging root.
DESTDIR    ?=
PREFIX     ?= /usr
BINDIR     ?= $(PREFIX)/bin
SYSCONFDIR ?= /etc
SCRIPTDIR  ?= $(SYSCONFDIR)/gss/scripts

GOFLAGS    ?= -trimpath -buildvcs=false
LDFLAGS    ?= -s -w -X main.version=$(VERSION)

.PHONY: run install

run: ## Build and run with dev scripts: make run ARGS="version"
	go build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(BIN) $(PKG)
	GSS_SCRIPTS_DIR=$(CURDIR)/scripts ./$(BIN) $(ARGS)

install: ## Build and install binary + scripts (honours DESTDIR, PREFIX)
	go build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(BIN) $(PKG)
	install -Dm755 $(BIN) $(DESTDIR)$(BINDIR)/$(BIN)
	install -d $(DESTDIR)$(SCRIPTDIR)
	install -m755 scripts/*.sh $(DESTDIR)$(SCRIPTDIR)/
