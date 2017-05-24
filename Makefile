PKG_SHORT = swift4hubic
PKG       = github.com/reimannf/${PKG_SHORT}
PREFIX   ?= /usr/local

GO            := GOPATH=$(CURDIR)/.gopath GOBIN=$(CURDIR)/build go
GO_BUILDFLAGS :=
GO_LDFLAGS    := -s -w

GOFMT          = gofmt
GOLINT         = golint
GLIDE          = glide

GO_ALLPKGS := $(shell go list $(PKG)/... | grep -v vendor)

M = $(shell printf "\033[34;1m▶\033[0m")

all: build/${PKG_SHORT}

build/swift4hubic: ALWAYS; $(info $(M) building…)
	$(GO) install $(GO_BUILDFLAGS) -ldflags '$(GO_LDFLAGS)' '$(PKG)'

check: all ALWAYS gofmt golint govet

gofmt: ALWAYS; $(info $(M) gofmt…)
	ret=0 && for d in $$($(GO) list -f '{{.Dir}}' ./... | grep -v /vendor/); do \
		$(GOFMT) -l -w $$d/*.go || ret=$$? ; \
	 done ; exit $$ret

golint: ALWAYS; $(info $(M) golint…)
	ret=0 && for pkg in $(GO_ALLPKGS); do \
		$(GOLINT) $$pkg || ret=$$? ; \
	 done ; exit $$ret

govet: ALWAYS; $(info $(M) go vet…)
	$(GO) vet . $(GO_ALLPKGS)

install: ALWAYS all
	install -D -m 0755 build/${PKG_SHORT} "$(DESTDIR)$(PREFIX)/bin/${PKG_SHORT}"

vendor: ALWAYS glide.yaml ; $(info $(M) updating dependencies…)
	$(GLIDE) update

clean: ALWAYS ; $(info $(M) cleaning…)
	@rm -rf build

.PHONY: ALWAYS
