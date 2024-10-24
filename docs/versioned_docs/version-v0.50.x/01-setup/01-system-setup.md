---
title: Setup Development Environment
sidebar_label: System Setup
sidebar_position: 1
slug: /install/system-setup
---

import Tabs from '@theme/Tabs';
import TabItem from '@theme/TabItem';

# Overview

:::note Synopsis
Setup your development environment with the essentials to get started building the blockchain.
:::


## System Requirements

Before you can install and interact with spawn, you must have the following core tools installed:
* [`Go 1.22+`](https://go.dev/doc/install)
* [`Docker`](https://docs.docker.com/get-docker/)
* [`Git`](https://git-scm.com/)

If you do not have these components installed, follow the instructions below to install them.

## Install Dependencies

<Tabs defaultValue="macos">
  <TabItem value="macos" label="MacOS">
  ```bash
  # Base
  brew install make
  brew install gcc
  brew install wget
  brew install jq

  # Github CLI - https://github.com/cli/cli
  brew install gh
  gh auth login

  # Golang
  brew install go

  # Docker
  brew install --cask docker
  open -a Docker # start docker desktop
  # settings -> General -> Start Docker Desktop when you sign in to your computer
  # Apply & Restart

  # Setup base git config
  git config --global user.email "yourEmail@gmail.com"
  git config --global user.name "Your Name"
  ```
  </TabItem>

  <TabItem value="windows" label="Windows (WSL)" default>
  ```bash
  # Install WSL in powershell
  wsl --install
  Restart-Computer

  # Setup WSL Ubuntu Image
  wsl.exe --install Ubuntu-24.04

  # Open wsl instance
  wsl

  # update and add snap if not already installed
  sudo apt update && sudo apt install snapd

  # Install Go (Snap)
  sudo snap install go --channel=1.23/stable --classic

  # Install Base
  sudo apt install make gcc git jq wget

  # Install github-cli
  sudo snap install gh

  # Install docker
  https://docs.docker.com/desktop/wsl/#turn-on-docker-desktop-wsl-2
  # or snap:
  sudo snap install docker

  # Fix versioning for interaction of commands
  sudo chmod 666 /var/run/docker.sock

  # Setup base git config
  git config --global user.email "yourEmail@gmail.com"
  git config --global user.name "Your Name"
  ```
  </TabItem>

  <TabItem value="ubuntu-linux" label="Linux (Ubuntu)">
  <!-- markdown-link-check-disable -->
  ```bash
  # Base
  sudo apt install make gcc git jq wget

  # (optional) Github CLI - https://github.com/cli/cli
  curl -sS https://webi.sh/gh | sh
  gh auth login

  # Golang
  GO_VERSION=1.23.0
  wget https://go.dev/dl/go$GO_VERSION.linux-amd64.tar.gz
  sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go$GO_VERSION.linux-amd64.tar.gz

  # Docker
  sudo apt -y install docker.io

  # Setup base git config
  git config --global user.email "yourEmail@gmail.com"
  git config --global user.name "Your Name"
  ```
  </TabItem>
  <!-- markdown-link-check-enable -->

  <TabItem value="cosmwasm-rust" label="CosmWasm (Rust)">
  Some tutorials require CosmWasm (Rust smart contracts) setup. This section is option, unless a tutorial is CosmWasm focused.

  CosmWasm requires [Rust](https://www.rust-lang.org/). You must have this installed as the contract will be built locally.
  ```bash
  # Install rust - https://www.rust-lang.org/tools/install
  curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh

  # Update shell env
  source $HOME/.cargo/env

  # or Update if you have it
  rustup update

  # Install other dependencies
  rustup target add wasm32-unknown-unknown

  cargo install cargo-generate --features vendored-openssl
  cargo install cargo-run-script
  ```
  </TabItem>
</Tabs>
