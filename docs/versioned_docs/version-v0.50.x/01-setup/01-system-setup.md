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

Install [VSCode](https://code.visualstudio.com/download) if you do not already have a file editor. Once installed, the following extensions are useful
- [Golang Syntax Highlighting](https://marketplace.visualstudio.com/items?itemName=golang.Go)
- [Protobuf 3 Highlighting](https://marketplace.visualstudio.com/items?itemName=zxh404.vscode-proto3)
- [Improved Errors](https://marketplace.visualstudio.com/items?itemName=usernamehw.errorlens)

<Tabs defaultValue="macos">
  <TabItem value="macos" label="MacOS">

  ```bash
  # Setup Homebrew (https://brew.sh/)
  /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
  sudo echo 'export PATH=$PATH:/opt/homebrew/bin' >> ~/.zshrc

  # Base
  brew install make wget jq gh gcc go

  # (optional) Github CLI login
  # gh auth login

  # Docker
  brew install --cask docker

  # Start docker desktop
  open -a Docker
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

  # Setup WSL Ubuntu Image
  wsl.exe --install -d Ubuntu-24.04

  # Open wsl instance
  wsl

  # update and add snap if not already installed
  sudo apt update && sudo apt install snapd

  # Install Go (https://go.dev/wiki/Ubuntu)
  sudo apt update
  sudo apt install golang -y

  # Install Base
  sudo apt install make gcc git jq wget -y

  # Optional: Install github-cli
  # sudo snap install gh

  # Install docker
  # If you can't use snap, setup from:
  # - https://docs.docker.com/desktop/wsl/#turn-on-docker-desktop-wsl-2
  sudo snap install docker

  # Fix versioning for interaction of commands
  sudo chmod 666 /var/run/docker.sock

  # Setup base git config
  git config --global user.email "yourEmail@gmail.com"
  git config --global user.name "Your Name"
  ```

  After installing VSCode and having your system setup, open a terminal and run `code` to open vscode. You can open the current directory with `code .`, where `.` in computer terms stands for the current directory.

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
