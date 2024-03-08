# Spawn Demo

[Spawn Introduction Article](http://TODO)

Introducing spawn, the developer tool that gets your cosmos chain up and running quickly. Allowing you to focus on what matters: your product.

Let's create a new chain called 'rollchain'. We are going to set some of the defining characteristics such as
- which modules to disable from the template *if any*
- Wallet prefix (bech32)
- Token name
- and the binary

```bash
spawn new rollchain --disable=globalfee --bech32=roll --denom=uroll --bin=rolld --org=rollchains
```

The chain is now created and we can start to write our application logic on top.

## Generate New Module

Let's build a nameservice module for this example.

Move into the 'rollchain' directory, then generate the new module with the following command:

```bash
cd rollchain

spawn module new nameservice
```

This creates a new template module with the name `nameservice` in the `x` and `proto` directories. This is also automatically connected to your app.go and ready for application use.

## Setup Messages

The protobuf files have been automatically generated for you with a default `Params` section. Building off this, we are adding our new messages to query and set a wallets name.

Open the proto/nameservice directory. Edit `tx.proto` *(proto/nameservice/v1/tx.proto)* to add the transaction setter message.

```proto

  rpc SetServiceName(MsgSetServiceName) returns (MsgSetServiceNameResponse);
}

message MsgSetServiceName {
  option (cosmos.msg.v1.signer) = "sender";

  string sender = 1 [(cosmos_proto.scalar) = "cosmos.AddressString"];

  string name = 2;
}

message MsgSetServiceNameResponse {}
```

![proto/nameservice/v1/tx.proto file](https://github.com/rollchains/spawn/assets/31943163/73a583e2-9edd-471f-ada6-1010d0dbf072)

Find `query.proto` *(proto/nameservice/v1/query.proto)* and add the following

```proto

  rpc ResolveName(QueryResolveNameRequest) returns (QueryResolveNameResponse) {
    option (google.api.http).get = "/nameservice/v1/names/{wallet}";
  }
}

message QueryResolveNameRequest {
  string wallet = 1;
}

message QueryResolveNameResponse {
  string name = 1;
}
```

![proto/nameservice/v1/query.proto file](https://github.com/rollchains/spawn/assets/31943163/234a13d7-be62-492d-961c-63e92d7543d9)

Now build the proto files into their .go counterparts. Once generated, the new interface requirements are automatically satisfied for us for both MsgServer & Querier.

```bash
make proto-gen
```

![make proto-gen](https://github.com/rollchains/spawn/assets/31943163/c51bf57c-e83a-4004-8041-9b1f3d3a24f4)


# Keeper Storage Structure

Now set the data structure map in the keeper *x/nameservice/keeper/keeper.go* to store the wallet to name pair.

```go

type Keeper struct {
	...

	NameMapping collections.Map[string, string]

  ...
}

...

func NewKeeper() Keeper {
  ...

  k := Keeper{
    ...

    NameMapping: collections.NewMap(sb, collections.NewPrefix(1), "name_mapping", collections.StringKey, collections.StringValue),
  }

}
```

![keeper NewKeeper NameMapping](https://github.com/rollchains/spawn/assets/31943163/47ed4a41-4df2-4a5d-9ac5-bfb23aeefd94)

---

# Application Logic

Update the `x/nameservice/keeper/msg_server.go` now to set on the newly created map.

```go
func (ms msgServer) SetServiceName(ctx context.Context, msg *types.MsgSetServiceName) (*types.MsgSetServiceNameResponse, error) {
	if err := ms.k.NameMapping.Set(ctx, msg.Sender, msg.Name); err != nil {
		return nil, err
	}

	return &types.MsgSetServiceNameResponse{}, nil
}
```

and also for the query in `x/nameservice/keeper/query_server.go`

```go
func (k Querier) ResolveName(goCtx context.Context, req *types.QueryResolveNameRequest) (*types.QueryResolveNameResponse, error) {
	v, err := k.Keeper.NameMapping.Get(goCtx, req.Wallet)
	if err != nil {
		return nil, err
	}

	return &types.QueryResolveNameResponse{
		Name: v,
	}, nil
}
```

---

# Setting the CLI Client

The Cosmos-SDK recently introduced the AutoCLI. This method simplifies the setup of the CLI client for transactions and queries.

### Query

Modify the Query RPC options to support the recently introduced ResolveName message.

```go
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{

                 ...

			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "ResolveName",
					Use:       "resolve [wallet]",
					Short:     "Resolve the name of a wallet address",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "wallet"},
					},
				},
                ...
      },
    }
  }
```

![AutoCLI Query](https://github.com/rollchains/spawn/assets/31943163/fefe8c7d-88b5-42d5-afd9-cb33cd22df16)


### Transaction

And also for setting the Transaction

```go
		Tx: &autocliv1.ServiceCommandDescriptor{

      ...

			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "SetServiceName",
					Use:       "set [name]",
					Short:     "Set the mapping to your wallet address",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "name"},
					},
				},
        ...
      }
    }
```

![AutoCLI Tx](https://github.com/rollchains/spawn/assets/31943163/e945c898-415c-4d22-8bb3-b8af34a44cee)


---

# Testnet

With the module now completed, it is time to run the local testnet to validate our additions.

The following will build the binary, add keys to the test keyring, and start a single validator instance with the latest chain binary.

```bash
make sh-testnet
```

The chain will begin to mint new blocks, which you can interact with.

### Interaction

Using the newly built binary (rolld from the --bin flag when we created the chain), we are going to execute the `set` transaction to "myname". This links user1's address (in the keyring) to the desired name in the keeper.

Then, we resolve this name with the nameservice lookup. `$(rolld keys show user1 -a)` is a substitute for the user1's address. You can also use just `roll1hj5fveer5cjtn4wd6wstzugjfdxzl0xpg2te87` here.

```bash
rolld tx nameservice set myname --from=user1 --yes

# rolld q tx 088382C43C35440676438359B88899D97A8092F34BBDADD32345498297D332BA

sleep 2

rolld q nameservice resolve $(rolld keys show user1 -a) --output=json
```

The expected result should be:

```json
{
  "name": "myname"
}
```
