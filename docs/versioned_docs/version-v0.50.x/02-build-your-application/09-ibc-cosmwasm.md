---
title: "IBC NameService Contract"
sidebar_label: "IBC Contract (Part 3)"
sidebar_position: 8
slug: /build/name-service-ibc-contract
---

# IBC Name Service Contract

You will build a new IBC contract with [CosmWasm](https://cosmwasm.com), enabling the same features we just built out in the IBC module. While this is a part 3 of the series, it can be done standalone as it requires a new chain. It is a similar concept to the previous parts 1 and 2, but with a smart contract focus instead of a chain.

## Prerequisites
- [System Setup](../01-setup/01-system-setup.md)
- [Install Spawn](../01-setup/02-install-spawn.md)
- [Rust + CosmWasm](../01-setup/01-system-setup.md)

## Setup the Chain

Build a new blockchain with CosmWasm enabled.

```bash
GITHUB_USERNAME=rollchains

spawn new cwchain \
--consensus=proof-of-stake \
--bech32=roll \
--denom=uroll \
--bin=rolld \
--disabled=block-explorer \
--org=${GITHUB_USERNAME}

# move into the chain directory
cd cwchain

# download latest dependencies
go mod tidy
```

## Build CosmWasm Contract

CosmWasm has a template repository that is used to generate new contracts. A minimal contract will be built with the `nameservice-contract` name provided on a [specific commit](https://github.com/CosmWasm/cw-template/commits/a2a169164324aa1b48ab76dd630f75f504e41d99/).

```bash
cargo generate --git https://github.com/CosmWasm/cw-template.git \
    --name nameservice-contract \
    --force-git-init \
    -d minimal=true --tag a2a169164324aa1b48ab76dd630f75f504e41d99
```

Open the contract in your code editor now to begin adding the application logic.

:::note Info
It is useful to install code rust extensions like [rust-analyzer](https://marketplace.visualstudio.com/items?itemName=rust-lang.rust-analyzer) and [even better toml](https://marketplace.visualstudio.com/items?itemName=tamasfe.even-better-toml) for an increased editing experience.
:::

```bash
# open using vscode in the terminal
code nameservice-contract/
```

## Update Contract Dependencies

This version of the CosmWasm template has some outdated versions. Update these in the `Cargo.toml` file and add the "ibc3" capability (for IBC support).

```toml title="Cargo.toml"
[dependencies]
# highlight-start
cosmwasm-schema = "1.5.7"
cosmwasm-std = { version = "1.5.7", features = [
    # "cosmwasm_1_3",
    "ibc3"
] }
# highlight-end
```

Update your local environment with the dependencies.

```bash
cargo update
```

## Setup State

This Rust code defines the structure for a name service in a CosmWasm smart contract. It saves a map of all channels (outside chain connections) to a list of wallet address and their associated names.

```rust title="src/state.rs"
// highlight-start
use std::collections::BTreeMap;

use cw_storage_plus::Map;

// Pair the wallet address to the name a user provides.
pub type WalletMapping = BTreeMap<String, String>;

/// create a new empty wallet mapping for a channel.
/// useful if a channel is opened and we have no data yet
pub fn new_wallet_mapping() -> WalletMapping {
    BTreeMap::new()
}

/// Name Service maps for each channel saved to a storage object
pub const NAME_SERVICE: Map<String, WalletMapping> = Map::new("nameservice");
// highlight-end
```

## Setup Interactions

Now that the state is setup, focus on modeling the users interaction with the contract. Users should be able to set a name. This also requires an input for a "channel" since a contract could connect to multiple chains. It could be written in a way that a user could set it to all channels, but for simplicity, we will require a channel to be specified. Just as is set, a user should get the name with the same format: a channel and a wallet address. Then a new message type is added specifically for IBCExecution messages. This is the packet transfered over the network, between chains, and gives the ability to set a name elsewhere on its contract.

The contract will call the IBCExecuteMsg when a user runs the ExecuteMsg.SendName function. This indirectly generates the packet and submits it for the user.

```rust title="src/msg.rs"
use cosmwasm_schema::{cw_serde, QueryResponses};

#[cw_serde]
pub struct InstantiateMsg {}

#[cw_serde]
pub enum ExecuteMsg {
    // highlight-next-line
    SetName { channel: String, name: String },
}

#[cw_serde]
#[derive(QueryResponses)]
pub enum QueryMsg {
    // highlight-start
    #[returns(GetNameResponse)]
    GetName { channel: String, wallet: String },
    // highlight-end
}

// highlight-start
#[cw_serde]
pub enum IbcExecuteMsg {
    SetName { wallet: String, name: String },
}

#[cw_serde]
pub struct GetNameResponse {
    pub name: String,
}
// highlight-end
```


## Contract Logic

Here are all the imports used in this. Replace your files top.

```rust title="src/contract.rs"
// highlight-start
#[cfg(not(feature = "library"))]
use cosmwasm_std::entry_point;
use cosmwasm_std::{Binary, Deps, DepsMut, Env, MessageInfo, Response, StdResult};

use crate::error::ContractError;
use crate::msg::{ExecuteMsg, IbcExecuteMsg, InstantiateMsg, QueryMsg};

use cosmwasm_std::{to_json_binary, IbcMsg, IbcTimeout, StdError};
// highlight-end
```

Instantiate creates a new version of this contract that you control. Rather than being unimplemented, return a basic response saying it was Ok (successful) and add some extra logging metadata.


```rust title="src/contract.rs"
#[cfg_attr(not(feature = "library"), entry_point)]
pub fn instantiate(
    _deps: DepsMut,
    _env: Env,
    _info: MessageInfo,
    _msg: InstantiateMsg,
) -> Result<Response, ContractError> {
    // highlight-next-line
    Ok(Response::new().add_attribute("method", "instantiate"))
}
```

The ExecuteMsg::SetName method is allowed to be interacted from anyone. Just like instantiate we return a new Ok response. This time an add_message function is added. This will generate the packet as the user interacts, performing a new action from a previous action. This uses the IbcMsg::SendPacket built in type to create it for the user. Notice the data field includes the IbcExecuteMsg::SetName we defined before. This is transferred to the other version of this contract on another chain and processed.

If the packet is not picked up by a relayer service provider within a few minutes, the packet will become void and stop attempting execution on the other chain's contract.

```rust title="src/contract.rs"
#[cfg_attr(not(feature = "library"), entry_point)]
pub fn execute(
    _deps: DepsMut,
    // highlight-start
    env: Env, // removes the underscore _
    info: MessageInfo,
    msg: ExecuteMsg,
    // highlight-end
) -> Result<Response, ContractError> {
    // highlight-start
    match msg {
        ExecuteMsg::SetName { channel, name } => {
            Ok(Response::new()
                .add_attribute("method", "set_name")
                .add_attribute("channel", channel.clone())
                // outbound IBC message, where packet is then received on other chain
                .add_message(IbcMsg::SendPacket {
                    channel_id: channel,
                    data: to_json_binary(&IbcExecuteMsg::SetName {
                        name: name,
                        wallet: info.sender.into_string(),
                    })?,
                    // default timeout of two minutes.
                    timeout: IbcTimeout::with_timestamp(env.block.time.plus_seconds(120)),
                }))
        }
    }
    // highlight-end
}
```

The users name is not set, but it is only useful if you can also get said data. Read from the `NAME_SERVICE` storage Map defined in `state.rs`. Using may load grabs the data if the channel has a name set. If no channel is found (no users have set a name from this chain), it returns an error to the user requesting. If a channel of pairs is found, it loads them and checks if the wallet address requested is set in it. If it is, return what the wallets name is set to. If the user with this wallet did not set a name, return an error.

```rust title="src/contract.rs"
#[cfg_attr(not(feature = "library"), entry_point)]
pub fn query(deps: Deps, _env: Env, msg: QueryMsg) -> StdResult<Binary> {
    // highlight-start
    match msg {
        QueryMsg::GetName { channel, wallet } => {
            crate::state::NAME_SERVICE
                .may_load(deps.storage, channel.clone())
                .and_then(|maybe_wallets| match maybe_wallets {
                    Some(wallets) => match wallets.get(&wallet) {
                        Some(wallet) => Ok(to_json_binary(&crate::msg::GetNameResponse {
                            name: wallet.clone(),
                        })?),
                        None => Err(StdError::generic_err("No name set for wallet")),
                    },
                    None => Err(StdError::generic_err("Channel not found")),
                })
        }
    }
    // highlight-end
}
```

The contract will receive a packet and must run logic to process it. This is called the `try_set_name` method. It updates a given channel to include a new wallet. If the wallet already exists, it will overwrite the name. It then returns the users name back, or an error, for our future IBC logic to handle.

```rust title="src/contract.rs"
// highlight-start
/// called on IBC packet receive in other chain
pub fn try_set_name(
    deps: DepsMut,
    channel: String,
    wallet: String,
    name: String,
) -> Result<String, StdError> {
    crate::state::NAME_SERVICE.update(deps.storage, channel, |wallets| -> StdResult<_> {
        let mut wallets = wallets.unwrap_or_default();
        wallets.insert(wallet, name.clone());
        Ok(wallets)
    })?;
    Ok(name)
}
// highlight-end
```

## Create Transaction acknowledgement

Create a new file, `ack.rs`, to handle the IBC ACK (acknowledgement) messages. This just returns back to the user if their interaction was a success or an error.

```bash
touch src/ack.rs
```

```rust title="src/ack.rs"
// highlight-start
use cosmwasm_schema::cw_serde;
use cosmwasm_std::{to_json_binary, Binary};

/// IBC ACK. See:
/// https://github.com/cosmos/cosmos-sdk/blob/f999b1ff05a4db4a338a855713864497bedd4396/proto/ibc/core/channel/v1/channel.proto#L141-L147
#[cw_serde]
pub enum Ack {
    Result(Binary),
    Error(String),
}

pub fn make_ack_success() -> Binary {
    let res = Ack::Result(b"1".into());
    to_json_binary(&res).unwrap()
}

pub fn make_ack_fail(err: String) -> Binary {
    let res = Ack::Error(err);
    to_json_binary(&res).unwrap()
}
// highlight-end
```

:::note Note
Rust has a lib.rs file that is the entry point for the Rust library. All files that are used must be mentioned here to have access to them.
:::

Add the ack logic to the lib.rs so the application can use it.

```rust title="src/lib.rs"
// highlight-next-line
pub mod ack;
pub mod contract;
...

pub use crate::error::ContractError;
```

## Setup Errors

If a relayer or contract try to connect to an unlike protocol, the InvalidVersion error will be returned to the attempted actor. This contract only supports 1 protocol version across networks because it must speak the same "language". If you speak english while another person speaks spanish, your interactions are incompatible. Contracts are like this too. They verify their protocol version in a format like "ics-20" or "ns-1" first to make sure they can communicate.

OrderedChannel is a type of flow control for network packets, or interactions. This tutorial uses unordered paths so any packet that times out or fails does not block future packets from going through. If a relayer tries to make this an ordered path, the contract returns this error to stop them from doing so.

```rust title="src/error.rs"
use cosmwasm_std::StdError;
use thiserror::Error;

#[derive(Error, Debug)]
pub enum ContractError {
    #[error("{0}")]
    Std(#[from] StdError),

    #[error("Unauthorized")]
    Unauthorized {},
    // Add any other custom errors you like here.
    // Look at https://docs.rs/thiserror/1.0.21/thiserror/ for details.

    // highlight-start
    #[error("invalid IBC channel version. Got ({actual}), expected ({expected})")]
    InvalidVersion { actual: String, expected: String },

    #[error("only unordered channels are supported")]
    OrderedChannel {},
    // highlight-end
}

// highlight-start
// There is an IBC specific error that is never returned.
#[derive(Error, Debug)]
pub enum Never {}
// highlight-end
```


## IBC Specific Logic

Create a new file `ibc.rs`. Add this to the lib.rs. This is where our core IBC logic will go.

```bash
touch src/ibc.rs
```

```rust title="src/lib.rs"
// highlight-next-line
pub mod ibc;
pub mod ack;
pub mod contract;
mod error;
pub mod helpers;
pub mod msg;
pub mod state;

pub use crate::error::ContractError;
```

Place the following in the ibc.rs file. Import all the types needed, set the IBC version to "ns-1" to stand for "nameservice-1", and add the basic validation logic for the contract. You must ensure contracts that try to talk together are verified to work together. This is that logic.

```rust title="src/ibc.rs"
// highlight-start
#[cfg(not(feature = "library"))]
use cosmwasm_std::entry_point;
use cosmwasm_std::{
    from_json, DepsMut, Env, IbcBasicResponse, IbcChannel, IbcChannelCloseMsg,
    IbcChannelConnectMsg, IbcChannelOpenMsg, IbcChannelOpenResponse, IbcOrder, IbcPacketAckMsg,
    IbcPacketReceiveMsg, IbcPacketTimeoutMsg, IbcReceiveResponse,
};

use crate::{
    ack::{make_ack_fail, make_ack_success},
    contract::try_set_name,
    msg::IbcExecuteMsg,
    state::NAME_SERVICE,
    ContractError,
};

pub const IBC_VERSION: &str = "ns-1";

pub fn validate_order_and_version(
    channel: &IbcChannel,
    counterparty_version: Option<&str>,
) -> Result<(), ContractError> {
    // We expect an unordered channel here. Ordered channels have the
    // property that if a message is lost the entire channel will stop
    // working until you start it again.
    if channel.order != IbcOrder::Unordered {
        return Err(ContractError::OrderedChannel {});
    }

    if channel.version != IBC_VERSION {
        return Err(ContractError::InvalidVersion {
            actual: channel.version.to_string(),
            expected: IBC_VERSION.to_string(),
        });
    }

    // Make sure that we're talking with a counterparty who speaks the
    // same "protocol" as us.
    //
    // For a connection between chain A and chain B being established
    // by chain A, chain B knows counterparty information during
    // `OpenTry` and chain A knows counterparty information during
    // `OpenAck`. We verify it when we have it but when we don't it's
    // alright.
    if let Some(counterparty_version) = counterparty_version {
        if counterparty_version != IBC_VERSION {
            return Err(ContractError::InvalidVersion {
                actual: counterparty_version.to_string(),
                expected: IBC_VERSION.to_string(),
            });
        }
    }

    Ok(())
}
// highlight-end
```

The contract verifies data on an attempted open of a new connection. Ensure the contracts talk the same protocol language, and that all the validation basic logic is connect. Then when a channel is closed, clear the data from storage for it. It is very rare you would want to close a channel.

```rust title="src/ibc.rs"
// highlight-start
/// Handles the `OpenInit` and `OpenTry` parts of the IBC handshake.
#[cfg_attr(not(feature = "library"), entry_point)]
pub fn ibc_channel_open(
    _deps: DepsMut,
    _env: Env,
    msg: IbcChannelOpenMsg,
) -> Result<IbcChannelOpenResponse, ContractError> {
    validate_order_and_version(msg.channel(), msg.counterparty_version())?;
    Ok(None)
}

#[cfg_attr(not(feature = "library"), entry_point)]
pub fn ibc_channel_close(
    deps: DepsMut,
    _env: Env,
    msg: IbcChannelCloseMsg,
) -> Result<IbcBasicResponse, ContractError> {
    let channel = msg.channel().endpoint.channel_id.clone();
    NAME_SERVICE.remove(deps.storage, channel.clone());
    Ok(IbcBasicResponse::new()
        .add_attribute("method", "ibc_channel_close")
        .add_attribute("channel", channel))
}
// highlight-end
```

When a successful connection is made, the contract saves a new blank wallet mapping to the channel's unique id. 'channel-0' is the first. All  future connections are channel-1, channel-2, etc. This is the first step in the IBC process. The contract is now ready to receive packets once the handler logic is put in place on receive.

```rust title="src/ibc.rs"
// highlight-start
#[cfg_attr(not(feature = "library"), entry_point)]
pub fn ibc_channel_connect(
    deps: DepsMut,
    _env: Env,
    msg: IbcChannelConnectMsg,
) -> Result<IbcBasicResponse, ContractError> {
    validate_order_and_version(msg.channel(), msg.counterparty_version())?;

    let channel = msg.channel().endpoint.channel_id.clone();
    NAME_SERVICE.save(
        deps.storage, channel.clone(), &crate::state::new_wallet_mapping(),
    )?;

    Ok(IbcBasicResponse::new()
        .add_attribute("method", "ibc_channel_connect")
        .add_attribute("channel_id", channel))
}
// highlight-end
```

ibc_packet_receive handles incoming packets from already connected networks. The packet is forwarded to this contract and processed in `do_ibc_packet_receive`. It takes the channel and the packet data *(the IbcMsg::SetName sent out from the ExecuteMsg earlier)*, and tries to set the name on a wallet for this channel. If successful, it returns an acknowledgment of success. If not, it returns an acknowledgment of failure. The user will see this in their log event output.

```rust title="src/ibc.rs"
// highlight-start
#[cfg_attr(not(feature = "library"), entry_point)]
pub fn ibc_packet_receive(
    deps: DepsMut,
    _env: Env,
    msg: IbcPacketReceiveMsg,
) -> Result<IbcReceiveResponse, crate::error::Never> {
    // Regardless of if our processing of this packet works we need to
    // commit an ACK to the chain. As such, we wrap all handling logic
    // in a septate function and on error write out an error ack.
    match do_ibc_packet_receive(deps, msg) {
        Ok(response) => Ok(response),
        Err(error) => Ok(IbcReceiveResponse::new()
            .add_attribute("method", "ibc_packet_receive")
            .add_attribute("error", error.to_string())
            .set_ack(make_ack_fail(error.to_string()))),
    }
}

pub fn do_ibc_packet_receive(
    deps: DepsMut,
    msg: IbcPacketReceiveMsg,
) -> Result<IbcReceiveResponse, ContractError> {
    // The channel this packet is being relayed along on this chain.
    let channel = msg.packet.dest.channel_id;
    let msg: IbcExecuteMsg = from_json(&msg.packet.data)?;

    match msg {
        IbcExecuteMsg::SetName { wallet, name } => {
            let name = try_set_name(deps, channel, wallet, name)?;

            Ok(IbcReceiveResponse::new()
                .add_attribute("method", "execute_increment")
                .add_attribute("name", name)
                .set_ack(make_ack_success()))
        }
    }
}
// highlight-end
```

Sometimes after a failed acknowledgement the contract may want to rollback some data or make note of it for future reference. This contract is simple enough so no rollback or refunds are required. We just return a basic response to the user for both the ack or a timeout. Think of this similarly as a NoOp (no operation).

```rust title="src/ibc.rs"
// highlight-start
#[cfg_attr(not(feature = "library"), entry_point)]
pub fn ibc_packet_ack(
    _deps: DepsMut,
    _env: Env,
    _ack: IbcPacketAckMsg,
) -> Result<IbcBasicResponse, ContractError> {
    Ok(IbcBasicResponse::new().add_attribute("method", "ibc_packet_ack"))
}

#[cfg_attr(not(feature = "library"), entry_point)]
pub fn ibc_packet_timeout(
    _deps: DepsMut,
    _env: Env,
    _msg: IbcPacketTimeoutMsg,
) -> Result<IbcBasicResponse, ContractError> {
    Ok(IbcBasicResponse::new().add_attribute("method", "ibc_packet_timeout"))
}
// highlight-end
```

## Build Contract From Source

The contract can now be compiled from its source into the .wasm file. This is the binary executable that will be uploaded to the chain.

```bash
cargo-run-script optimize
```

---


### Start the chains and connect

Make sure you are in the `cwchain` directory to begin interacting and uploading the contract to the chain. It is time to start the cosmwasm chain and launch a testnet that connects to itself. The `self-ibc` chain is automatically generated for you on the creation with spawn. It launches 2 of your networks, localchain-1 and localchain-2, and connects them with a relayer operator at startup.

```bash
# move back into the cwchain directory
cd ../

# Install heighliner if you have not already
# (Easily docker builder tool)
make get-heighliner

# Build docker image, set configs, keys, and install binary
#
# Error 1 (ignored) codes are okay here if you already have
# the keys and configs setup. If so you only have to `make local-image`
# in future runs :)
make setup-testnet

local-ic start self-ibc
```

### Store the Contract on both chains

Get the [RPC ](https://www.techtarget.com/searchapparchitecture/definition/Remote-Procedure-Call-RPC) interaction addresses for each network from the local-interchain testnet API. Upload the contract source to both chains using the different RPC addresses.

```bash
RPC_1=`curl http://127.0.0.1:8080/info | jq -r .logs.chains[0].rpc_address`
RPC_2=`curl http://127.0.0.1:8080/info | jq -r .logs.chains[1].rpc_address`
echo "Using RPC_1=$RPC_1 and RPC_2=$RPC_2"

CONTRACT_SOURCE=./nameservice-contract/artifacts/nameservice_contract.wasm
rolld tx wasm store $CONTRACT_SOURCE --from=acc0 --gas=auto --gas-adjustment=2.0 --yes --node=$RPC_1
# rolld q wasm list-code --node=$RPC_1

rolld tx wasm store $CONTRACT_SOURCE --from=acc0 --gas=auto --gas-adjustment=2.0 --yes --node=$RPC_2 --chain-id=localchain-2
# rolld q wasm list-code --node=$RPC_2
```

### Instantiate our Contract on both chains

You can now create your contract from the source on each chain.

```bash
rolld tx wasm instantiate 1 '{}' --no-admin --from=acc0 --label="ns-1" --gas=auto --gas-adjustment=2.0 --yes --node=$RPC_1
rolld tx wasm instantiate 1 '{}' --no-admin --from=acc0 --label="ns-1" --gas=auto --gas-adjustment=2.0 --yes --node=$RPC_2 --chain-id=localchain-2

rolld q wasm list-contracts-by-creator roll1hj5fveer5cjtn4wd6wstzugjfdxzl0xpg2te87
rolld q wasm list-contracts-by-creator roll1hj5fveer5cjtn4wd6wstzugjfdxzl0xpg2te87 --node=$RPC_2

NSERVICE_CONTRACT=roll14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9sjczpjh
```

### Relayer connect

The relayer must now connect the contracts together and create an IBC connection, link, between them. Use the Local-Interchain helper methods to connect the contracts across the chains. This command will take a second then show a bunch of logs. `Error context canceled` is fine to see. You will verify they were opened in the next step.

```bash
# This will take a moment.
# Connects the 2 contracts together with the IBC relayer
local-ic interact localchain-1 relayer-exec "rly tx link localchain-1_localchain-2 --src-port wasm.${NSERVICE_CONTRACT} --dst-port=wasm.${NSERVICE_CONTRACT} --order unordered --version ns-1"
```


## Verify channels

Verify the channels were created. Query either with the application binary of the relayer itself. If you see both a channel-0 and channel-1 in your logs, it was a success. If you only see channel-0 re-run the above relayer exec tx link command.

```bash
# app binary
rolld q ibc channel channels

# relayer
local-ic interact localchain-1 relayer-exec "rly q channels localchain-1"
```

## Transaction interaction

Using the ExecuteMsg::SetName method, set a name. This will be transferred to chain 2 behind the scenes. Flushing the relayer will force it to auto pick up pending IBC packets and transfer them across. Not running this may take up to 30 seconds for the relayer to automatically pick it up.

```bash
# Set the name from chain 1
MESSAGE='{"set_name":{"channel":"channel-1","name":"myname"}}'
rolld tx wasm execute $NSERVICE_CONTRACT "$MESSAGE" --from=acc0 --gas=auto --gas-adjustment=2.0 --yes

# This will take a moment
# 'account sequence mismatch' errors are fine.
local-ic interact localchain-1 relayer-exec "rly tx flush"
```

---

### Verify data

After the packet is sent over the network, processed, and acknowledged *(something that can be done in less than 30 seconds)*, you can query the data on chain 2. You can also dump all the contract data out to get HEX and BASE64 encoded data for what the contract state storage looks like.

```bash
# query the name on chain 2, from chain 1
rolld q wasm state smart $NSERVICE_CONTRACT '{"get_name":{"channel":"channel-1","wallet":"roll1hj5fveer5cjtn4wd6wstzugjfdxzl0xpg2te87"}}' --node=$RPC_2

# dump contract state from the other chain
rolld q wasm state all $NSERVICE_CONTRACT --node=$RPC_2
```

## Summary

You just build an IBC protocol using cosmwasm! It allowed you to set a name on another network entirely and securely with IBC.

## What you Learned

* Scaffolding an CosmWasm contract
* Adding business logic for an IBC request
* Implementing IBC in a contract
* Connecting two CosmWasm contracts with a custom IBC protocol
* Sending a packet from contract A to contract B
