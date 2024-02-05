#!/usr/bin/make -f

CWD := $(dir $(abspath $(firstword $(MAKEFILE_LIST))))

# ldflags = -X main.MakeFileInstallDirectory=$(CWD)
# ldflags := $(strip $(ldflags))
# BUILD_FLAGS := -ldflags '$(ldflags)'

.PHONY: build
build:
	go build $(BUILD_FLAGS) -o ./bin/spawn ./cmd/spawn

.PHONY: run
run:
	go run ./cmd/spawn $(filter-out $@,$(MAKECMDGOALS))

.PHONY: install
install:
	go install $(BUILD_FLAGS) ./cmd/spawn


###############################################################################
###                                  heighliner                             ###
###############################################################################

get-heighliner:
	git clone https://github.com/strangelove-ventures/heighliner.git
	cd heighliner && go install

.PHONY: get-heighliner