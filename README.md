<div align="center">
  <h1>Spawn</h1>
</div>

Spawn is the easiest way to build, maintain and scale a Cosmos SDK blockchain. Spawn solves all the key pain points engineers face when building new Cosmos-SDK networks.
  - **Tailor-fit**: Pick and choose modules to create a network for your needs.
  - **Commonality**: Use native Cosmos tools and standards you're already familiar with.
  - **Integrations**: Github actions and end-to-end testing are configured right from the start.
  - **Iteration**: Quickly test between your new chain and established networks like the local Cosmos-Hub devnet.

## Spawn in Action

In this 4 minute demo we:
- Create a new chain, customizing the modules and genesis
- Create a new `nameservice` module
  - Add the new message structure for transactions and queries
  - Store the new data types
  - Add the application logic
  - Connect it to the command line
- Runs `make sh-testnet` to build and launch the chain locally
- Interacts with the chain's nameservice logic, settings a name, and retrieving it

https://github.com/rollchains/spawn/assets/31943163/ecc21ce4-c42c-4ff2-8e73-897c0ede27f0

## Requirements

- [`go 1.22+`](https://go.dev/doc/install)
- [`Docker`](https://docs.docker.com/get-docker/)

[MacOS + Ubuntu Setup](./docs/SYSTEM_SETUP.md)

---

## Getting Started
In this tutorial, we'll create and interact with a new Cosmos-SDK blockchain called "rollchain", with the token denomination "uroll". This chain has tokenfactory and Proof of Authority, but we'll disable cosmwasm.

1. Clone this repo and install

```shell
git clone https://github.com/rollchains/spawn.git
cd spawn
git checkout v0.50.4
make install
```

2. Create your chain using the `spawn` command and customize it to your needs!

```shell
GITHUB_USERNAME=rollchains

# Available Features:
# * tokenfactory,globalfee,ibc-packetforward,ibc-ratelimit,cosmwasm,wasm-light-client,optimistic-execution,ignite-cli,block-explorer

spawn new rollchain \
--consensus=proof-of-authority `# proof-of-authority,proof-of-stake,interchain-security` \
--bech32=roll `# the prefix for addresses` \
--denom=uroll `# the coin denomination to create` \
--bin=rolld `# the name of the binary` \
--disabled=cosmwasm,globalfee,block-explorer `# disable features.` \
--org=${GITHUB_USERNAME} `# the github username or organization to use for the module imports, optional`


```

> *NOTE:* `spawn` creates a ready to use repository complete with `git` and GitHub CI. It can be quickly pushed to a new repository getting you and your team up and running quickly.

3. Spin up a local testnet for your chain

```shell
cd rollchain

# Starts 2 networks for the IBC testnet at http://127.0.0.1:8080.
# - Builds the docker image of your chain
# - Launches a testnet with IBC automatically connected and relayed
#
# Note: you can run a single node, non IBC testnet, with `make sh-testnet`.
make testnet
```

4. Open a new terminal window and send a transaction on your new chain

```shell
# list the keys that have been provisioned with funds in genesis
rolld keys list

# send a transaction from one account to another
rolld tx bank send acc0 $(rolld keys show acc1 -a) 1337uroll --chain-id=localchain-1

# enter "y" to confirm the transaction
# then query your balances tfor proof the transaction executed successfully
rolld q bank balances $(rolld keys show acc1 -a)
```

5. (optional) Send an IBC transaction

```shell
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

```shell
# Create a new repository on GitHub from the gh cli
gh repo create rollchain --source=. --remote=upstream --push --private
```

> You can also push it the old fashioned way with https://github.com/new

In this tutorial, we configured a new custom chain, launched a testnet for it, tested a simple token transfer, and confirmed it was successful. This tutorial demonstrates just how easy it is to create a brand new custom Cosmos-SDK blockchain from scratch, saving developers time.

## Add a Feature

If you wish to add a feature to spawn, reference the [NEW_FEATURE.md](./docs/NEW_FEATURE.md) for guidance.