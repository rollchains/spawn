---
title: "CW Validator Reviews"
sidebar_label: "CW Validator Reviews"
slug: /demo/cw-validator-reviews
---

# CosmWasm Validator Reviews

- New network with cosmwasm
- Write a contract which takes in all validators
- Write a module with an endblock that updates all validators every X blocks
- Prove this works and is set, set reviews


## Prerequisites
- [System Setup](../01-setup/01-system-setup.md)
- [Install Spawn](../01-setup/02-install-spawn.md)

## CosmWasm Setup

```bash
# install rust & things
cargo install cargo-generate --features vendored-openssl
cargo install cargo-run-script


```

## CosmWasm Build Contract

```bash
# Build the template
cargo generate --git https://github.com/CosmWasm/cw-template.git --name validator-reviews-contract -d minimal=true --tag a2a169164324aa1b48ab76dd630f75f504e41d99

# run the build optimizer (from source -> contract wasm binary)
cargo run-script optimize
```

## Chain Setup

```bash
GITHUB_USERNAME=rollchains

spawn new rollchain \
--consensus=proof-of-stake \
--bech32=roll \
--denom=uroll \
--bin=rolld \
--disabled=block-explorer \
--org=${GITHUB_USERNAME}


# move into the chain directory
cd rollchain

# build module
spawn module new reviews
make proto-gen

# - Installs the binary
# - Setups the default keys with funds
# - Starts the chain in your shell
make sh-testnet
# TODO: setup rust if not already, read through wasm guides





```

## Test Deployment

```bash
rolld tx wasm store $HOME/Desktop/validator-reviews-contract/artifacts/validator_reviews_contract.wasm --from=acc0 --gas=auto --gas-adjustment=2.0 --yes
# rolld q wasm list-code

rolld tx wasm instantiate 1 '{}' --no-admin --from=acc0 --label="validator_reviews" --gas=auto --gas-adjustment=2.0 --yes
# rolld q wasm list-contracts-by-creator roll1hj5fveer5cjtn4wd6wstzugjfdxzl0xpg2te87

REVIEWS_CONTRACT=roll14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9sjczpjh

rolld q wasm state smart $REVIEWS_CONTRACT '{"validators":{}}'


MESSAGE='{"write_review":{"val_addr":"rollvaloper1hj5fveer5cjtn4wd6wstzugjfdxzl0xpmhf3p6","review":"hi this is a review"}}'
rolld tx wasm execute $REVIEWS_CONTRACT "$MESSAGE" --from=acc0 --gas=auto --gas-adjustment=2.0 --yes

rolld q wasm state smart $REVIEWS_CONTRACT '{"reviews":{"address":"rollvaloper1hj5fveer5cjtn4wd6wstzugjfdxzl0xpmhf3p6"}}'

```
