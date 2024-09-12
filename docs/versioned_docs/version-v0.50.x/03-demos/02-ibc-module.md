---
title: "IBC Module"
sidebar_label: "IBC NameService Module"
sidebar_position: 1
slug: /demo/ibc-module
---

# IBC NameService Module

In this tutorial, you will build on the [nameservice tutorial](../02-build-your-chain/01-nameservice.md) to add cross chain functionality. This will allow you to sent a name from another network.

## Prerequisites
- [System Setup](../01-setup/01-system-setup.md)
- [Install Spawn](../01-setup/02-install-spawn.md)
- [Build Your Name Service Chain Turotial](../02-build-your-chain/01-nameservice.md)

## Create your chain

You should already have a network, `rollchain`, with the nameservice module from the [nameservice tutorial](../02-build-your-chain/01-nameservice.md). If you do not, complete that tutorial now.

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
```

## Import the other NameService Module

You now reference the nameservice module you built within this new IBC module. This will allow you to save the name mapping on the name service, making it available for both IBC and native chain interactions.

```go title="x/nsibc/keeper/keeper.go"
import (
	...
	nameservicekeeper "github.com/rollchains/rollchain/x/nameservice/keeper"
)

type Keeper struct {
	...
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
	nsk *nameservicekeeper.Keeper,
) Keeper {
    ...

	k := Keeper{
		...
		NameServiceKeeper: nsk,
	}
```
<details>
	<summary>NewKeeper Image</summary>

	![View](https://github.com/user-attachments/assets/7639e468-a354-468d-8368-6bedd3c142a7)
</details>

## Provide NameService to the IBC Module

You must now give the IBC module access to nameservice keeper. It needs this reference so that the logic and connections can be shared. This is done in the `app/app.go` file. Find where the EvidenceKeeper line is, and overwrite the current lines with the following.

This moves the `NameserviceKeeper` above the IBC keeper (since you need to set that first so you can use it), and adds the `&app.NameserviceKeeper,` to the IBC Name Service keeper.

```go title="app/app.go"
	// If evidence needs to be handled for the app, set routes in router here and seal
	app.EvidenceKeeper = *evidenceKeeper

	// Create the nameservice Keeper
	app.NameserviceKeeper = nameservicekeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[nameservicetypes.StoreKey]),
		logger,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	// Create the nsibc IBC Module Keeper
	app.NsibcKeeper = nsibckeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[nsibctypes.StoreKey]),
		app.IBCKeeper.ChannelKeeper,
		app.IBCKeeper.PortKeeper,
		scopedNsibc,
		&app.NameserviceKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
```

<details>
	<summary>Application NameService Reference Image</summary>

	![View](https://github.com/user-attachments/assets/af456634-d7b7-475f-b468-7c14411803da)
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
	return im.keeper.NameServiceKeeper.NameMapping.Set(ctx, data.Sender, data.SomeData)
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

# verify the binary works
rolld

# build docker image
make local-image

# run testnet between itself and an IBC relayer
# this will take a minute
local-ic start self-ibc
```

## Import Testnet Helpers

Pasting the following lines in your terminal will import helper functions to interact with the testnet.
The source is publicly available on GitHub to review.

```bash
# Import the testnet interaction helper functions
# for local-interchain
source <(curl -s https://raw.githubusercontent.com/strangelove-ventures/interchaintest/main/local-interchain/bash/source.bash)
API_ADDR="http://localhost:8080"

# Waits for the testnet to start
ICT_POLL_FOR_START $API_ADDR 50 && echo "Testnet started"
```

## Connect Your IBC Modules

You are ready to connect the two chains with your IBC module protocol. The [cosmos/relayer](https://github.com/cosmos/relayer) is already running between the 2 networks now. You will just send a command to connect the new logic.

:::note
A Channel is a connection between two chains, like a highway. A port is a specific protocol (or logic) that can connect with itself on another chain.
For example; transfer to transfer, nsibc to nsibc, but transfer to nsibc can not be done. The version is just extra information (metadata) about the connection.

These values are found in the keys.go file as the module name. By default version is just the module name + "-1".
:::


```bash
# This will take a minute.
ICT_RELAYER_EXEC $API_ADDR "localchain-1" "rly tx connect localchain-1_localchain-2 --src-port=nsibc --dst-port=nsibc --order=unordered --version=nsibc-1"

# Running the channels command now shows 2 channels, one for `transfer`
# and 1 for `nsibc`, marked as channel-1.
echo `ICT_RELAYER_CHANNELS $API_ADDR "localchain-1"`
```

## Submit Name Service Name Over IBC
```bash
# Set the IBC name from chain 1.
# view this command in x/nsibc/client/tx.go
rolld tx nsibc example-tx nsibc channel-1 testname --from acc0 --chain-id localchain-1 --yes

# View the logs
rolld q tx 8A2009667022BE432B60158498C2256AEED0E86E9DFF79BD11CC9EA70DEC4A8A

# Verify chain 2 set the name (
# `rolld keys show -a acc0` from chain-1
ICT_QUERY "http://localhost:8080" "localchain-2" "nameservice resolve roll1hj5fveer5cjtn4wd6wstzugjfdxzl0xpg2te87"
```
