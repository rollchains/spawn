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
if [ `ps -p $$ -o 'comm='` == "bash" ]; then
  echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.bashrc
  source ~/.bashrc
elif [ `ps -p $$ -o 'comm='` == "zsh" ]; then
  echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.zshrc
  source ~/.zshrc
fi

```
