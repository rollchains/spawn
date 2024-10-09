---
title: "IBC NameService Contract"
sidebar_label: "IBC Contract (Part 3)"
sidebar_position: 8
slug: /build/name-service-ibc-contract
---

# IBC Name Service Contract

You will build a new IBC contract with [CosmWasm](https://cosmwasm.com), enabling the same features we just built out in the IBC module.

## Prerequisites
- [System Setup](../01-setup/01-system-setup.md)
- [Install Spawn](../01-setup/02-install-spawn.md)
- [Rust + CosmWasm](../01-setup/01-system-setup.md#cosmwasm)

## Setup the Chain

Build a new blockchain with CosmWasm enabled. THen generate a new module from the template for the reviews.

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

go mod tidy
```

Once the chain is started, continue on to the contract building steps

## CosmWasm Build Contract

CosmWasm has a template repository that is used to generate new contracts. A minimal contract will be built with the `nameservice-contract` name provided on a specific commit.

```bash
cargo generate --git https://github.com/CosmWasm/cw-template.git \
    --name nameservice-contract \
    --force-git-init \
    -d minimal=true --tag a2a169164324aa1b48ab76dd630f75f504e41d99
```

```bash
code nameservice-contract/
```

<!-- TODO: WIP here: https://github.com/Reecepbcups/nameservice-contract -->


<!-- TODO: modify cargo.toml -->
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

```bash
cargo update
```

## Setup State

```rust title="src/state.rs"
// highlight-start
use std::collections::BTreeMap;

use cw_storage_plus::Map;

// Wallet  -> Name
pub type WalletMapping = BTreeMap<String, String>;

pub fn new_wallet_mapping() -> WalletMapping {
    BTreeMap::new()
}

pub const NAME_SERVICE: Map<String, WalletMapping> = Map::new("nameservice");
// highlight-end
```

## Setup Interactions

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

Here are all the imports used in this. Replace your files top with these.

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

<!-- TODO: -->
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

<!-- TODO: -->
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

<!-- TODO: -->
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

<!-- TODO: -->
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

Create a new file, `ack.rs`, to handle the IBC ACK messages.

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

Add it to the lib.rs so the application can use it.

```rust title="src/lib.rs"
// highlight-next-line
pub mod ack;
pub mod contract;
...

pub use crate::error::ContractError;
```

## Setup Errors

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
    #[error("only unordered channels are supported")]
    OrderedChannel {},

    #[error("invalid IBC channel version. Got ({actual}), expected ({expected})")]
    InvalidVersion { actual: String, expected: String },
    // highlight-end
}

// highlight-start
#[derive(Error, Debug)]
pub enum Never {}
// highlight-end
```


## IBC Specific Logic

Create a new file `ibc.rs`. Add this to the lib.rs

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

Now begin working on the IBC logic.

<!-- TODO: input what the below does -->

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

<!-- TODO: -->

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

<!-- TODO: input what the below does -->

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

<!-- TODO: input what the below does -->

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

<!-- TODO: -->

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

```bash
cargo-run-script optimize
```


---


### Upload Contract to both networks

Make sure you are in the `cwchain` directory to begin interacting and uploading the contract to the chain.

```bash
# Build docker image, set configs, keys, and install binary
# Error 1 (ignored) codes are okay here if you already have
# the keys and configs setup. If so you only have to `make local-image`
# in future runs :)
make setup-testnet

local-ic start self-ibc
```


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

### Instantiate our Contract on oth chains

```bash
rolld tx wasm instantiate 1 '{}' --no-admin --from=acc0 --label="ns-1" --gas=auto --gas-adjustment=2.0 --yes --node=$RPC_1
rolld tx wasm instantiate 1 '{}' --no-admin --from=acc0 --label="ns-1" --gas=auto --gas-adjustment=2.0 --yes --node=$RPC_2 --chain-id=localchain-2

rolld q wasm list-contracts-by-creator roll1hj5fveer5cjtn4wd6wstzugjfdxzl0xpg2te87
rolld q wasm list-contracts-by-creator roll1hj5fveer5cjtn4wd6wstzugjfdxzl0xpg2te87 --node=$RPC_2

NSERVICE_CONTRACT=roll14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9sjczpjh
```

### Relayer connect

```bash
# Import the testnet interaction helper functions
# for local-interchain
curl -s https://raw.githubusercontent.com/strangelove-ventures/interchaintest/main/local-interchain/bash/source.bash > ict_source.bash
source ./ict_source.bash

API_ADDR="http://localhost:8080"

# This will take a moment.
ICT_RELAYER_EXEC "$API_ADDR" "localchain-1" \
    "rly tx link localchain-1_localchain-2 --src-port wasm.${NSERVICE_CONTRACT} --dst-port=wasm.${NSERVICE_CONTRACT} --order unordered --version ns-1"
```


## Verify channels

```bash
# app binary
rolld q ibc channel channels

# relayer
ICT_RELAYER_EXEC "$API_ADDR" "localchain-1" "rly q channels localchain-1"
```

## Transaction against

```bash
# Set the name from chain 1
MESSAGE='{"set_name":{"channel":"channel-1","name":"myname"}}'
rolld tx wasm execute $NSERVICE_CONTRACT "$MESSAGE" --from=acc0 --gas=auto --gas-adjustment=2.0 --yes

# This will take a moment
# 'account sequence mismatch' errors are fine.
ICT_RELAYER_EXEC "$API_ADDR" "localchain-1" "rly tx flush"
```

---

### Verify data

```bash
# query the name on chain 2, from chain 1
rolld q wasm state smart $NSERVICE_CONTRACT '{"get_name":{"channel":"channel-1","wallet":"roll1hj5fveer5cjtn4wd6wstzugjfdxzl0xpg2te87"}}' --node=$RPC_2

# dump contract state from the other chain
rolld q wasm state all $NSERVICE_CONTRACT --node=$RPC_2
```
