# Spawn Demo

[Spawn Introduction Article](http://TODO)

Introducing spawn, the developer tool that gets your cosmos chain up and running quickly. Allowing you to focus on what matters: your product.

Let's create a new chain called 'rollchain'. We are going to set some of the defining characteristics such as
- which modules to disable from the template *if any*
- Wallet prefix (bech32)
- Token name
- and the binary

```bash
spawn new rollchain --disable=cosmwasm,globalfee --bech32=roll --denom=uroll --bin=rolld --org=rollchains
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

Head into the proto/nameservice directory and find `query.proto` *(proto/nameservice/v1/query.proto)* and add the following lines in

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

## TODO: image of the file proto/nameservice/v1/query.proto

Then edit `tx.proto` *(proto/nameservice/v1/tx.proto)* to add the transaction setter message.

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

## TODO: image of the file proto/nameservice/v1/tx.proto

Now we build the proto files into their .go counterparts. Once generated, the new interface requirements are automatically satisfied for us for both MsgServer & Querier.

```bash
make proto-gen
```

## TODO: image of the terminal here after a make proto-gen

```
$ make proto-gen

Generating Protobuf files
Generating gogo proto code
Generating pulsar proto code
 [+] Moving: ./nameservice to ./api/nameservice
3INF Applied RPC Stub module=nameservice type=query name=ResolveName req=QueryResolveNameRequest res=QueryResolveNameResponse file=rollchain/x/nameservice/keeper/query_server.go
3INF Applied RPC Stub module=nameservice type=tx name=SetServiceName req=MsgSetServiceName res=MsgSetServiceNameResponse file=rollchain/x/nameservice/keeper/msg_server.go
```

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

# TODO: Show image of the Query side entirely

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

sleep 2

rolld q nameservice resolve $(rolld keys show user1 -a) --output=json
```

The expected result should be:

```json
{
  "name": "myname"
}
```