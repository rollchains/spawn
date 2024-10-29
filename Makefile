#!/usr/bin/make -f

CWD := $(dir $(abspath $(firstword $(MAKEFILE_LIST))))

# don't override user values
ifeq (,$(VERSION))
  VERSION := $(shell git tag --sort=version:refname | tail -1)
  # if VERSION is empty, then populate it with branch's name and raw commit hash
  ifeq (,$(VERSION))
    VERSION := $(BRANCH)-$(COMMIT)
  endif
endif

DATE := $(shell date '+%Y-%m-%dT%H:%M:%S')
HEAD = $(shell git rev-parse HEAD)
# MacOS sequoia + XCode 16 causes the binary to instantly die
# Stripping the debug info with -s -w fixes this so the binary
# is properly signed
# ref: https://github.com/rollchains/spawn/issues/248
LD_FLAGS = -X main.SpawnVersion=$(VERSION) -s -w
BUILD_FLAGS = -mod=readonly -ldflags='$(LD_FLAGS)'

## mod-tidy: go mod tidy spawn, simapp, and interchaintest with proper go.mod suffixes
mod-tidy:
	go mod tidy
	mv simapp/interchaintest/go.mod_ simapp/interchaintest/go.mod
	cd simapp && go mod tidy
	cd simapp/interchaintest && go mod tidy
	mv simapp/interchaintest/go.mod simapp/interchaintest/go.mod_

## install: Install the binary.
install:
	@echo ⏳ Installing Spawn...
	go install $(BUILD_FLAGS) ./cmd/spawn
	@echo ✅ Spawn installed

## build: Build to ./bin/spawn.
build:
	go build $(BUILD_FLAGS) -o ./bin/spawn ./cmd/spawn

## run: Run the raw source.
run:
	go run ./cmd/spawn $(filter-out $@,$(MAKECMDGOALS))

.PHONY: install build run

# ---- Downloads ----

## get-heighliner: Install the cosmos docker utility.
get-heighliner:
	@echo ⏳ Installing heighliner...
	git clone --depth 1 https://github.com/strangelove-ventures/heighliner.git
	cd heighliner && go install
	@sleep 0.1
	@echo ✅ heighliner installed to $(shell which heighliner)

get-localic:
	@echo "Installing local-interchain"
	git clone --depth 1 --branch v8.7.1 https://github.com/strangelove-ventures/interchaintest.git interchaintest-downloader
	cd interchaintest-downloader/local-interchain && make install
	@sleep 0.1
	@echo ✅ local-interchain installed $(shell which local-ic)

.PHONY: get-heighliner

.PHONY: build-docs
## build-docs: Build the documentation.
build-docs:
	@cd docs && npm ci && npm run build

## serve-docs: Run docs server
serve-docs:
	@cd docs && npm run serve

## dev-docs: Run dev docs instance
dev-docs:
	@cd docs && npm run dev

help: Makefile
	@echo
	@echo " Choose a command run in "spawn", or just run 'make' for install"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo

.PHONY: help

# ---- Developer Templates ----
template-staking: install
	spawn new myproject --consensus=proof-of-stake --debug --bech32=cosmos --bin=appd --bypass-prompt --log-level=debug --org=reece

template-poa: install
	spawn new myproject --consensus=proof-of-authority --debug --no-git --bin=rolld --bech32=roll --denom=uroll --bypass-prompt --log-level=debug

template-ics: install
	spawn new myproject --consensus=interchain-security --debug --no-git --bin=rolld --bech32=roll --denom=uroll --bypass-prompt --log-level=debug


.DEFAULT_GOAL := install
