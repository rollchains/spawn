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
git clone https://github.com/rollchains/spawn.git --depth 1 --branch v0.50.10

# Change to this directory
cd spawn

# Clear Go modules cache for a fresh install
go clean -modcache

# Install Spawn
make install

# Install Local Interchain (testnet runner)
make get-localic

# Install docker container builder
make get-heighliner

# Verify installations were successful
spawn

local-ic

heighliner
```

## Command not found error

If you get "command 'spawn' not found", run:

```bash
# Gets your operating system
unameOut="$(uname -s)"
case "${unameOut}" in
    Linux*)     machine=Linux;;
    Darwin*)    machine=Mac;;
    CYGWIN*)    machine=Cygwin;;
    MINGW*)     machine=MinGw;;
    MSYS_NT*)   machine=MSys;;
    *)          machine="UNKNOWN:${unameOut}"
esac
echo "Your operating system is: $machine"
echo -e "\nAdding the go binary location to your PATH for global access.\n\tIt will now prompt you for your password."

# Adds the location of the binaries to your PATH for global execution.
cmd='export PATH=$PATH:$(go env GOPATH)/bin'
if [ $machine == "Linux" ]; then
    sudo echo "$cmd" >> ~/.bashrc && source ~/.bashrc
elif [ $machine == "Mac" ]; then
    sudo echo "$cmd" >> ~/.zshrc && source ~/.zshrc
else
    echo 'Please add the following to your PATH: $(go env GOPATH)/bin'
fi
```
