<div align="center">
  <h1>Spawn</h1>
</div>

Spawn is the best development platform for building custom modular Cosmos-SDK blockchains. Pick and choose modules to create a network tailor-fit to your needs. Use the native Cosmos tools and standards you're already familiar with. Quickly test interoperability between your new chain and established networks like the Cosmos-Hub across local devnet, public testnet, and mainnet through native InterchainTest support. Take advantage of the latest innovations, such as Proof-Of-Authority consensus and Celestia Data Availability. Build without limits.

## Spawn in Action

In this demo @Reecepbcups:
- Downloads and installs `spawn`
- Creates a new chain, customizing the modules and genesis
- Installs the chain binary and spins up a local testnet
- Interacts with the chain using the `appd` CLI

https://github.com/rollchains/spawn/assets/10821110/0de7bf37-c82a-4826-a3e3-13def6a53327

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
spawn new rollchain \
--bech32=roll `# the prefix for addresses` \
--denom=uroll `# the coin denomination to create` \
--bin=rolld `# the name of the binary` \
--disabled=cosmwasm `# modules to disable. By default all modules are enabled [tokenfactory, PoA, globalfee, cosmwasm]` \
--org={your_github_username} `# the github username or organization to use for the module imports`
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

5. Push your new chain to a new repository

> [Create a new repository on GitHub](https://github.com/new)

```shell
# git init, add, and commit are all handled by default unless you add the `--no-git` flag on create
git remote add origin https://github.com/{your_github_username}/rollchain.git
git push -u origin master
```
In this tutorial, we configured a new custom chain, launched a testnet for it, tested a simple token transfer, and then pushed the custom chain code to a new git repo. This tutorial demonstrates just how easy it is to create a brand new custom Cosmos-SDK blockchain from scratch, saving developers time.

## Spawn Product Goals

- Easy templating for a new chain from base

- Local-Interchain nested, easy JSON configured starts

- Chain-Chores like features (pull all patches from upstream based off the value of the current spawn instance. i.e. spawn v1.0 pulls from the v1.0 branch)

- Easily add CI/CD in line with the template repo (could just pull from this Repos CI so we can confirm all works upstream. Then wget down)

- Base for a new module into your repo (spawn module new <module-name>). Regex import into your apps SDK without any comments in the app.go

- Easily import or clone upstream modules into your repo (spawn module import <module-name>). Module name can also be a git repo (even one we do not own) that we can pull the files from directly for the user. So if we want SDK v0.50 tokenfactory, we can pull from repo X or repo Y depending on our needs. May require a unique go.mod for each module, unsure atm. Maybe we can abstract this away and handle ourselves?
