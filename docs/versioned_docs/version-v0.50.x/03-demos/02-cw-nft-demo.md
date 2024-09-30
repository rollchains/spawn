---
title: "CW NFTs"
sidebar_label: "CosmWasm NFTs"
slug: /demo/cw-nft
# sidebar_position: 1
---

# Non-fungible Token Demo

You will build a new chain with [CosmWasm](https://cosmwasm.com/), enabling support for smart contracts on a new Cosmos-SDK application. You will download a pre-built contract, upload it, and interact with it to transfer the ownership of some data.

If you do not know what an NFT is, you can read about them here: [investopedia.com/non-fungible-tokens-nft](https://www.investopedia.com/non-fungible-tokens-nft-5115211).

:::note Warning
Some parts of this tutorial will not have the added context about spawn's inner workings or how commands work. Run through [Build Your Application](../02-build-your-application/01-nameservice.md) for this context.
:::

## Prerequisites
- [System Setup](../01-setup/01-system-setup.md)
- [Install Spawn](../01-setup/02-install-spawn.md)

:::note Danger
Some machines like Windows will not work with running the testnet. This is a limitation of the operating system hardware with wasm and required C language libraries / DLLs.

For the best experience, try `make testnet` or use a Linux machine or a cloud-based linux instance from [Hetzner](https://www.hetzner.com/cloud/) or [Digital Ocean](https://www.digitalocean.com/pricing/droplets) for $6 per month.
:::


## Create your chain

Build a new chain that has CosmWasm configured.

```bash
GITHUB_USERNAME=rollchains

spawn new rollchain \
--consensus=proof-of-stake \
--bech32=roll \
--denom=uroll \
--bin=rolld \
--disabled=block-explorer \
--org=${GITHUB_USERNAME}
```

<details>

<summary>View UI Selector</summary>

If you remove the `--disabled` flag; a more intuitive UI selection approach will be taken. Make sure CosmWasm is selected with the green arrow, then press `done`.

![Image](https://github.com/user-attachments/assets/16698f3f-143b-4258-9ff2-fc429764b58c)

</details>


## Start the testnet

:::note Note
If `make sh-testnet` does not start due to a port bind error, you can kill your previously running testnet with `killall -9 rolld`, then try again.
:::


```bash
# move into the chain directory
cd rollchain

# - Installs the binary
# - Setups the default keys with funds
# - Starts the chain in your shell
make sh-testnet
```

## Verify CosmWasm is enabled

```bash
rolld q wasm params
```

<details>

<summary>Expected Output</summary>

```bash
code_upload_access:
  addresses: []
  permission: Everybody
  instantiate_default_permission: Everybody
```

</details>

## Upload the contract to the network

You will use the [CW721](https://github.com/public-awesome/cw-nfts) contract for your NFT journey. CW721 stands for CosmWasm 721. [721 corresponds to the Ethereum specification for NFTs](https://www.coinbase.com/learn/crypto-glossary/what-is-erc-721). Understanding this is out of scope for this tutorial. Just know you can create, transfer, and query data.

Download the contract code from the CosmWasm NFTs repository, then upload it to the network with the application binary.

```bash
# Download the the NFT contract to your machine
curl -LO https://github.com/public-awesome/cw-nfts/releases/download/v0.19.0/cw721_base.wasm

# Upload the source code to the chain
# - gas is is amount of compute resources to allocate.
rolld tx wasm store ./cw721_base.wasm --from=acc0 \
    --gas=auto --gas-adjustment=2.0 --yes
```

## Verify the code was uploaded

```bash
# Code id: "1" is available
rolld q wasm list-code

# See the details (A lot of spam)
rolld q tx 4601FBACBDF93E4309D92E968F8B4E7F9177BCB132B65AA363AFDC26FE6B5CB6
```

<details>

<summary>Expected Code Info</summary>

```bash
(main) -> $ rolld q wasm list-code
code_infos:
- code_id: "1"
  creator: roll1hj5fveer5cjtn4wd6wstzugjfdxzl0xpg2te87
  data_hash: E13AA30E0D70EA895B294AD1BC809950E60FE081B322B1657F75B67BE6021B1C
  instantiate_permission:
    addresses: []
    permission: Everybody
pagination:
  next_key: null
  total: "0"
```

</details>



## Create a new NFT collection

With the source now uploaded, anyone can create a new NFT collection with this base contract code now on the chain. This will be a new contract that only you control. Now, instantiate the contract to create the new NFT collection.

<details>

<summary>Instantiate Format Source</summary>

You can find the instantiate, execute, and query messages (json) formats in the contract source code.

```rust reference title="packages/cw721/src/msg.rs"
https://github.com/public-awesome/cw-nfts/blob/v0.19.0/packages/cw721/src/msg.rs#L126-L143
```
</details>

:::note Warning
Notice the MESSAGE= below has no spaces in the JSON. This is required for the command line to parse it correctly. Failure to do so will result in the error

`ERR failure when running app err="accepts 2 arg(s), received 3"`
:::

```bash
# Get our account address for the acc0 wallet / key.
rolld keys show acc0 -a # roll1hj5fveer5cjtn4wd6wstzugjfdxzl0xpg2te87

# Create the NFT collection with our account
# as the authorized minter / creator for new NFTs.
MESSAGE='{"name":"Roll","symbol":"ROLL","minter":"roll1hj5fveer5cjtn4wd6wstzugjfdxzl0xpg2te87"}'

# Create the NFT collection
rolld tx wasm instantiate 1 $MESSAGE --no-admin --from=acc0 --label="my-nft" \
    --gas=auto --gas-adjustment=2.0 --yes
```

## Contract address

A contract address is where all the collection and information is stored. It never changes and is the unique identifier for interaction. Think of this similar to a website, google.com always brings you to google search. `NFT_CONTRACT` is always the RollNFTs collection.

```bash
# View all contract addresses a wallet has created
rolld q wasm list-contracts-by-creator roll1hj5fveer5cjtn4wd6wstzugjfdxzl0xpg2te87

# The contract address for the NFT collection just created
NFT_CONTRACT=roll14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9sjczpjh
```



## Create an NFT in the collection

The acc0 account now must create the first NFT in the collection since it is the minter. Specify the unique ID (1), the owner (acc0), and some data to be associated with this NFT. Set the url of a sunflower image as the metadata for this tutorial.

:::note Note
The `token_uri` is a URL that points to the metadata of the NFT. This can be a [JSON object](https://eips.ethereum.org/EIPS/eip-721#specification) or a URL to a JSON object.
This URL can be a link to an external service like [IPFS](https://ipfs.tech/), or the raw text directly. The contract does not care, it is up to you to manage the data and build the services around it.
:::

<details>

<summary>Execute Format Source</summary>

```rust reference title="packages/cw721/src/msg.rs"
https://github.com/public-awesome/cw-nfts/blob/v0.19.0/packages/cw721/src/msg.rs#L80-L91
```
</details>


```bash
MESSAGE='{"mint":{"token_id":"1","owner":"roll1hj5fveer5cjtn4wd6wstzugjfdxzl0xpg2te87","token_uri":"https://onlinejpgtools.com/images/examples-onlinejpgtools/sunflower.jpg"}}'

rolld tx wasm execute $NFT_CONTRACT $MESSAGE --from=acc0 \
    --gas=auto --gas-adjustment=2.0 --yes
```

## Grab this NFT data

There is now an NFT with the ID of 1 owned by `roll1hj5fveer5cjtn4wd6wstzugjfdxzl0xpg2te87`. Now query the contract to see the data and verify it is correct.

<details>

<summary>Query Format Source</summary>

```rust reference title="packages/cw721/src/msg.rs"
https://github.com/public-awesome/cw-nfts/blob/v0.19.0/packages/cw721/src/msg.rs#L157-L161
```

```rust reference title="packages/cw721/src/msg.rs"
https://github.com/public-awesome/cw-nfts/blob/v0.19.0/packages/cw721/src/msg.rs#L236-L240
```

</details>

```bash
# Get who is the owner of ID 1
rolld q wasm state smart $NFT_CONTRACT '{"owner_of":{"token_id":"1"}}'

# Retrieve the NFT info
rolld q wasm state smart $NFT_CONTRACT '{"nft_info":{"token_id":"1"}}'
```

## Transfer the NFT to another account

Now move the token from the originally minted account (acc0) to another account (acc1). This is a simple transfer of ownership to move who owns the data.

<details>

<summary>Execute Format Source</summary>

```rust reference title="packages/cw721/src/msg.rs"
https://github.com/public-awesome/cw-nfts/blob/v0.19.0/packages/cw721/src/msg.rs#L44-L48
```

</details>

```bash
# Recipient account
rolld keys show acc1 -a # roll1efd63aw40lxf3n4mhf7dzhjkr453axur57cawh

MESSAGE='{"transfer_nft":{"recipient":"roll1efd63aw40lxf3n4mhf7dzhjkr453axur57cawh","token_id":"1"}}'
rolld tx wasm execute $NFT_CONTRACT $MESSAGE --from=acc0 --gas=auto --gas-adjustment=2.0 --yes

# Get who is the owner of 1
# Moved to: roll1efd63aw40lxf3n4mhf7dzhjkr453axur57cawh
rolld q wasm state smart $NFT_CONTRACT '{"owner_of":{"token_id":"1"}}'
```

## Conclusion

In this tutorial, you built a new chain with CosmWasm enabled, launched a testnet for it, and launched an NFT collection! You uploaded a contract, created an NFT, and transferred it to another account. This is the foundation for building a new unique marketplace or game on the Interchain.

<!-- TODO: ICS721, cross chain NFTs -->
