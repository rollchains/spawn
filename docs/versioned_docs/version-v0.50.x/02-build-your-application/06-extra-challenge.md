---
title: "Name Service"
sidebar_label: "Bonus"
sidebar_position: 6
slug: /build/name-service-bonus
---

# Extra Challenges

## Challenge 1: Limit Input

It seems the nameservice will let you set any name length you want. Add a validation check in `SetServiceName` to ensure the name is less than 32 characters long.

<details>
<summary>Hint #1</summary>
<p>The `SetServiceName` in the msg_server.go looks like an interesting place to start. It should return an error if the name is too long.</p>
</details>

<details>
<summary>Solution</summary>

If a user attempts to submit a name longer than 32 characters, it will return an error that is not allowed.
```go title="x/nameservice/keeper/msg_server.go"
// SetServiceName implements types.MsgServer.
func (ms msgServer) SetServiceName(ctx context.Context, msg *types.MsgSetServiceName) (*types.MsgSetServiceNameResponse, error) {
	if len(msg.Name) > 32 {
		return nil, fmt.Errorf("name cannot be longer than 32 characters")
	}

	if err := ms.k.NameMapping.Set(ctx, msg.Sender, msg.Name); err != nil {
		return nil, err
	}

	return &types.MsgSetServiceNameResponse{}, nil
}
```
</details>


## Challenge 2: Resolve Wallet From Name

Currently the nameservice only allows you to resolve a name given a wallet. If someone has a name they should be able to resolve the wallet address. Add a new query to the `query_server` and autocli client to resolve a wallet address from a name.

> This challenge is signinicantly harder and will some previous Go programming knowledge with iterators. You can also just copy the solutions.

<details>
<summary>Hint #1</summary>
<p>Create a new query.proto for ResolveWallet that takes in a name string</p>
</details>

<details>
<summary>Solution #1</summary>

```protobuf title="proto/nameservice/v1/query.proto"
// Query provides defines the gRPC querier service.
service Query {
  // Params queries all parameters of the module.
  rpc Params(QueryParamsRequest) returns (QueryParamsResponse) {
    option (google.api.http).get = "/nameservice/v1/params";
  }

  // ResolveName allows a user to resolve the name of an account.
  rpc ResolveName(QueryResolveNameRequest) returns (QueryResolveNameResponse) {
    option (google.api.http).get = "/nameservice/v1/name/{wallet}";
  }

  // highlight-start
  // ResolveWallet allows a user to resolve the wallet of a name.
  rpc ResolveWallet(QueryResolveWalletRequest) returns (QueryResolveWalletResponse) {
    option (google.api.http).get = "/nameservice/v1/wallet/{name}";
  }
  // highlight-end
}

// highlight-start
message QueryResolveWalletRequest {
  string name = 1;
}

message QueryResolveWalletResponse {
  string wallet = 1;
}
// highlight-end
```

```bash
make proto-gen
```

</details>

<details>
<summary>Hint #2</summary>
<p>Iterate through the `k.Keeper.NameMapping`, check the Value(). if it matches the name requested, return that wallet (Key)</p>
</details>

<details>
<summary>Solution #2</summary>

```go title="x/nameservice/keeper/query_server.go"
// ResolveWallet implements types.QueryServer.
func (k Querier) ResolveWallet(goCtx context.Context, req *types.QueryResolveWalletRequest) (*types.QueryResolveWalletResponse, error) {
	// highlight-start
	// create a way to iterate over all the name mappings.
	iter, err := k.Keeper.NameMapping.Iterate(goCtx, nil)
	if err != nil {
		return nil, err
	}
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		// get the value (name)
		v, err := iter.Value()
		if err != nil {
			return nil, err
		}

		// if current name matches the requested name,
		// return the wallet address for the name
		if v == req.Name {
			walletAddr, err := iter.Key()
			if err != nil {
				return nil, err
			}

			return &types.QueryResolveWalletResponse{
				Wallet: walletAddr,
			}, nil
		}
	}

	return nil, fmt.Errorf("wallet not found for name %s", req.Name)
	// highlight-end
}


```
This is not the most efficient way to do this. If you would like, create a new WalletMapping collection that maps name->sender when `SetServiceName` is called. This way you can resolve the wallet from the name in O(1) time (i.e. instant) instead of looping through all possible wallets.

</details>


<details>
<summary>Hint #3</summary>
<p>Add the AutoCLI method to `ResolveWallet` with the `ProtoField` "name" to match the .proto file</p>
</details>


<details>
<summary>Solution #3</summary>

```go title="x/nameservice/autocli.go"
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: modulev1.Query_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "ResolveName",
					Use:       "resolve [wallet]",
					Short:     "Resolve the name of a wallet address",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "wallet"},
					},
				},
				// highlight-start
				{
					RpcMethod: "ResolveWallet",
					Use:       "wallet [name]",
					Short:     "Resolve the wallet address from a given name",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "name"},
					},
				},
				// highlight-end
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "Query the current module parameters",
				},
			},
		},
		...
```

Then `make install` and re-run the testnet to verify `rolld q nameservice wallet <name>` returns the expected wallet address.

</details>
