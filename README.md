<div align="center">
  <h1>Spawn</h1>
</div>

Spawn is the best development platform for building custom modular Cosmos-SDK blockchains. Pick and choose modules to create a network tailor-fit to your needs. Use the native Cosmos tools and standards you're already familiar with. Quickly test interoperability between your new chain and established networks like the Cosmos-Hub across local devnet, public testnet, and mainnet through native InterchainTest support. Take advantage of the latest innovations, such as Proof-Of-Authority consensus and Celestia Data Availability. Build without limits.

## 1 Minute Demo

https://github.com/rollchains/spawn/assets/10821110/0de7bf37-c82a-4826-a3e3-13def6a53327

## Getting Started
1. Install `go 1.21+` - [official site](https://go.dev/doc/install)
2. Clone this repo and install
    ```shell
    $ git clone https://github.com/rollchains/spawn.git
    $ cd spawn
    $ make install
    ```
3. Scaffold a new custom chain
    ```shell
    $ spawn new customchain --bech32=customchain --disable=tokenfactory
    ```
4. Install the binary, and instantiate a testnet for your new chain
   ```shell
    $ cd customchain
    $ make testnet
    ```
5. Open a new terminal and interact with your new chain. The command defaults to `appd`
   ```shell
    $ appd keys list
    # copy the address for acc1 and replace the "..." in the next line
    $ appd tx bank send acc0 ... 1stake --chain-id=chainid-1
    # copy the tx hash and replace the "..." in the next line
    $ appd q tx ...
    ```

## Goals
- Easy templating for a new chain from base

- Local-Interchain nested, easy JSON configured starts

- Chain-Chores like features (pull all patches from upstream based off the value of the current spawn instance. i.e. spawn v1.0 pulls from the v1.0 branch)

- Easily add CI/CD in line with the template repo (could just pull from this Repos CI so we can confirm all works upstream. Then wget down)

- Base for a new module into your repo (spawn module new <module-name>). Regex import into your apps SDK without any comments in the app.go

- Easily import or clone upstream modules into your repo (spawn module import <module-name>). Module name can also be a git repo (even one we do not own) that we can pull the files from directly for the user. So if we want SDK v0.50 tokenfactory, we can pull from repo X or repo Y depending on our needs. May require a unique go.mod for each module, unsure atm. Maybe we can abstract this away and handle ourselves?
