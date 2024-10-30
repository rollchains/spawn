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

## How a Name Service works

Imagine you have a set of labeled containers, each with a unique name like "Essentials" "Electronics" and "Books". The label on each container is called the key, and what’s stored inside is the value. For example, the key "Books" leads you to a container full of books, while "Essentials" might have things like toiletries or basic supplies.

In a nameservice, this key-value system lets you quickly find or access specific data by referencing a unique name, or key, which reliably points you to the related value. This is useful for mapping names to specific information or resources, so with just the name, you can always find exactly what you're looking for.

For this tutorial we map a human readable name (like `"alice"`) to a complex wallet address (like `roll1efd63aw40lxf3n4mhf7dzhjkr453axur57cawh`) so it is easier to understand and view as a user.

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
# If you get a cannot find module error
# go clean -modcache
make proto-gen
```

This creates a new template module with the name `nameservice` in the `x` and `proto` directories. It also automatically connected to your application and is ready for use.
