---
title: "Token Factory"
sidebar_label: "Token Factory"
slug: /demo/tokenfactory
# sidebar_position: 1
---

# Tokenfactory

You will build a new chain with [TokenFactory](https://github.com/strangelove-ventures/tokenfactory), enabling any account to create, transfer, and interact with fractionalized native tokens.

:::note Warning
Some parts of this tutorial will not have the added context about spawn's inner workings or how commands work. Run through [Build Your Application](../02-build-your-application/01-nameservice.md) for this context.
:::

## Prerequisites
- [System Setup](../01-setup/01-system-setup.md)
- [Install Spawn](../01-setup/02-install-spawn.md)


## Create your chain

Build a new chain that has TokenFactory configured. By default, it is enabled.

```bash
GITHUB_USERNAME=rollchains

spawn new rollchain \
--consensus=proof-of-stake \
--bech32=roll \
--denom=uroll \
--bin=rolld \
--disabled=cosmwasm,block-explorer \
--org=${GITHUB_USERNAME}
```


## Start the testnet

:::note Note
If `make sh-testnet` does not start due to a port bind error, you can kill your previously running testnet with `killall -9 rolld`, then try again.
:::


```bash
# move into the chain directory
cd rollchain

# - Installs the binary
# - Setups the default keys with funds
# - Starts the chain in your shell
make sh-testnet
```

## Confirm tokenfactory is enabled

```bash
rolld q tokenfactory params
```

<details>
<summary>params output</summary>

The `denom_creation_fee` is a cost the application can set for creating new tokens by default, there is no cost.

The `denom_creation_gas_consume` is the amount of indirect resource cost to consume for creating a new token.
It is a more indirect approach to charging and is a better experience overall for developers on a network.

```bash
params:
  denom_creation_fee: []
  denom_creation_gas_consume: "100000"
```
</details>

## Create a token


```bash
# Create a denom (native token)
# - gas is is amount of compute resources to allocate.
rolld tx tokenfactory create-denom mytoken --from=acc0 --chain-id=localchain-1 --yes
```

## Verify the token was created

```bash
# Get our account address for the acc0 wallet / key.
# acc0 is roll1hj5fveer5cjtn4wd6wstzugjfdxzl0xpg2te87
rolld q tokenfactory denoms-from-creator $(rolld keys show acc0 -a)
```

<details>
<summary>denoms-from-creator output</summary>
```bash
denoms:
- factory/roll1hj5fveer5cjtn4wd6wstzugjfdxzl0xpg2te87/mytoken
```
</details>

The output shows a denom with the named `factory/roll1hj5fveer5cjtn4wd6wstzugjfdxzl0xpg2te87/mytoken`

:::note Note
Why did it add extra data to the token?

Imagine there are 2 people, both named John. If only the name John is used, which John is it talking about? More information must be added to the name to make it unique. This is the same concept, but with tokens. The extra data is added to ensure the token is unique while it can still contain the same base name. With tokenfactory, the creators name is placed in the token. Read more about [naming collisions](https://en.wikipedia.org/wiki/Naming_collision).
:::


## Modify token metadata

Clients (websites, frontends, users) may wish to see more information about the token. This is where metadata comes in. You can add a ticker symbol, description, and decimal places to the token.

The Interchain uses 6 decimal places as the default standard. This process of expressing fractions of a value in whole numbers is called [fixed-point arithmetic](https://en.wikipedia.org/wiki/Fixed-point_arithmetic) and is used for financial precision. This means that 1 token is really 1,000,000 (10^6) of these micro base tokens. If I want to send you 0.5 of a token, I really send you 500,000 of these micro base tokens on the backend.

```bash
# 'Denom' is short for denomination.
DENOM=factory/roll1hj5fveer5cjtn4wd6wstzugjfdxzl0xpg2te87/mytoken
DESCRIPTION="My token description"

rolld tx tokenfactory modify-metadata $DENOM MYTOKEN "$DESCRIPTION" 6 --from acc0 --yes
```

## Verify the token metadata

```bash
rolld q bank denom-metadata $DENOM
```

<details>
<summary>bank denom-metadata output</summary>
```bash
metadata:
  base: factory/roll1hj5fveer5cjtn4wd6wstzugjfdxzl0xpg2te87/mytoken
  denom_units:
  - aliases:
    - MYTOKEN
    denom: factory/roll1hj5fveer5cjtn4wd6wstzugjfdxzl0xpg2te87/mytoken
  - aliases:
    - factory/roll1hj5fveer5cjtn4wd6wstzugjfdxzl0xpg2te87/mytoken
    denom: MYTOKEN
    exponent: 6
  description: My token description
  display: MYTOKEN
  name: factory/roll1hj5fveer5cjtn4wd6wstzugjfdxzl0xpg2te87/mytoken
  symbol: MYTOKEN
```
</details>

## Create new tokens to transfer

The base token structure is created, but no tokens actually exists yet. Mint new tokens to then be able to transfer them between accounts.

```bash
# Mint 5,000,000 micro mytoken. By default this goes to the token creator.
rolld tx tokenfactory mint 5000000$DENOM --from acc0 --yes

# Verify token creator balance: roll1hj5fveer5cjtn4wd6wstzugjfdxzl0xpg2te87
rolld q bank balances $(rolld keys show acc0 -a)
```


<details>
<summary>bank balances output</summary>
```bash
balances:
- amount: "5000000"
  denom: factory/roll1hj5fveer5cjtn4wd6wstzugjfdxzl0xpg2te87/mytoken
- amount: "900"
  denom: test
- amount: "9000000"
  denom: uroll
pagination:
  total: "3"
```
</details>


## Create new tokens for another account

While you could mint tokens followed by a manual `tx bank send` transfer, you can also mint-to another account directly.

```bash
# Mint 1,000,000 to another account
rolld tx tokenfactory mint-to $(rolld keys show acc1 -a) 1000000$DENOM --from acc0 --yes

rolld q bank balances $(rolld keys show acc1 -a)
```

<details>
<summary>mint-to output</summary>
```bash
balances:
- amount: "1000000"
  denom: factory/roll1hj5fveer5cjtn4wd6wstzugjfdxzl0xpg2te87/mytoken
- amount: "800"
  denom: test
- amount: "10000000"
  denom: uroll
pagination:
  total: "3"
```

note, you can check for just a specific token balance with
```bash
rolld q bank balance $(rolld keys show acc0 -a) $DENOM
```

```bash
balance:
  amount: "5000000"
  denom: factory/roll1hj5fveer5cjtn4wd6wstzugjfdxzl0xpg2te87/mytoken
```
</details>

## Transfer tokens

Now with tokens minted, either user can transfer them as they please between any accounts. Even ones that do not yet have any tokens.

```bash
# Send 7 base micro tokens from acc0 to acc1
rolld tx bank send acc0 $(rolld keys show acc1 -a) 7$DENOM --from acc0 --yes

# Verify the 7 base tokens sent and has increased to 1000007, or 1.000007
rolld q bank balances $(rolld keys show acc1 -a)
```

## Burn tokens

If you wish to remove tokens from the system, you can burn them from the admin account.

```bash
# Burn micro tokens from account
rolld tx tokenfactory burn 123$DENOM --from acc0 --yes

# Verify the tokens have been reduced
rolld q bank balances $(rolld keys show acc0 -a)
```

## Conclusion

In this tutorial, you built a new chain with the TokenFactory feature, launched a testnet for it, and created a new native token. You minted tokens, transferred them between accounts, and burned them. These tokens could be kept internally for some personal or application based accounting, or transferred to other chains via IBC. This is showcased in the [IBC Transfer Demo](../03-demos/01-ibc-transfer-demo.md).
