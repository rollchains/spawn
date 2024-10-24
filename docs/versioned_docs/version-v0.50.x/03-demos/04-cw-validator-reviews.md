---
title: "CW Validator Reviews"
sidebar_label: "CosmWasm Validator Reviews"
slug: /demo/cw-validator-reviews
---

# CosmWasm Validator Reviews

You will build a new chain with [CosmWasm](https://cosmwasm.com), enabling a proof-of-stake validator review system. You will write a contract to collect and manage validator reviews, integrate it with the chain, and update validator data automatically through a Cosmos-SDK endblocker module.

There are easy ways to get validators in a cosmwasm smart contract. The goal of this tutorial is to teach how to pass data from the SDK to a contract.

## Prerequisites
- [System Setup](../01-setup/01-system-setup.md)
- [Install Spawn](../01-setup/02-install-spawn.md)
- [Rust + CosmWasm](../01-setup/01-system-setup.md)

## Setup the Chain

Build a new blockchain with CosmWasm enabled. THen generate a new module from the template for the reviews.

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

Once the chain is started, continue on to the contract building steps

## CosmWasm Build Contract

CosmWasm has a template repository that is used to generate new contracts. A minimal contract will be built with the `validator-reviews-contract` name provided on a specific commit.

```bash
cargo generate --git https://github.com/CosmWasm/cw-template.git \
    --name validator-reviews-contract \
    -d minimal=true --tag a2a169164324aa1b48ab76dd630f75f504e41d99
```

A new folder will be created with the contract template.

```bash
code validator-reviews-contract/
```


### Set State

The contract state and base structure is set in the state.rs file. There are 2 groups of data that must be managed, validators and the reviews for validators.

- `Validators` have unique addresses and name stored on the chain. This data will be passed from the Cosmos-SDK to the contract.
- `Reviews` will save a user wallets address and their text reviews for a validator.

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

### Set Inputs

By default contracts get 3 messages, `InstantiateMsg`, `ExecuteMsg`, and `QueryMsg`.

- **Instantiate** allows initial contracts setup with a configuration desired. This is not used for us. Keep it empty.
- **Execute** is where the main logic of the contract is. Add a `WriteReview` message to allow users to write reviews. The user must know who they want to write a review for and what it says.
- **Query** is for reading data from the contract. Add 2 queries, one to get all validators available and one to get reviews for a specific validator.

The `SudoMsg` is a default type not typically used. `Sudo` stands for `Super User DO` where the super user is the chain. **Only** the chain can submit data to this message type. A `SetValidators` message is added to allow the chain to update the validators list within the contract. This is the pass through from the Cosmos-SDK to the contract.

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

### Set new error

For a better experience, a new error is added to the contract. This will be used when a validator is not found in the list of validators. Users should not be allowed to post reviews for non existent validators.

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
    #[error("The validator is not found")]
    NoValidatorFound {},
    // highlight-end
}
```

### Imports

The imports required for this next section are provided here. Paste these at the top of the file to get syntax highlighting.

```rust title="src/contract.rs"
// highlight-start
use crate::state::{Reviews, REVIEWS, VALIDATORS};
use cosmwasm_std::to_json_binary;
use std::collections::BTreeMap;
// highlight-end
```

### Modify Instantiate Message

Even though no extra data is passed through to the setup method, an empty list of validators is saved to storage. This way if we try to get validators from the contract **before** any have been set by the chain, it returns nothing instead of an error.

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

### Add Execute Logic

When the user sends a message, the contract needs to first check if the validator exist. It does this by loading the validators state and looping through all the validators to see if the one the user requested if in the list. If it is not, it returns to the user that the validator is not found. If it is found then the contract loads all reviews a validator has. If there are none, it creates an empty list of reviews since this will be the first one. The user's review is then added to the list and saved back to storage.

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

### Add Queries

It is only useful to set reviews if you can also get them back. The first query for `Validators` is just a helper method so users can see who they are allowed to review. The second query is for `Reviews` and takes a validator address as a parameter. This will return all reviews for that validator.

To get reviews for all validators, a user would need to query `Validators`, then iterate through all the addresses and query `Reviews` for each one.

```rust title="src/contract.rs"
#[cfg_attr(not(feature = "library"), entry_point)]
// highlight-start
pub fn query(deps: Deps, _env: Env, msg: QueryMsg) -> StdResult<Binary> {
    // note: ^^ deps & msg are not underscored ^^
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

### Add New Sudo Message Type

The chain extended portion of this contract is now added. It is where the validator logic is actually set and saved to storage. As the validator set changes (nodes stop, new nodes come online), the chain will update the contract right away.

```rust title="src/contract.rs"
// highlight-start
// Insert at the bottom of the file
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



### Build the Contract

Build the contract to get the cosmwasm wasm binary. This converts it from english programming rust text to 0s and 1s that the chain can understand. The `optimize` script requires docker to be installed and running. Make sure you followed the setup prerequisites at the top of the page and have the docker service or docker desktop installed.

```bash
# run the build optimizer (from source -> contract wasm binary)
cargo run-script optimize
```

The .wasm file is then saved to `./artifacts/validator_reviews_contract.wasm`.

## Modify the Module

The contract is complete but we need to pass the data into the contract from the chain. This is done through the cosmos-sdk reviews module generated earlier. The module will be updated to include the wasm contract and the endblocker will be updated to pass the validator data to the contract.

### Setup the Keeper

We must give our code access to other modules on the chain. The wasm module is required to interact with the contract. The staking module is required to get the list of validators. This keeper is the access point for all the specific logic.

Add the imports, keeper setup, and new keeper output.

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

### 'Fix' keeper_test


:::note Warning
Testing wasm requires significantly more setup for the test environment. For now, just add a blank reference here.
:::

```go title="x/reviews/keeper/keeper_test.go"
func SetupTest(t *testing.T) *testFixture
    ...

	// Setup Keeper.
    // highlight-next-line
	f.k = keeper.NewKeeper(encCfg.Codec, runtime.NewKVStoreService(keys[types.StoreKey]), logger, &wasmkeeper.Keeper{}, f.stakingKeeper, f.govModAddr)
}
```

### Dependency Inject (v2)

Similar to the keeper_test issue, CosmWasm does not have support for Cosmos-SDK v2 depinject. This will be updated in the future. For now, set the keeper to nil and provide Staking reference. You do not need to know what this does. Just resolve the error on the line with a copy paste.

```go title="x/reviews/depinject.go"
func ProvideModule(in ModuleInputs) ModuleOutputs {
    ...

    // highlight-next-line
    k := keeper.NewKeeper(in.Cdc, in.StoreService, log.NewLogger(os.Stderr), nil, &in.StakingKeeper, govAddr)
}
```

### Fix app.go references

The main application needs the access to the wasm and staking logic as well. Fix the missing imports and add them in like so.

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

it is now time to use these modules and write the logic to pass data to the contract from the chain.

### Module Core Logic

The CosmWasm contract requires data in a specific format. You can review this back in the `src/state.rs` file. Since the chain only is passing validator data over, we need to convert this into Go code manually. The `Validator` struct is created to match the contract. The CosmWasm contract expects a JSON formatted input. This input is put together with the `Formatted` method on a list of validators. The chain could have just 1 validator, or several hundred. This method will convert them all into the correct format for the list we are to pass.

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

    // remove the trailing comma from the last output append.
	return output[:len(output)-1]
}
// highlight-end
```

### Implement the EndBlocker

To pass data we must first get the data. Using the `GetAllValidators` method from the staking module, all validators are now accessible for the logic to use. Loop through these validators and only add the ones that are bonded (active) to the list of validators. If they are bonded, they are added to the list.

Once all validators have been processed the `endBlockSudoMsg` gets them into the JSON format required. The format is out of scope but a high level overview
- `SetValidators` in the code becomes `set_validators`, called [snake case](https://simple.wikipedia.org/wiki/Snake_case).
- The `SetValidators` type in rust has the element called `validators` which is an array of `Validator` objects. This is the `validators` array in the JSON.
- Each `Validator` object has an `address` and `moniker` field. These are the `address` and `moniker` fields in the JSON, called from the Formatted() method.

```go title="x/reviews/module.go"
// Paste this anywhere within the file
// highlight-start
func (am AppModule) EndBlock(ctx context.Context) error {
	stakingVals, err := am.keeper.StakingKeeper.GetAllValidators(ctx)
	if err != nil {
		return err
	}

	validators := Validators{}
	for _, val := range stakingVals {
        // if it is not active, skip it
		if !val.IsBonded() {
			continue
		}

		validators = append(validators, Validator{
			Address: val.OperatorAddress,
			Moniker: val.Description.Moniker,
		})
	}

    // The first contract created from acc0
    addr := "roll14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9sjczpjh"
	contract := sdk.MustAccAddressFromBech32(addr)

    // SudoMsg format for the contract input.
    // example: {"set_validators":{"validators":[{"address":"ADDRESS","moniker": "NAME"}]}}
    endBlockSudoMsg := fmt.Sprintf(`{"set_validators":{"validators":[%s]}}`, validators.Formatted())
	fmt.Println("EndBlockSudoMessage Format:", endBlockSudoMsg)

    // When the network first starts up there is no contract to execute against (until uploaded)
    // This returns an error but is expected behavior initially.
    // You can not return errors in the EndBlocker as it is not a transaction. It will halt the network.
    //
    // This is why the error is only printed to the logs and not returned.
    //
    // A more proper solution would set the contract via a SDK message after it is uploaded.
    // This is out of scope for this tutorial, but a future challenge for you.
	res, err := am.keeper.ContractKeeper.Sudo(sdk.UnwrapSDKContext(ctx), contract, []byte(endBlockSudoMsg))
    if err != nil {
        fmt.Println("EndBlockSudoResult", res)
        fmt.Println("EndBlockSudoError", err)
    }

	return nil
}
// highlight-end
```

## Test Deployment

### Start Testnet

Begin the CosmWasm testnet with the custom EndBlocker logic. You will see errors every block. This is expected and is explained in the EndBlock code why this is the case.

```bash
make sh-testnet
```

### Upload Contract

Make sure you are in the `rollchain` directory to begin interacting and uploading the contract to the chain.

```bash
CONTRACT_SOURCE=./validator-reviews-contract/artifacts/validator_reviews_contract.wasm
rolld tx wasm store $CONTRACT_SOURCE --from=acc0 --gas=auto --gas-adjustment=2.0 --yes
# rolld q wasm list-code
```

### Instantiate our Contract

```bash
rolld tx wasm instantiate 1 '{}' --no-admin --from=acc0 --label="reviews" \
    --gas=auto --gas-adjustment=2.0 --yes

rolld q wasm list-contracts-by-creator roll1hj5fveer5cjtn4wd6wstzugjfdxzl0xpg2te87

REVIEWS_CONTRACT=roll14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9sjczpjh
```

### Verify data

```bash
rolld q wasm state smart $REVIEWS_CONTRACT '{"validators":{}}'
```

## Write a review

```bash
MESSAGE='{"write_review":{"val_addr":"rollvaloper1hj5fveer5cjtn4wd6wstzugjfdxzl0xpmhf3p6","review":"hi reviewing."}}'
rolld tx wasm execute $REVIEWS_CONTRACT "$MESSAGE" --from=acc0 \
    --gas=auto --gas-adjustment=2.0 --yes
```

### Verify the review

```bash
rolld q wasm state smart $REVIEWS_CONTRACT '{"reviews":{"address":"rollvaloper1hj5fveer5cjtn4wd6wstzugjfdxzl0xpmhf3p6"}}'
```

### Write a review for a non-validator

```bash
MESSAGE='{"write_review":{"val_addr":"NotAValidator","review":"hi this is a review"}}'

rolld tx wasm execute $REVIEWS_CONTRACT "$MESSAGE" --from=acc0 \
    --gas=auto --gas-adjustment=2.0 --yes
```
