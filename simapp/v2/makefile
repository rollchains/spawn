PACKAGES_NOSIMULATION=$(shell go list ./... | grep -v '/simulation')
PACKAGES_SIMTEST=$(shell go list ./... | grep '/simulation')
export VERSION := $(shell echo $(shell git describe --tags --always --match "v*") | sed 's/^v//')
export COMMIT := $(shell git log -1 --format='%H')
LEDGER_ENABLED ?= true
BINDIR ?= $(GOPATH)/bin
BUILDDIR ?= $(CURDIR)/build
SIMAPP = simapp
MOCKS_DIR = $(CURDIR)/tests/mocks
HTTPS_GIT := https://github.com/cosmos/cosmos-sdk.git
DOCKER := $(shell which docker)
PROJECT_NAME = $(shell git remote get-url origin | xargs basename -s .git)

.PHONY: install
install:
	go install ./simdv2/simdv2.go


.PHONY: init
init: install
	rm -rf /home/reece/.simappv2
	simdv2 init moniker
	echo -n "easy visit salt program one absurd strike warm remind badge cake muscle survey monitor swamp nation holiday escape vault tuna claw essay deliver nominee" > $HOME/.simappv2/mnemonic.txt

	# TODO: no override command to skip interactive prompt?
	simdv2 keys add val --recover --source $HOME/.simappv2/mnemonic.txt --keyring-backend=test

	simdv2 keys list --keyring-backend=test

	simdv2 genesis add-genesis-account val 1000000stake --keyring-backend=test
	simdv2 genesis gentx val 1000000stake --keyring-backend=test --chain-id=gcosmos
	simdv2 genesis collect-gentxs

	simdv2 start
