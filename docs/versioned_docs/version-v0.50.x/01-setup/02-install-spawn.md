---
title: Install Spawn
sidebar_label: Install Spawn
sidebar_position: 2
slug: /install/install-spawn
---


# Overview

:::note Synopsis
Install the Spawn CLI tool to your local machine
:::


## Install Spawn

Install Spawn from source.

```bash
# Install from latest source
git clone https://github.com/rollchains/spawn.git --depth 1 --branch v0.50.7

# Change to this directory
cd spawn

# Install Spawn
make install

# Install Local Interchain (testnet runner)
make get-localic

# Verify installations were successful
spawn

local-ic

# If you get "command 'spawn' not found", run the following
# Linux / Windows / Some MacOS
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.bashrc
source ~/.bashrc

# MacOS
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.zshrc
source ~/.zshrc

# Legacy MacOS Go
echo 'export PATH=$PATH:$HOME/go/bin' >> ~/.zshrc
source ~/.zshrc

# Sometimes it can be good to also clear your cache
# especially WSL users
go clean -cache
```
