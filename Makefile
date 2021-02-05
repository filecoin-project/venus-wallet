SHELL=/usr/bin/env bash

GOVERSION:=$(shell go version | cut -d' ' -f 3 | cut -d. -f 2)
ifeq ($(shell expr $(GOVERSION) \< 13), 1)
$(warning Your Golang version is go 1.$(GOVERSION))
$(error Update Golang to version $(shell grep '^go' go.mod))
endif

# git modules that need to be loaded
MODULES:=

CLEAN:=
BINS:=

git=$(subst -,.,$(shell git describe --always --match=NeVeRmAtCh --dirty 2>/dev/null || git rev-parse --short HEAD 2>/dev/null))

ldflags=-X=github.com/ipfs-force-community/venus-wallet/build.CurrentCommit='+git$(git)'

ifneq ($(strip $(LDFLAGS)),)
	ldflags+=-extldflags=$(LDFLAGS)
endif

GOFLAGS+=-ldflags="$(ldflags)"

## MAIN BINARIES

CLEAN+=build/.update-modules

wallet: show-env $(BUILD_DEPS)
	rm -f wallet
	go build $(GOFLAGS) -o venus-wallet ./cmd/wallet
	./venus-wallet --version


show-env:
	@echo '_________________build_environment_______________'
	@echo '| CC=$(CC)'
	@echo '| CGO_CFLAGS=$(CGO_CFLAGS)'
	@echo '| git commit=$(git)'
	@echo '-------------------------------------------------'


.PHONY: wallet
BINS+=wallet

# MISC

clean:
	rm -rf $(CLEAN) $(BINS)
.PHONY: clean

dist-clean:
	git clean -xdff
	git submodule deinit --all -f
.PHONY: dist-clean

print-%:
	@echo $*=$($*)
