---
title: "Name Service"
sidebar_label: "Build a Name Service"
sidebar_position: 1
slug: /build/name-service
---


# Overview

:::note Synopsis
Building your first Cosmos-SDK blockchain with Spawn. This tutorial focuses on a 'nameservice' where you set your account to a name you choose.

* Generating a new chain
* Creating a new module
* Adding custom logic
* Run locally
* Interacting with the network
:::

## Prerequisites
- [System Setup](../01-setup/01-system-setup.md)
- [Install Spawn](../01-setup/02-install-spawn.md)

## Video Walkthrough

<iframe width="560" height="315" src="https://www.youtube.com/embed/4gFSuLUlP4I?si=A_VqEwhOh2ZPxNsb" title="YouTube video player" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share; fullscreen" referrerpolicy="strict-origin-when-cross-origin" allowfullscreen></iframe>

## Generate a New Chain

Let's create a new chain called 'rollchain'. You are going to set defining characteristics such as
- Which modules to disable from the template *if any*
- Proof of Stake consensus
- Wallet prefix (bech32)
- Token name (denom)
- Binary executable (bin)

```bash
spawn new rollchain --consensus=pos --disable=cosmwasm --bech32=roll --denom=uroll --bin=rolld
```

🎉 Your new blockchain 'rollchain' is now generated!

## Scaffold the Module
Now it is time to build the nameservice module structure.

Move into the 'rollchain' directory and generate the new module with the following commands:

```bash
# moves into the rollchain directory you just generated
cd rollchain

# scaffolds your new nameservice module
spawn module new nameservice

# proto-gen proto files to go
#
# If you get a /.cache permission error, run:
# sudo chmod -R 777 $(pwd)/.cache
# sudo chown -R $USER $(pwd)/.cache
#
# If you get a cannot find module error
# go clean -modcache
#

make proto-gen
```

This creates a new template module with the name `nameservice` in the `x` and `proto` directories. It also automatically connected to your application and is ready for use.
