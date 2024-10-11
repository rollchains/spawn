---
slug: /
sidebar_position: 0
title: Meet Spawn
sidebar_label: "Spawn Documentation"
---

Spawn is the easiest way to build, maintain and scale a Cosmos SDK blockchain. Spawn solves all the key pain points engineers face when building new Cosmos-SDK networks.
  - **Tailor-fit**: Pick and choose modules to create a network for your needs.
  - **Commonality**: Use native Cosmos tools and standards you're already familiar with.
  - **Integrations**: Github actions and end-to-end testing are configured right from the start.
  - **Iteration**: Quickly test between your new chain and established networks like the local Cosmos-Hub devnet.

## NameService Demo

[Follow Along with the NameService demo](./02-build-your-application/01-nameservice.md)

<video src="https://github.com/rollchains/spawn/assets/31943163/ecc21ce4-c42c-4ff2-8e73-897c0ede27f0" width="100%" height="100%" controls></video>


## Testimonials

> "Spawn is a marked transformation in CosmosSDK protocol development, allowing scaffolding and upgrading from 0.47 to 0.50 to be achievable and understandable. Without the tool, this would have been a dedicated multi-month effort" - Ash, [Burnt.com](https://twitter.com/burnt_xion)


> "Spawn has truly streamlined the developer onboarding process into the Cosmos ecosystem, seamless and efficient." - [Anil](https://x.com/anilcse_/status/1840444855576846355) [VitWit](https://www.vitwit.com/)

---

## Spawn Overview

Setting up a new blockchain used to take at least a week, requiring manual edits, debugging, and configuring tests. Now, with Spawn, you can create a custom network in just a few clicks. It generates a personalized network tailored to your project, letting you focus on writing product logic. The modular approach allows you to include or remove features, so you can start building quickly without the hassle of setting up everything from scratch. Spawn simplifies the process, especially for new developers, by removing guesswork and speeding up the setup.

## New Development

Get started building using the `new-chain` command. Spawn will guide you through the process of selecting the modules you need and configuring your new chain. Using `--help` will showcase examples and other options you may want to consider for your new network.

```bash
spawn new mychain --help
```

```bash
Create a new project

Usage:
  spawn new-chain [project-name] [flags]

Aliases:
  new-chain, new, init, create

Flags:
  -b, --binary string          Application binary name (default "simd")
      --bypass-prompt          Bypass UI prompter
      --denom string           Bank token denomination (default "token")
      --org string             Github organization name (default "rollchains")
      --skip-git               No git repository created
      --wallet-prefix string   Users wallet namespace (default "cosmos")
```

### Security Selection

You can read about different security models in the [Consensus Security](./04-learn/01-consensus-algos.md) section. If you don't know which to select, use proof of authority.

```bash
spawn new mychain
```

After running the new command, navigate with your arrow keys and press 'enter' to select the module you want to use. You can only use 1 from this consensus list. Then select done.

```
Consensus Selector (( enter to toggle ))
  Done
  ✔ proof-of-authority
  proof-of-stake
  interchain-security
```

### Feature Selection

You now select which features you want to include in your base application. Usually you would have to do these manually, each taking about 15 minutes to integrate. With spawn, you select them right away. It automatically configures them **and** give you testing for the assurance it works.

An information guide will be displayed for each feature at the bottom of the UI, sharing information about what the feature does. Select the following then press 'enter' on done to continue.

```bash
Feature Selector (( enter to toggle ))
  Done
  ✔ tokenfactory
  ✔ ibc-packetforward
  ✔ ibc-ratelimit
  cosmwasm
  wasm-light-client
  ✔ optimistic-execution
  ✔ block-explorer
tokenfactory: Native token minting, sending, and burning on the chain
```

Just like that, an entire network is generated. Everything you need to get started and more! Let's dive in.

## Structure

Opening up this newly generated `mychain/` gives you a general view into the entire layout.

```bash title="ls -laG"
.github/
app/
chains/
cmd/
contrib/
explorer/
interchaintest/
proto/
scripts/

.gitignore
.goreleaser.yaml
chain_metadata.json
chain_registry_assets.json
chain_registry.json
chains.yaml
docker-compose.yml
Dockerfile
go.mod
go.sum
Makefile
README.md
```

### .github/

This directory contains all the workflow actions for native github integration out of the box. It handles
- Integration & Unit tests for every code change
- Docker images saved to [ghcr](https://github.com/features/packages) on a new version tag
- Public cloud or private hosted testnets
- App binary releases
- PR title formatting
- Markdown file valid link reviews

### app/

App is the main location for all of the application connection logic.

- **decorators/** - Initial logic as new transactions are received. Used to override input data, block requests, or add additional logic before the action begins initial processing.
- **upgrades/** - You have to run an upgrade when you add or remove logic and nodes are already running different logic. This is where you put the upgrade information and state migrations.
- **ante.go** - The decorators for the entire network, wired together.
- **app.go** - The entire application connected and given access to the cosmos-sdk. The brain of the program.
- **upgrades.go** - Registers the upgrades/ folder logic when one is pending processing.

### chains/

The chains/ directory is where the local and public testnet configuration files are placed. Reference the [testnets](#testnets) section for more information

### cmd/

The cmd/ directory is the entry point for the wiring connections and is where the main.go file is located. This is where the application is started and the chain is initialized when you run the binary. By default, `simd` is the binary name and is saved to your $GOPATH (/home/user/go/bin/).

### explorer/

If you enabled the explorer in the feature selection, this is where the [ping.pub](https://ping.pub/) explorer files are located. When running a testnet with `make sh-testnet` or `make testnet`, you can launch the explorer along side the chain to view activity in real time. Blocks, transactions, uptime, connections, and more are all viewable. Easily launch it with the `docker compose up` command in the root of the directory.

### interchaintest/

Interchaintest is a generalized integration test environment for the Interchain and beyond. It supports Cosmos, Ethereum, UTXO (Bitcoin), and other chain types. By default you will see many test like `ibc_test.go`, `ibc_rate_limit_test.go` and `tokenfactory_test.go` after generation. Any features you select are placed here automatically to confirm your network is working as expected. This are run with the github action automatically on every code change **or** you can run them manually with `make local-image && make ictest-*`, where the * is the testname *(ictest-ibc, ictest-tokenfactory, etc)*.

### proto/

[Proto, also called protocol buffers](https://protobuf.dev/), are a generalized way to define the structure of data. Discussed this more in the [Modules](#modules) sub section.

### scripts/

Scripts automate some more complex requirements list setting up a fast testnet or generating code on the fly. You should not need to modify anything here until you are more advanced. These are shown in the `make help` command to abstract away complexity.

### chain_metadata.json

A cosmetic file showcasing a format for the network. Fill in the data here once you push to the public so developers can easily see what your network is about. This is required for [ICS consumer networks](./04-learn/01-consensus-algos.md#create-an-ics-consumer-network). If you do not use ICS, you can delete this file if you wish.

### chain_registry.json & assets

These files are the format needed to upload to [https://cosmos.directory/](https://cosmos.directory/) ([github](https://github.com/cosmos/chain-registry)). Frontends use this data to connect to the network, especially in the [local-interchain testnet tool](#testnets).

## Modules

We're all here to build new logic on top. The SDK calls these modules, or e**x**tensions, x/ for short. To make this easy spawn has a build in generator for a module.

```bash
spawn module new --help
```

```bash
Usage:
  spawn module new [name] [flags]

Aliases:
  new, c, create

Examples:
  spawn module new mymodule [--ibc-module]

Flags:
  --ibc-middleware   Set the module as an IBC Middleware
  --ibc-module       Set the module as an IBC Module
```

All you need to have is the name you wish to call it, and if you want standard or an IBC module. IBC enables cross network communication of the logic. This is a powerful feature that allows you to build a network of networks. You can try this out with the [IBC module demo](./02-build-your-application/08-ibc-module.md) demo.

For now, just create a default module called `example`

```bash
spawn module new example
```

```
🎉 New Module 'example' generated!
🏅 Commands:
  - $ make proto-gen     # convert proto files into code
```

This created a new x/example module and the [proto/](#proto) files in the expected structure. `genesis.proto` contains the data saved and more hardcoded. `query.proto` is how you allow external actors to grab data from the network and `tx.proto` is how you allow external actors to send data to the network. Spawn also connects it to the application if you look through your `app/app.go`.

Learn how to make a new module with the [Name Service](./02-build-your-application/01-nameservice.md) guide.

## Testnets

This uses the [local-interchain](https://github.com/strangelove-ventures/interchaintest/tree/main/local-interchain) format and supports JSON or YAML. By default, 2 IBC network defaults are included. **self-ibc** and **testnet**. Run the testnet with `make testnet` to automatically build, setup, and launch a complex network simply.

Self IBC is really only useful if you are building [IBC Modules](./02-build-your-application/08-ibc-module.md). Follow that guide to see how to use it.
