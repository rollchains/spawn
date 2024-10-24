<div align="center">
  <h1>Spawn</h1>
</div>

Spawn is the easiest way to build, maintain and scale a Cosmos SDK blockchain. Spawn solves all the key pain points engineers face when building new Cosmos-SDK networks.
  - **Tailor-fit**: Pick and choose modules to create a network for your needs.
  - **Commonality**: Use native Cosmos tools and standards you're already familiar with.
  - **Integrations**: Github actions and end-to-end testing are configured right from the start.
  - **Iteration**: Quickly test between your new chain and established networks like the local Cosmos-Hub devnet.

## Documentation

<https://rollchains.github.io/spawn/v0.50/>

## Installation

### Prerequisite Setup

If you do not have [`go 1.22+`](https://go.dev/doc/install), [`Docker`](https://docs.docker.com/get-docker/), or [`git`](https://git-scm.com/) installed, follow the instructions below.

* [MacOS, Windows, and Ubuntu Setup](./docs/versioned_docs/version-v0.50.x/01-setup/01-system-setup.md)

### Install Spawn

```bash
# Download the the Spawn repository
git clone https://github.com/rollchains/spawn.git --depth=1 --branch v0.50.10
cd spawn

# Install Spawn
make install

# Install Local-Interchain (testnet runner)
make get-localic

# Attempt to run a command
spawn help

# Potential workaround if spawn is immediately killed when attempting run "spawn help" on a m* macbook.
# This may happen because of an xcode16 update.
# To work around this as a temporary measure please change this line in the Make file.
LD_FLAGS = -X main.SpawnVersion=$(VERSION)
# to
LD_FLAGS = -X main.SpawnVersion=$(VERSION) -s -w


# If you get "command 'spawn' not found", add to path.
# Run the following in your terminal to test
# Then add to ~/.bashrc (linux / windows) or ~/.zshrc (mac)
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.bashrc
source ~/.bashrc
```

## Spawn in Action

In this 4 minute demo we:
- Create a new chain, customizing the modules and genesis
- Create a new `nameservice` module
  - Add the new message structure for transactions and queries
  - Store the new data types
  - Add the application logic
  - Connect it to the command line
- Build and launch a chain locally
- Interact with the chain's nameservice logic, settings a name, and retrieving it

[Follow Along with the nameservice demo in the docs](https://rollchains.github.io/spawn/v0.50/build/name-service/) | [source](./docs/versioned_docs/version-v0.50.x/02-build-your-application/01-nameservice.md)

https://github.com/rollchains/spawn/assets/31943163/ecc21ce4-c42c-4ff2-8e73-897c0ede27f0
