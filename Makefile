#!/usr/bin/make -f

CWD := $(dir $(abspath $(firstword $(MAKEFILE_LIST))))

# don't override user values
ifeq (,$(VERSION))
  VERSION := $(shell git describe --tags)
  # if VERSION is empty, then populate it with branch's name and raw commit hash
  ifeq (,$(VERSION))
    VERSION := $(BRANCH)-$(COMMIT)
  endif
endif

DATE := $(shell date '+%Y-%m-%dT%H:%M:%S')
HEAD = $(shell git rev-parse HEAD)
LD_FLAGS = -X main.SpawnVersion=$(VERSION)
BUILD_FLAGS = -mod=readonly -ldflags='$(LD_FLAGS)'



## install: Install the binary.
install:
	@echo ⏳ Installing Spawn...
	go install $(BUILD_FLAGS) ./cmd/spawn
	@echo ✅ Spawn installed to $(shell which spawn)

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
	git clone https://github.com/strangelove-ventures/heighliner.git
	cd heighliner && go install
	@echo ✅ heighliner installed to $(shell which heighliner)

## get-localic: Install local interchain testnet manager.
get-localic:
	@echo ⏳ Installing local-interchain...
	git clone --branch v8.1.0 https://github.com/strangelove-ventures/interchaintest.git
	cd interchaintest/local-interchain && make install
	@echo ✅ local-interchain installed to $(shell which local-ic)

.PHONY: get-heighliner

help: Makefile
	@echo
	@echo " Choose a command run in "spawn", or just run 'make' for install"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo

.PHONY: help

# ---- Developer Templates ----
template-default: install
	spawn new myproject --debug --bech32=cosmos --bin=appd --disable=tokenfactory

template-specific: install
	spawn new myproject --no-git --bin=rolld --bech32=roll --denom=uroll --disable=globalfee,poa


.DEFAULT_GOAL := install