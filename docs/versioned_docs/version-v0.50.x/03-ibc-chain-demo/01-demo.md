---
title: "IBC Demo"
sidebar_label: "IBC Network"
sidebar_position: 1
slug: /build/ibc-network
---

# Getting Started

In this tutorial, we'll create and interact with a new Cosmos-SDK blockchain called "rollchain", with the token denomination "uroll". This chain has tokenfactory and Proof of Authority, but we'll disable cosmwasm.

1. Clone this repo and install

```bash
git clone https://github.com/rollchains/spawn.git --depth 1 --branch v0.50.7
cd spawn

make install
make get-localic

# If you get "command 'spawn' not found", add to path
# Run the following in your terminal to test
# Then add to ~/.bashrc (linux) or ~/.zshrc (mac)
export PATH="$PATH:$(go env GOPATH)/bin"
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.bashrc
```

2. Create your chain using the `spawn` command and customize it to your needs!

```bash
GITHUB_USERNAME=rollchains

# If the `--consensus` or `--disabled` flags are not present,
# a selector UI will appear in your terminal to see all options.
spawn new rollchain \
--consensus=proof-of-authority \
--bech32=roll \
--denom=uroll \
--bin=rolld \
--disabled=cosmwasm,globalfee,block-explorer \
--org=${GITHUB_USERNAME}
```

> *NOTE:* `spawn` creates a ready to use repository complete with `git` and GitHub CI. It can be quickly pushed to a new repository getting you and your team up and running quickly.

3. Spin up a local testnet for your chain

```bash
cd rollchain

# Starts 2 networks for the IBC testnet at http://127.0.0.1:8080.
# - Builds the docker image of your chain
# - Launches a testnet with IBC automatically connected and relayed
#
# Note: you can run a single node, non IBC testnet, with `make sh-testnet`.
make testnet
```

4. Open a new terminal window and send a transaction on your new chain

```bash
# list the keys that have been provisioned with funds in genesis
rolld keys list

# send a transaction from one account to another
rolld tx bank send acc0 $(rolld keys show acc1 -a) 1337uroll --chain-id=localchain-1

# enter "y" to confirm the transaction
# then query your balances for proof the transaction executed successfully
rolld q bank balances $(rolld keys show acc1 -a)
```

5. (optional) Send an IBC transaction

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

6. Push your new chain to a github repository

```bash
# Create a new repository on GitHub from the gh cli
gh repo create rollchain --source=. --remote=upstream --push --private
```

> You can also push it the old fashioned way with https://github.com/new

In this tutorial, we configured a new custom chain, launched a testnet for it, tested a simple token transfer, and confirmed it was successful. This tutorial demonstrates just how easy it is to create a brand new custom Cosmos-SDK blockchain from scratch, saving developers time.
