---
title: "Spawns Purpose"
sidebar_label: "What Spawn Does"
slug: /learn/what-spawn-does
# sidebar_position: 1
---

# What Spawn Does

Learn about what spawn does, how it works, and how it can help you build your applications faster than ever before

## Overview

Spawn is a powerful tool designed to streamline the process of building, maintaining, and scaling Cosmos SDK blockchains. As a developer, Spawn offers several compelling reasons to incorporate it into your workflow.

Before, you would have at least a week setting up a new chain. You would manually copy paste some other codebase, modify it to your needs, debug issues, add custom test, fix github actions, and then you have a base network. Now you can do all this in just a flew clicks. Spawn takes in your context and with magic, generates a new network for you fitting your needs.

## New Development

Say you have a project idea, and you want to get started on writing the logic. Spawn gets you from 0 to 1 with the template matching your exact needs. The modular approach allows you to pick what you just have, such as smart contracts, and remove the things you don't. This way, you can focus on what matters to you.

Get started building your first chain using the `new-chain` command. Spawn will guide you through the process of selecting the modules you need and configuring your new chain.

```bash
spawn new mychain --help
```

Using `--help` will showcase examples and other options you may want to consider for your new network.

```bash
Create a new project

Usage:
  spawn new-chain [project-name] [flags]

Aliases:
  new-chain, new, init, create

Flags:
  -b, --binary string          binary name (default "simd")
      --bypass-prompt          bypass UI prompt
      --debug                  enable debugging
      --denom string           bank token denomination (default "token")
  -h, --help                   help for new-chain
      --org string             github organization (default "rollchains")
      --skip-git               ignore git init
      --wallet-prefix string   chain bech32 wallet prefix (default "cosmos")
```

### Security Selection

You can read about different security models in the [Consensus Security](./01-consensus-algos.md) section. If you don't know which to select, use proof of authority for your application.

```bash
spawn new mychain
```

After running the new command, use your arrow keys and use 'enter' to select the module you want to use. You can only use 1 from this list. Then select done.

```
Consensus Selector (( enter to toggle ))
  Done
  ✔ proof-of-authority
  proof-of-stake
  interchain-security
```

### Feature Selection

You now select which features you want to include in your base application. Usually you would have to do these manually, each taking about 15 minutes to add. With spawn, we let you select them right away, automatically configure them, **and** give you testing to give you the assurance it works.

As you scroll through features,

```bash
Feature Selector (( enter to toggle ))
  Done
  ✔ tokenfactory
  ✔ ibc-packetforward
  ✔ ibc-ratelimit
  cosmwasm
  wasm-light-client
  ✔ optimistic-execution
  ignite-cli
  ✔ block-explorer
tokenfactory: Native token minting, sending, and burning on the chain
```



---

# Testing

# Modules

# testnets
