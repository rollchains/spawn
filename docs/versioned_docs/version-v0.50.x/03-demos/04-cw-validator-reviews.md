---
title: "CW Validator Reviews"
sidebar_label: "CW Validator Reviews"
slug: /demo/cw-validator-reviews
---

# CosmWasm Validator Reviews

<!-- TODO: walk through what we are to do -->
- New network with cosmwasm
- Write a contract which takes in all validators
- Write a module with an endblock that updates all validators every X blocks
- Prove this works and is set, set reviews


## Prerequisites
- [System Setup](../01-setup/01-system-setup.md)
- [Install Spawn](../01-setup/02-install-spawn.md)

## CosmWasm Setup

```bash
# Install rust - https://www.rust-lang.org/tools/install
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh

# Update if you have it
rustup update

# Install wasm target
rustup target add wasm32-unknown-unknown

# Install generate and run script features
cargo install cargo-generate --features vendored-openssl
cargo install cargo-run-script
```

## Setup Chain

```bash
GITHUB_USERNAME=rollchains

spawn new rollchain \
--consensus=proof-of-stake \
--bech32=roll \
--denom=uroll \
--bin=rolld \
--disabled=block-explorer \
--org=${GITHUB_USERNAME}


# move into the chain directory
cd rollchain

# Generate the Cosmos-SDK reviews module
spawn module new reviews

# build the proto to code
make proto-gen
```

We will come back to this later. Let's build the contract first.

## CosmWasm Build Contract

```bash
# Build the template
cargo generate --git https://github.com/CosmWasm/cw-template.git --name validator-reviews-contract -d minimal=true --tag a2a169164324aa1b48ab76dd630f75f504e41d99

# open validator-reviews-contract/ in your text editor
code validator-reviews-contract/
```


## Set State Structure

```rust title="src/state.rs"
use std::collections::BTreeMap;

use cosmwasm_schema::cw_serde;
use cw_storage_plus::{Item, Map};

#[cw_serde]
pub struct Validator {
    pub address: String,
    pub moniker: String,
}
pub const VALIDATORS: Item<Vec<Validator>> = Item::new("validators");

// user -> text
pub type Reviews = BTreeMap<String, String>;

// validator_address -> reviews
pub const REVIEWS: Map<&str, Reviews> = Map::new("reviews");
```

## Set Input Structure

```rust title="src/msg.rs"
use cosmwasm_schema::{cw_serde, QueryResponses};

#[cw_serde]
pub struct InstantiateMsg {}

#[cw_serde]
pub enum ExecuteMsg {
    // highlight-next-line
    WriteReview { val_addr: String, review: String },
}

#[cw_serde]
#[derive(QueryResponses)]
pub enum QueryMsg {
    // highlight-start
    #[returns(Vec<crate::state::Validator>)]
    Validators {},

    #[returns(crate::state::Reviews)]
    Reviews { address: String },
    // highlight-end
}

// highlight-start
#[cw_serde]
pub enum SudoMsg {
    SetValidators {
        validators: Vec<crate::state::Validator>,
    },
}
// highlight-end
```

## Add a new Error

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
    #[error("The validator is not found in this set")]
    NoValidatorFound {},
    // highlight-end
}
```

## Add imports at the top

You can also import them during, but it is much easier to just do so now.

```rust title="src/contract.rs"
// highlight-start
use crate::state::{Reviews, REVIEWS, VALIDATORS};
use cosmwasm_std::to_json_binary;
use std::collections::BTreeMap;
// highlight-end
```

## Modify Instantiate Message

```rust title="src/contract.rs"
#[cfg_attr(not(feature = "library"), entry_point)]
pub fn instantiate(
    // highlight-next-line
    deps: DepsMut, // removes the underscore
    _env: Env,
    _info: MessageInfo,
    _msg: InstantiateMsg,
) -> Result<Response, ContractError> {
    // highlight-start
    VALIDATORS.save(deps.storage, &Vec::new())?;
    Ok(Response::default())
    // highlight-end
}
```

## Add Execute Logic

```rust title="src/contract.rs"
#[cfg_attr(not(feature = "library"), entry_point)]
pub fn execute(
    // highlight-next-line
    deps: DepsMut, // removes the underscore
    _env: Env,
    _info: MessageInfo,
    // highlight-next-line
    msg: ExecuteMsg,  // removes the underscore
) -> Result<Response, ContractError> {
    // highlight-start
    match msg {
        ExecuteMsg::WriteReview { val_addr, review } => {
            let active_validators = VALIDATORS.load(deps.storage)?;
            if active_validators.iter().find(|v| v.address == val_addr).is_none() {
                return Err(ContractError::NoValidatorFound {});
            }

            // Get current validator reviews if any. If there are none, create a new empty review map.
            let mut all_revs: Reviews = match REVIEWS.may_load(deps.storage, &val_addr) {
                Ok(Some(rev)) => rev,
                _ => BTreeMap::new(),
            };

            // Set this users review for the validator.
            all_revs.insert(val_addr.clone(), review);

            // Save the updated reviews
            REVIEWS.save(deps.storage, &val_addr, &all_revs)?;
        }
    }

    Ok(Response::default())
    // highlight-end
}
```

## Add Queries

```rust title="src/contract.rs"
#[cfg_attr(not(feature = "library"), entry_point)]
// highlight-start
pub fn query(deps: Deps, _env: Env, msg: QueryMsg) -> StdResult<Binary> {
    match msg {
        QueryMsg::Validators {} => {
            let validators = VALIDATORS.load(deps.storage)?;
            Ok(to_json_binary(&validators)?)
        }
        QueryMsg::Reviews { address } => {
            let reviews = REVIEWS.load(deps.storage, &address)?;
            Ok(to_json_binary(&reviews)?)
        }
    }
// highlight-end
}
```

## Add New Sudo Message

```rust title="src/contract.rs"
// highlight-start
#[cfg_attr(not(feature = "library"), entry_point)]
pub fn sudo(deps: DepsMut, _env: Env, msg: crate::msg::SudoMsg) -> Result<Response, ContractError> {
    match msg {
        crate::msg::SudoMsg::SetValidators { validators } => {
            VALIDATORS.save(deps.storage, &validators)?;
            Ok(Response::new())
        }
    }
}
// highlight-end
```



## Build Contract

We now build the contract within docker. Make sure you have docker installed and running.

```bash
# run the build optimizer (from source -> contract wasm binary)
cargo run-script optimize

# Output: ./artifacts/validator_reviews_contract.wasm
```


## Modify the Module

## Setup the Keeper

```go title="x/reviews/keeper.go"
import (
	...

    // highlight-start
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
    // highlight-end
)

type Keeper struct
    ...

    // highlight-start
    WasmKeeper     *wasmkeeper.Keeper
	ContractKeeper wasmtypes.ContractOpsKeeper
	StakingKeeper  *stakingkeeper.Keeper
    // highlight-end
}

func NewKeeper(
    ...
    // highlight-start
    wasmKeeper *wasmkeeper.Keeper, // since wasm may not be created yet.
	stakingKeeper *stakingkeeper.Keeper,
    // highlight-end
    authority string,
) Keeper {
    ...

    k := Keeper{
        ...

        // highlight-start
        WasmKeeper:     wasmKeeper,
        ContractKeeper: wasmkeeper.NewDefaultPermissionKeeper(wasmKeeper),
        StakingKeeper:  stakingKeeper,
        // highlight-end

        authority: authority,
	}
}
```

## Fix keeper_test

Testing wasm takes significantly more setup for the keeper. For now, we will just add a blank reference here.

```go title="x/reviews/keeper/keeper_test.go"
	// Setup Keeper.
    // highlight-next-line
	f.k = keeper.NewKeeper(encCfg.Codec, storeService, logger, &wasmkeeper.Keeper{}, f.stakingKeeper, f.govModAddr)
```

## Dep Inject (v2)

Reference same warning above, cosmwasm does not have depinject support at the time of writing.

```go title="x/reviews/keeper/keeper_test.go"
func ProvideModule(in ModuleInputs) ModuleOutputs {
    ...

    // highlight-start
    k := keeper.NewKeeper(in.Cdc, in.StoreService, log.NewLogger(os.Stderr), nil, &in.StakingKeeper, govAddr)
}
```

## Fix app.go references

```go title="app/app.go"
	app.ReviewsKeeper = reviewskeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[reviewstypes.StoreKey]),
		logger,
        // highlight-start
		&app.WasmKeeper,
		app.StakingKeeper,
        // highlight-end
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
```


## Module

```go title="x/reviews/module.go"
// Add this below AppModule struct {

// highlight-start
type Validator struct {
	Address string `json:"address"`
	Moniker string `json:"moniker"`
}

type Validators []Validator

func (vs Validators) Formatted() string {
	output := ""
	for _, val := range vs {
		output += fmt.Sprintf(`{"address":"%s","moniker":"%s"},`, val.Address, val.Moniker)
	}

	return output[:len(output)-1]
}
// highlight-end
```

## Implement the endblock

```go title="x/reviews/module.go"
var _ appmodule.HasEndBlocker = AppModule{} // optional

func (am AppModule) EndBlock(ctx context.Context) error {
	stakingVals, err := am.keeper.StakingKeeper.GetAllValidators(ctx)
	if err != nil {
		return err
	}

	validators := Validators{}
	for _, val := range stakingVals {
		if !val.IsBonded() {
			continue
		}

		validators = append(validators, Validator{
			Address: val.OperatorAddress,
			Moniker: val.Description.Moniker,
		})
	}

	endBlockSudoMsg := fmt.Sprintf(`{"set_validators":{"validators":[%s]}}`, validators.Formatted())

	addr := sdk.MustAccAddressFromBech32("roll14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9sjczpjh")

	res, err := am.keeper.ContractKeeper.Sudo(sdk.UnwrapSDKContext(ctx), addr, []byte(endBlockSudoMsg))
	fmt.Println("EndBlockSudoMessage", endBlockSudoMsg)
	fmt.Println("EndBlockSudoResult", res, "error", err)

	return nil
}
```

## Start Testnet
```bash
make sh-testnet

```



----

## Test Deployment

<!-- TODO: make reference that they should have already built the contract by here -->

```bash
rolld tx wasm store ./validator-reviews-contract/artifacts/validator_reviews_contract.wasm --from=acc0 --gas=auto --gas-adjustment=2.0 --yes
# rolld q wasm list-code

rolld tx wasm instantiate 1 '{}' --no-admin --from=acc0 --label="validator_reviews" --gas=auto --gas-adjustment=2.0 --yes

rolld q wasm list-contracts-by-creator roll1hj5fveer5cjtn4wd6wstzugjfdxzl0xpg2te87
REVIEWS_CONTRACT=roll14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9sjczpjh

rolld q wasm state smart $REVIEWS_CONTRACT '{"validators":{}}'


MESSAGE='{"write_review":{"val_addr":"rollvaloper1hj5fveer5cjtn4wd6wstzugjfdxzl0xpmhf3p6","review":"hi this is a review"}}'
rolld tx wasm execute $REVIEWS_CONTRACT "$MESSAGE" --from=acc0 --gas=auto --gas-adjustment=2.0 --yes

rolld q wasm state smart $REVIEWS_CONTRACT '{"reviews":{"address":"rollvaloper1hj5fveer5cjtn4wd6wstzugjfdxzl0xpmhf3p6"}}'

MESSAGE='{"write_review":{"val_addr":"notavaldiator","review":"hi this is a review"}}'
rolld tx wasm execute $REVIEWS_CONTRACT "$MESSAGE" --from=acc0 --gas=auto --gas-adjustment=2.0 --yes
```
