<div align="center">
  <h1>Spawn</h1>
</div>

Spawn is the easiest way to build, maintain and scale a Cosmos SDK blockchain. Spawn solves all the key pain points engineers face when building new Cosmos-SDK networks.
  - **Tailor-fit**: Pick and choose modules to create a network for your needs.
  - **Commonality**: Use native Cosmos tools and standards you're already familiar with.
  - **Integrations**: Github actions and end-to-end testing are configured right from the start.
  - **Iteration**: Quickly test between your new chain and established networks like the local Cosmos-Hub devnet.

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

https://github.com/rollchains/spawn/assets/31943163/ecc21ce4-c42c-4ff2-8e73-897c0ede27f0

## Requirements

- [`go 1.22+`](https://go.dev/doc/install)
- [`Docker`](https://docs.docker.com/get-docker/)
- [`git`](https://git-scm.com/)

## System Setup

* [MacOS](./docs/versioned_docs/version-v0.50.x/01-setup/01-system-setup.md#macos)
* [Ubuntu](./docs/versioned_docs/version-v0.50.x/01-setup/01-system-setup.md#linux-ubuntu)
* [Windows](./docs/versioned_docs/version-v0.50.x/01-setup/01-system-setup.md#windows)

## Install Spawn

```bash
git clone https://github.com/rollchains/spawn.git --depth 1 --branch v0.50.7
cd spawn

make install
make get-localic

# If you get "command 'spawn' not found", add to path
# Run the following in your terminal to test
# Then add to ~/.bashrc (linux) or ~/.zshrc (mac)
export PATH="$PATH:$(go env GOPATH)/bin"
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.bashrc
```

## IBC Demo

- [Getting Started Demo](./docs/versioned_docs/version-v0.50.x/03-ibc-chain-demo/01-demo.md)

## Add a Feature

If you wish to add a feature to spawn, reference the [NEW_FEATURE.md](./docs/dev/NEW_FEATURE.md) for guidance.
