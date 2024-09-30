---
title: "Network Security"
sidebar_label: "Network Security Types"
slug: /learn/consensus-security
# sidebar_position: 1
---

# Network Security Types

## Prerequisites
- [System Setup](../01-setup/01-system-setup.md)
- [Install Spawn](../01-setup/02-install-spawn.md)

Learn about the different network security methods you can create with spawn and the different pros and cons of each. These are called [consensus algorithms](https://en.wikipedia.org/wiki/Consensus_(computer_science)) and are how the network agrees on what actions get put through, and which are rejected.

:::note Danger
This topic is a little more advanced due to technical speak. It is condensed to as a high level of an overview as possible. Please contribute if you can make it more accessible to average readers.
:::

## Choose for me (TLDR)

If you just want to build an application and don't want to focus on tokenomics or game theory, use [proof of authority](#proof-of-authority-poa). If a token is part of your product or required to financially reward users with some lock mechanism and long term incentives, use [proof of stake](#proof-of-stake-pos). If you have plans for a large amount of value to be secured or don't want to run your own network nodes, use [interchain security](#interchain-security-ics).

## Proof of Authority (PoA)

### Default
If you do not know which security module best fits for you, use this one. The source code for this feature can be found [here](https://github.com/strangelove-ventures/poa). The most popular example of this security model is [Circle's](https://www.circle.com/en/) USDC issuance network, [Noble](https://www.noble.xyz/) ([twitter](https://twitter.com/noble_xyz)).

### What it does
If you have an application and you want the network to run as efficiently as possible with a trusted set of internal nodes or external trusted 3 parties (validators), use [Proof of Authority](https://en.wikipedia.org/wiki/Proof_of_authority).

### Create a PoA network
To create a new network with proof-of-authority, use the `--consensus=proof-of-authority` flag. If `--consensus` is not present, a selector UI will appear in your terminal to see all options.

```bash
-> $ spawn new mychain
Consensus Selector (( enter to toggle ))
  Done
  ✔ proof-of-authority
  proof-of-stake
  interchain-security
```

```bash
spawn new rollchain \
  --consensus=proof-of-authority \
  --disabled=cosmwasm,block-explorer
```

### Considerations

#### Benefits
- Fast Transactions: Fewer people need to agree, so it can process transactions very quickly.
- Less Power Usage: Doesn’t require much electricity, making it more eco-friendly.
- Easy to Maintain: Only a few trusted people are in charge, making it simpler to run.
- Stable Performance: Because only a few people make decisions, things tend to run smoothly and predictably.

#### Downsides
- Centralized Control: A small group of people are in charge, which can lead to concerns about too much power in one place.
- Requires Trust: You have to trust the people in charge to make fair decisions, which can be risky.
- Less Diversity: Having fewer people in control means less variety in opinions and locations, which might be a problem if those people get compromised.
- Less Community Involvement: Regular users don’t have much of a role in helping the system, so it feels less like a community effort.

## Proof of Stake (PoS)

### What it does
You can have the value of a network back itself by users risking their own tokens to prove they are trustworthy. This is called [Proof of Stake](https://en.wikipedia.org/wiki/Proof_of_stake). Believers in an application lock their tokens to earn a small portion of rewards, similar to a bank account. However, if they misbehave by trying to cheat the system or submit bad actions, the network will take a portion of their value (usually 5 - 10%) as a penalty.

This security type is useful when you want a more distributed network that can be run by anyone with some holdings in the network. It is a trustless way to secure a network and the most popular security model in the ecosystem currently (Sept 2024).

### Create a PoS network
To create a new network with proof-of-stake, use the `--consensus=proof-of-stake` flag. If `--consensus` is not present, a selector UI will appear in your terminal to see all options.

```bash
-> $ spawn new mychain
Consensus Selector (( enter to toggle ))
  Done
  proof-of-authority
  ✔ proof-of-stake
  interchain-security
```

```bash
spawn new rollchain \
  --consensus=proof-of-stake \
  --disabled=cosmwasm,block-explorer
```

### Considerations

#### Benefits
- More People Involved: Anyone can participate if they’re willing to invest, which makes the system feel more balanced and community-driven.
- Eco-Friendly: Like PoA, PoS doesn’t use much electricity, making it good for the environment.
- More Fairness: The system allows many people to help make decisions, reducing the chance of one group having too much control.
- Grows with the Community: More people can get involved as the system grows, making it scalable and inclusive.

#### Downsides
- Slower Decision Making: Because more people are involved, it can take longer to reach a decision due to governance and politics
- Wealthy Have More Power: The more you invest, the more influence you have, which can lead to rich people having more control.
- Complex to Get Started: It can be harder for someone new to understand how to participate compared to systems with fewer decision-makers.
- Risk of Losing Investment: If you make a mistake or act dishonestly, you could lose your money, which adds some financial risk.

## Interchain Security (ICS)

### What it does

Interchain security shares the economic proof of stake security of a larger parent provider with a sub network, called a consumer *(since they consume security)*. This is useful when you want to create a new network that is secure from day one, without having to bootstrap a new set of network operators and validators, and have alignment with the parent. Current networks utilizing this are [CosmosHub](https://cosmos.network/interchain-security/), [Stride](https://www.stride.zone/), and [Lido's Neutron](https://www.neutron.org/). If you are from ethereum, this is similar to an [Actively Validated Services (AVS) on Eigenlayer](https://app.eigenlayer.xyz/avs).

The cost of running these networks is relatively low as you just pay a portion of your networks fees. This is a great way to take an application from a testnet to a mainnet with a trusted security model, especially if your application deals with a lot of possible monetary value. To compromise the network, an attacker would have to compromise the more secure parent network, which is a very high bar.

### Create an ICS Consumer network
To create a new network with interchain-security, use the `--consensus=interchain-security` flag. If `--consensus` is not present, a selector UI will appear in your terminal to see all options.

::note Note
Spawn does not support creating provider networks.
:::

```bash
-> $ spawn new mychain
Consensus Selector (( enter to toggle ))
  Done
  proof-of-authority
  proof-of-stake
  ✔ interchain-security
```

```bash
spawn new rollchain \
  --consensus=interchain-security \
  --disabled=cosmwasm,block-explorer
```

### Considerations

#### Benefits
- No Need for Consumer Chain to Build Its Own Security: The smaller or newer chain doesn't have to recruit its own set of people to protect it. It automatically benefits from the protection of the parent network, saving time and effort.
- Easier Node Bootstrapping: Since the parent network provides the security, the consumer chain doesn’t have to build a large number of participants to protect the system from scratch. This makes launching a new chain much faster and easier.
- Shared Trust: By using the same security as a well-known, established chain, the consumer chain inherits the trust and credibility of the parent network, making it more appealing to users and developers.
- Aligned Interests: Since the parent network's security also protects the consumer chain, both chains have an interest in maintaining a secure, well-functioning system. This alignment reduces the risk of conflicts between the two.



#### Downsides
- Dependence on Parent Network: The consumer chain becomes dependent on the parent network. If something goes wrong with the parent network’s security or operations, the consumer chain is also affected, even if it is unrelated to the issue.
- Limited Autonomy: The consumer chain may have less control over its own security decisions, since it’s tied to the security model of the parent network. This could limit flexibility in responding to specific needs or changes.
- Potential Congestion: If many consumer chains share the same parent network, the shared security system might become congested or stretched thin, leading to slower response times or performance issues.
- Complexity in Governance: Any changes to the shared security might require coordination between both the parent network and the consumer chain, adding complexity to decision-making and governance.
- Risk of Centralization: The reliance on a single parent network’s security model might lead to centralization, where a handful of large chains dominate the ecosystem. This reduces the diversity of security models and could concentrate power.

## Conclusion

You have now learned about different network security types, how to select different ones, and the pros and cons of each. You can now create a new chain with the security model that best fits your application's needs.
