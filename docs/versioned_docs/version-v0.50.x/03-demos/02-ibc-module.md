---
title: "IBC Module"
sidebar_label: "IBC Module"
sidebar_position: 1
slug: /demo/ibc-module
---

# IBC Module

In this tutorial, we'll build on the nameservice tutorial and add an IBC module to the chain. This will allow us to sent our name on another network.

## Prerequisites
- [System Setup](../01-setup/01-system-setup.md)
- [Install Spawn](../01-setup/02-install-spawn.md)
- [Build Your Name Service Chain](../02-build-your-chain/01-nameservice.md)

## Create your chain

You should already have a network called `rollchain` with the nameservice module enabled from the [nameservice tutorial](../02-build-your-chain/01-nameservice.md). If you do not, complete that tutorial now.

## Scaffold the IBC Module

```bash
# if not already in the rollchain directory
cd rollchain

# scaffolds your new nameService IBC module
spawn module new nsibc --ibc-module
```

## Add reference to the previous Name Service keeper

```go title="x/nsibc/keeper/keeper.go"
import (
	...
	nameservicekeeper "github.com/rollchains/rollchain/x/nameservice/keeper"
)

// Keeper defines the module keeper.
type Keeper struct {
	...
	NameServiceKeeper *nameservicekeeper.Keeper
}
```
![View](https://github.com/user-attachments/assets/4dd3e50d-1528-4ae4-91a2-a27612bf69d7)

```go title="x/nsibc/keeper/keeper.go"
// NewKeeper creates a new swap Keeper instance.
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
![View](https://github.com/user-attachments/assets/7639e468-a354-468d-8368-6bedd3c142a7)

## Update it in the application

Ensure that `app.NameserviceKeeper` comes before the `app.NsibcKeeper` in the `app/app.go` file. Then fix the `nsibckeeper.NewKeeper` function to use the nameservice keeper logic.

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
![View](https://github.com/user-attachments/assets/af456634-d7b7-475f-b468-7c14411803da)

## Use NameService Logic

The `OnRecvPacket` method in this file has a placeholder for where your logic should run. Find the `OnRecvPacket` in the ibc_module.go file, then find where the `handleOnRecvLogic` method resides

```go title="x/nsibc/ibc_module.go"
// Find this method in the file
func (im ExampleIBCModule) handleOnRecvLogic(ctx context.Context, data types.ExamplePacketData) error {
	...
	return nil
}
```
![View](https://github.com/user-attachments/assets/011cb6cb-6664-47b9-a09e-fe1b62862987)

Once found, remove the lines within it and replace it with the following.

```go title="x/nsibc/ibc_module.go"
func (im ExampleIBCModule) handleOnRecvLogic(ctx context.Context, data types.ExamplePacketData) error {
	return im.keeper.NameServiceKeeper.NameMapping.Set(ctx, data.Sender, data.SomeData)
}
```

This will set the name mapping, from the sender to some data (the name) in the original nameservice module.

::note
This is for example to show cross module interaction / extension with IBC.
You could just as easily write the NameMapping in the ibc keeper store as well.
:::

## Running testnet

```bash
# build chain binary
make install

# verify the binary works
rolld

# build docker image
make local-image

# run testnet between itself
local-ic start self-ibc
```

# Connect New IBC Module

```bash
source <(curl -s https://raw.githubusercontent.com/strangelove-ventures/interchaintest/main/local-interchain/bash/source.bash)
API_ADDR="http://localhost:8080"

ICT_POLL_FOR_START $API_ADDR 50 && echo "Testnet started"

# only 1 channel (ics-20) is auto created on start of the testnet
echo `ICT_RELAYER_CHANNELS $API_ADDR "localchain-1"`

# We will then create a new channel between localchain-1 and localchain-2
# Using the new nameservice module we just created.
ICT_RELAYER_EXEC $API_ADDR "localchain-1" "rly tx connect localchain-1_localchain-2 --src-port=nsibc --dst-port=nsibc --order=unordered --version=nsibc-1"

# We can then run the channels command again to verify the new channel was createds
echo `ICT_RELAYER_CHANNELS $API_ADDR "localchain-1"`
```

## Really Interaction
```bash
# Set the IBC name from chain 1.
rolld tx nsibc example-tx nsibc channel-1 testname --from acc0 --chain-id localchain-1 --yes

# View the logs
rolld q tx 8A2009667022BE432B60158498C2256AEED0E86E9DFF79BD11CC9EA70DEC4A8A

# Verify chain 2 set the name (
# rolld keys show -a acc0 from chain-1
ICT_QUERY "http://localhost:8080" "localchain-2" "nameservice resolve roll1hj5fveer5cjtn4wd6wstzugjfdxzl0xpg2te87"
```
