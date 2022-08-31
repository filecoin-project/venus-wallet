export CGO_CFLAGS_ALLOW=-D__BLST_PORTABLE__
export CGO_CFLAGS=-D__BLST_PORTABLE__

SHELL=/usr/bin/env bash

GOVERSION:=$(shell go version | cut -d' ' -f 3 | cut -d. -f 2)
ifeq ($(shell expr $(GOVERSION) \< 13), 1)
$(warning Your Golang version is go 1.$(GOVERSION))
$(error Update Golang to version $(shell grep '^go' go.mod))
endif

CLEAN:=
BINS:=./venus-wallet

git=$(subst -,.,$(shell git describe --always --match=NeVeRmAtCh --dirty 2>/dev/null || git rev-parse --short HEAD 2>/dev/null))

ldflags=-X=github.com/filecoin-project/venus-wallet/version.CurrentCommit='+git$(git)'
ifneq ($(strip $(LDFLAGS)),)
	ldflags+=-extldflags=$(LDFLAGS)
endif

GOFLAGS+=-ldflags="$(ldflags)"

wallet: show-env $(BUILD_DEPS)
	rm -f venus-wallet
	go build $(GOFLAGS) -o venus-wallet ./cmd/wallet/main.go
	./venus-wallet --version


linux: 	show-env $(BUILD_DEPS)
	rm -f venus-wallet
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 CC=x86_64-linux-musl-gcc CGO_LDFLAGS="-static" go build $(GOFLAGS) -o venus-wallet ./cmd/wallet/main.go

show-env:
	@echo '_________________build_environment_______________'
	@echo '| CC=$(CC)'
	@echo '| CGO_CFLAGS=$(CGO_CFLAGS)'
	@echo '| git commit=$(git)'
	@echo '-------------------------------------------------'

lint:
	gofmt -s -w ./
	golangci-lint run

clean:
	rm -rf $(CLEAN) $(BINS)
.PHONY: clean

print-%:
	@echo $*=$($*)



.PHONY: docker


docker:
	docker build --build-arg https_proxy=$(BUILD_DOCKER_PROXY) -t venus-wallet .



TAG:=test
docker:
	curl -O https://raw.githubusercontent.com/filecoin-project/venus-docs/master/script/dockerfile
	docker build --build-arg https_proxy=$(BUILD_DOCKER_PROXY) --build-arg BUILD_TARGET=venus-wallet -t venus-wallet .
	docker tag venus-auth filvenus/venus-auth:$(TAG)
