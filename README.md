<div align="center">
  <h1>Spawn</h1>
</div>

Spawn is the best development platform for building custom modular Cosmos-SDK blockchains. Pick and choose modules to create a network tailor-fit to your needs. Use the native Cosmos tools and standards you're already familiar with. Quickly test interoperability between your new chain and established networks like the Cosmos-Hub across local devnet, public testnet, and mainnet through native InterchainTest support. Take advantage of the latest innovations, such as Proof-Of-Authority consensus and Celestia Data Availability. Build without limits.

## Spawn in Action

In this 4 minute demo Jack:
- Creates a new chain, customizing the modules and genesis
- Creates a new `nameservice` module
  - Adds the new message structure for transactions and queries
  - Stores the new data types
  - Adds the application logic
  - Connect it to the command line
- Runs the `sh-testnet` to build and launch the chain locally
- Interacts with the chain's nameservice logic, settings a name, and retrieving it

https://github.com/rollchains/spawn/assets/31943163/ecc21ce4-c42c-4ff2-8e73-897c0ede27f0

## Requirements

- `go 1.21+` - [official site](https://go.dev/doc/install)
- Docker - [official site](https://docs.docker.com/get-docker/)

## Getting Started
In this tutorial, we'll create and interact with a new Cosmos-SDK blockchain called "rollchain", with the token denomination "uroll". This chain has tokenfactory, POA, and globalfee modules enabled, but we'll disable cosmwasm.

1. Clone this repo and install

```shell
git clone https://github.com/rollchains/spawn.git
cd spawn
make install
```

2. Create your chain using the `spawn` command and customize it to your needs!

```shell
GITHUB_USERNAME=rollchains

spawn new rollchain \
--bech32=roll `# the prefix for addresses` \
--denom=uroll `# the coin denomination to create` \
--bin=rolld `# the name of the binary` \
--disabled=cosmwasm `# modules to disable. [proof-of-authority,tokenfactory,globalfee,packetforward,cosmwasm,wasm-lc,ignite]` \
--org=${GITHUB_USERNAME} `# the github username or organization to use for the module imports`
```

> *NOTE:* `spawn` creates a ready to use repository complete with `git` and GitHub CI. It can be quickly pushed to a new repository getting you and your team up and running quickly.

3. Spin up a local testnet for your chain

```shell
cd rollchain
make testnet
```

4. Open a new terminal window and send a transaction on your new chain

```shell
# list the keys that have been provisioned with funds in genesis
rolld keys list

# send a transaction from one account to another
rolld tx bank send acc0 $(rolld keys show acc1 -a) 1337uroll --chain-id=chainid-1

# enter "y" to confirm the transaction
# then query your balances tfor proof the transaction executed successfully
rolld q bank balances $(rolld keys show acc1 -a)
```

5. (optional) Send an IBC transaction

```shell
# submit a cross chain transfer from acc0 to the other address
rolld tx ibc-transfer transfer transfer channel-0 cosmos1hj5fveer5cjtn4wd6wstzugjfdxzl0xpxvjjvr 7uroll --from=acc0 --chain-id=chainid-1 --yes

# Query the other side to confirm it went through
sleep 10

curl -X POST -H "Content-Type: application/json" -d '{
  "chain_id": "localcosmos-1",
  "action": "query",
  "cmd": "bank balances cosmos1hj5fveer5cjtn4wd6wstzugjfdxzl0xpxvjjvr"
}' http://127.0.0.1:8080/
```

6. Push your new chain to a new repository

> [Create a new repository on GitHub](https://github.com/new)

```shell
# git init, add, and commit are all handled by default unless you add the `--no-git` flag on create
git remote add origin https://github.com/{your_github_username}/rollchain.git
git push -u origin master
```
In this tutorial, we configured a new custom chain, launched a testnet for it, tested a simple token transfer, and then pushed the custom chain code to a new git repo. This tutorial demonstrates just how easy it is to create a brand new custom Cosmos-SDK blockchain from scratch, saving developers time.
