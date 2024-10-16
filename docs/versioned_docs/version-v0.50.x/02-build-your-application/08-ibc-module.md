---
title: "IBC NameService Module"
sidebar_label: "IBC NameService (Part 2)"
sidebar_position: 8
slug: /build/name-service-ibc-module
---

# IBC Name Service Module

In this section, you will build on top of the Name Service tutorial to add cross chain functionality. This will allow you to sent a name from another network.

## Prerequisites
- [System Setup](../01-setup/01-system-setup.md)
- [Install Spawn](../01-setup/02-install-spawn.md)
- [Build Your Name Service Chain Tutorial](./01-nameservice.md)

## Create your chain

You should already have a network, `rollchain`, with the nameservice module from the nameservice tutorial. If you do not, complete that tutorial now.

:::note warning
Make sure you do not have the previous testnet still running by stopping it with: `killall -9 rolld`
:::

## Scaffold the IBC Module

```bash
# if you are not already in the chain directory:
cd rollchain

# scaffold the base IBC module for The
# cross chain name service
spawn module new nsibc --ibc-module

# compile latest code with matching module name
# failure to do this will result in: `panic: reflect: New(nil)`
make proto-gen
```

## Use the NameService Module

You now use the nameservice module you built previously within this new IBC module. This will allow you to save the name mapping on the name service, making it available for both IBC and native chain interactions.

```go title="x/nsibc/keeper/keeper.go"
import (
	...
	// highlight-next-line
	nameservicekeeper "github.com/rollchains/rollchain/x/nameservice/keeper"
)

type Keeper struct {
	...
	// highlight-next-line
	NameServiceKeeper *nameservicekeeper.Keeper
}
```

<details>
	<summary>Keeper Setup Image</summary>

	![View](https://github.com/user-attachments/assets/4dd3e50d-1528-4ae4-91a2-a27612bf69d7)
</details>


```go title="x/nsibc/keeper/keeper.go"
// NewKeeper creates a new Keeper instance.
func NewKeeper(
	...
	// highlight-next-line
	nsk *nameservicekeeper.Keeper, // placed on the last line of NewKeeper
) Keeper {
    ...

	k := Keeper{
		...
		// highlight-next-line
		NameServiceKeeper: nsk,
	}
```
<details>
	<summary>NewKeeper Image</summary>

	![View](https://github.com/user-attachments/assets/7639e468-a354-468d-8368-6bedd3c142a7)
</details>

## Provide NameService to the IBC Module

You must now give the IBC module access to nameservice keeper. It needs this reference so that the logic and connections can be shared. This is done in the `app/app.go` file. Find where the NameService IBC line is and update it to include the `&app.NameserviceKeeper,` reference.

You can find the `NameserviceKeeper` set just after the `NsibcKeeper` is set. If you would like to re-order the original NameService keeper, you can do so.

```go title="app/app.go"
	// Create the nsibc IBC Module Keeper
	app.NsibcKeeper = nsibckeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[nsibctypes.StoreKey]),
		app.IBCKeeper.ChannelKeeper,
		app.IBCKeeper.PortKeeper,
		scopedNsibc,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		// highlight-next-line
		&app.NameserviceKeeper, // placed on the last line of NewKeeper
	)
```

<details>
	<summary>Application NameService Reference Image</summary>

	![View](https://github.com/user-attachments/assets/6da58e1d-481b-46ba-bb66-d6c4a96971d0)
</details>



## Set Name on IBC Packet

Now that the IBC module has access to the nameservice, you can add the logic to set a name received from another chain (called the counterparty). To implement, the `OnRecvPacket` method has a placeholder for where the logic should run called `handleOnRecvLogic`. Find the `OnRecvPacket` in the ibc_module.go file, then find where the `handleOnRecvLogic` method resides.

```go title="x/nsibc/ibc_module.go"
// Find this method in the file
func (im ExampleIBCModule) handleOnRecvLogic(ctx context.Context, data types.ExamplePacketData) error {
	...
	return nil
}
```

<details>
	<summary>handleOnRecvLogic location</summary>

	![View](https://github.com/user-attachments/assets/011cb6cb-6664-47b9-a09e-fe1b62862987)
</details>



Once found, remove the lines within and replace with the following return.

```go title="x/nsibc/ibc_module.go"
func (im ExampleIBCModule) handleOnRecvLogic(ctx context.Context, data types.ExamplePacketData) error {
	// highlight-start
	if len(data.SomeData) > 32 {
		return fmt.Errorf("name cannot be longer than 32 characters")
	}
	return im.keeper.NameServiceKeeper.NameMapping.Set(ctx, data.Sender, data.SomeData)
	// highlight-end
}
```

This sets the name mapping from the sender to some data (the name) in the original nameservice module.

:::note
This is for example to show cross module interaction / extension with IBC.
You could just as easily write the NameMapping in the ibc keeper store as well.
:::

## Start Testnet

```bash
# build chain binary
make install

# verify the binary works. if you get a panic,
# `make proto-gen`, then re `make install`
rolld

# Build docker image, set configs, keys, and install binary
#
# Error 1 (ignored) codes are okay here if you already have
# the keys and configs setup. If so you only have to `make local-image`
# in future runs :)
make setup-testnet

# run testnet between itself and an IBC relayer
# this will take a minute
local-ic start self-ibc
```

## Import Testnet Helpers

Pasting the following lines in your terminal will import helper functions to interact with the testnet.
The source is publicly available on GitHub to review. It gives you the ability to interact with the testnet easily using special `ICT_` commands.

:::note
Local Interchain comes with a bash client helper for interaction as well. Check it out if you love bash

```bash
curl -s https://raw.githubusercontent.com/strangelove-ventures/interchaintest/main/local-interchain/bash/source.bash > ict_source.bash
source ./ict_source.bash
```

Further `local-ic interact` sub commands can be found in the [interchaintest/local-interchain documentation](https://github.com/strangelove-ventures/interchaintest/tree/main/local-interchain).
:::

## Connect Your IBC Modules

You are ready to connect the two chains with your IBC module protocol. The [cosmos/relayer](https://github.com/cosmos/relayer) is already running between the 2 networks now.

:::note
A Channel is a connection between two chains, like a highway. A port is a specific protocol (or logic) that can connect with itself on another chain.
For example; transfer to transfer, nsibc to nsibc, but transfer to nsibc can not be done. The version is just extra information (metadata) about the connection.

These values are found in the keys.go file as the module name. By default version is just the module name + "-1".
:::

Execute the command on the testnet to connect the two chains with the IBC module with the relayer.

```bash
# This will take a minute as it connects multiple
# networks together via the relayer process
local-ic interact localchain-1 relayer-exec "rly tx connect localchain-1_localchain-2 --src-port=nsibc --dst-port=nsibc --order=unordered --version=nsibc-1"

# Running the channels command now shows 2 channels, one for `transfer`
# and 1 for `nsibc`, marked as channel-1.
local-ic interact localchain-1 get_channels
```

## Submit Name Service Name Over IBC
```bash
# Set the IBC name from chain 1.
# view this command in x/nsibc/client/tx.go
rolld tx nsibc example-tx nsibc channel-1 testname --from acc0 --chain-id localchain-1 --yes

# View the logs
rolld q tx DCDBFA80AE95F7C010296623CDB89E1A537CCC6A279A952B5C66C82DA15652B0

# Verify chain 2 set the name (
# `rolld keys show -a acc0` from chain-1
local-ic interact localchain-2 query "nameservice resolve roll1hj5fveer5cjtn4wd6wstzugjfdxzl0xpg2te87"
```

## Summary

You just build an IBC module that interacts with your other nameservice module! It allowed you to set your name from a different network entirely and securely with IBC.

## What you Learned

* Scaffolding an IBC module
* Importing another module
* Adding business logic for an IBC request
* Connecting two chains with a custom IBC protocol
* Sending your first IBC packet from chain A
* Processing the packet on chain B and verifying it was set
