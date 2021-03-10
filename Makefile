SHELL=/usr/bin/env bash

GOVERSION:=$(shell go version | cut -d' ' -f 3 | cut -d. -f 2)
ifeq ($(shell expr $(GOVERSION) \< 13), 1)
$(warning Your Golang version is go 1.$(GOVERSION))
$(error Update Golang to version $(shell grep '^go' go.mod))
endif

CLEAN:=
BINS:=./venus-wallet

git=$(subst -,.,$(shell git describe --always --match=NeVeRmAtCh --dirty 2>/dev/null || git rev-parse --short HEAD 2>/dev/null))

ldflags=-X=github.com/ipfs-force-community/venus-wallet/version.CurrentCommit='+git$(git)'

ifneq ($(strip $(LDFLAGS)),)
	ldflags+=-extldflags=$(LDFLAGS)
endif

GOFLAGS+=-ldflags="$(ldflags)"

wallet: show-env $(BUILD_DEPS)
	rm -f venus-wallet
	go build $(GOFLAGS) -o venus-wallet ./cmd/wallet/main.go
	./venus-wallet --version


show-env:
	@echo '_________________build_environment_______________'
	@echo '| CC=$(CC)'
	@echo '| CGO_CFLAGS=$(CGO_CFLAGS)'
	@echo '| git commit=$(git)'
	@echo '-------------------------------------------------'

lint:
	gofmt -s -w ./
	golangci-lint run

# MISC

clean:
	rm -rf $(CLEAN) $(BINS)
.PHONY: clean

print-%:
	@echo $*=$($*)
