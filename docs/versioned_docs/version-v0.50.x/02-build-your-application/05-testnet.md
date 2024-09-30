---
title: "Name Service"
sidebar_label: "Testnet"
sidebar_position: 5
slug: /build/name-service-testnet
---

# Running your Application

:::note Synopsis
Congrats!! You built your first network already. You are ready to run a local testnet environment to verify it works.

* Building your application executable
* Running a local testnet
* Interacting with the network
:::

### Launch The Network

Use the `sh-testnet` command *(short for shell testnet)* to quickly build your application, generate example wallet accounts, and start the local network on your machine.

```bash
# Run a quick shell testnet
make sh-testnet
```

The chain will begin to create (mint) new blocks. You can see the logs of the network running in the terminal.

### Interact Set Name

Using the newly built binary *(rolld from the --bin flag when the chain was created)*, you are going to execute the `set` transaction to your name. In this example, use "alice". This links account `acc1` address to the desired name in the keeper.

Then, resolve this name with the nameservice lookup. `$(rolld keys show acc1 -a)` is a substitute for the acc1's address. You can also use just `roll1hj5fveer5cjtn4wd6wstzugjfdxzl0xpg2te87` here.

```bash
rolld tx nameservice set alice --from=acc1 --yes

# You can verify this transaction was successful
# By querying it's unique ID.
rolld q tx 565CE77057ACBF6FB5D174231455E61E65009CD628971937C19201328E0A1FFD
```

## Interaction Get Name

Now you are going to get the name of a wallet. A nested command `$(rolld keys show acc1 -a)` gets the unique address of the acc1 account added when you started the testnet.

```bash
rolld q nameservice resolve roll1efd63aw40lxf3n4mhf7dzhjkr453axur57cawh --output=json

rolld q nameservice resolve $(rolld keys show acc1 -a) --output=json
```

The expected result should be:

```json
{
  "name": "alice"
}
```

:::note
When you are ready to stop the testnet, you can use `ctrl + c` or `killall -9 rolld`.
:::


Your network is now running and you have successfully set and resolved a name! 🎉
