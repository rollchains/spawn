---
title: "IBC Transfers"
sidebar_label: "IBC Transfers"
slug: /demo/ibc
# sidebar_position: 1
---

# IBC Demo

In this tutorial, we'll create and interact with a new Cosmos-SDK blockchain called "rollchain", with the token denomination "uroll". This chain has tokenfactory and Proof of Authority, but we'll disable cosmwasm.

## Prerequisites
- [System Setup](../01-setup/01-system-setup.md)
- [Install Spawn](../01-setup/02-install-spawn.md)

## Create your chain

Create your chain using the spawn command line tool. Change the `GITHUB_USERNAME` to your github username.
If you do not have a github, leaving it as default is fine.

Create a [Proof of Authority](https://en.wikipedia.org/wiki/Proof_of_authority) network to focus on the application logic rather than worry about the security model. This is a great starting point for new chains.

```bash
GITHUB_USERNAME=rollchains

# If the `--consensus` or `--disabled` flags are not present,
# a selector UI will appear in your terminal to see all options.
spawn new rollchain \
--consensus=proof-of-authority \
--bech32=roll \
--denom=uroll \
--bin=rolld \
--disabled=cosmwasm,block-explorer \
--org=${GITHUB_USERNAME}
```

> *NOTE:* `spawn` creates a ready to use repository complete with `git` and GitHub CI. It can be quickly pushed to a new repository getting you and your team up and running quickly.

## Spin up an IBC testnet

The `chains/testnet.json` file contains the configuration for the testnet. It is a simple JSON file that contains the chain configurations for the testnet. By default it starts 2 networks, configures a relayer, and connects the two chains via IBC.

```bash
cd rollchain

# Starts 2 networks for the IBC testnet at http://127.0.0.1:8080.
# - Builds the docker image of your chain
# - Launches a testnet with IBC automatically connected and relayed
#
# Note: you can run a single node, non IBC testnet, with `make sh-testnet`.
make testnet
```

## Send a Transaction

```bash
# list the keys that have been provisioned with funds at launch
rolld keys list

# send a transaction from one account to another
rolld tx bank send acc0 $(rolld keys show acc1 -a) 1337uroll --chain-id=localchain-1

# enter "y" to confirm the transaction
# then query your balances for proof the transaction executed successfully
rolld q bank balances $(rolld keys show acc1 -a)
```

## Send an IBC transaction

```bash
# submit a cross chain transfer from acc0 to the other address
rolld tx ibc-transfer transfer transfer channel-0 cosmos1hj5fveer5cjtn4wd6wstzugjfdxzl0xpxvjjvr 7uroll --from=acc0 --chain-id=localchain-1 --yes

# Query the other side to confirm it went through
sleep 10

# Interact with the other chain without having to install the cosmos binary
# - Endpoints found at: GET http://127.0.0.1:8080/info
# - make get-localic
local-ic interact localcosmos-1 query 'bank balances cosmos1hj5fveer5cjtn4wd6wstzugjfdxzl0xpxvjjvr'
```

## (optional) Push to GitHub

```bash
# Create a new repository on GitHub from the gh cli
gh repo create rollchain --source=. --remote=upstream --push --private
```

> You can also push it the old fashioned way with https://github.com/new

## Conclusion

In this tutorial, you configured a new custom chain, launched a testnet for it, tested a cross chain token transfer, and confirmed it was successful. This tutorial demonstrates just how easy it is to create a brand new custom Cosmos-SDK blockchain from scratch with spawn.
