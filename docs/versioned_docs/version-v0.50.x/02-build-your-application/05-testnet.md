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

Using the newly built binary executable *(rolld from the --bin flag when the chain was created)*, you are going to execute the `set` action to your name. In this example, use "alice". This links account `acc1` address to the desired name in the keeper.

You can either query or set data in the network using the command executable. If you wish to perform an action you submit a transaction (tx). If you wish to read data  you are querying (q). The next sub command specifies which module will receive the action on. In this case, the `nameservice` module since our module is named nameservice. Then the `set` command is called, which was defined in the autocli.go.

```bash
rolld tx nameservice set alice --from=acc1 --yes

# You can verify this transaction was successful
# By querying it's unique ID.
rolld q tx EC3FBF3248E24B5FEB6A5F7F35BBB4634E9C75587119E3FBCF5C1FED05E5A399
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
