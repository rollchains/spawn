---
title: "Name Service"
sidebar_label: "Configure Client"
sidebar_position: 4
slug: /build/name-service-client
---

# Command Line Client

Using the Cosmos-SDKs AutoCLI, you will easily set up the CLI client for transactions and queries.

### Query

Update the autocli to allow someone to get the name of a wallet account.

```go title="x/nameservice/autocli.go"
		Query: &autocliv1.ServiceCommandDescriptor{
            Service: modulev1.Query_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				// highlight-start
				{
					RpcMethod: "ResolveName",
					Use:       "resolve [wallet]",
					Short:     "Resolve the name of a wallet address",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "wallet"},
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
```

<details>
  <summary>AutoCLI Query</summary>

  ![AutoCLI Query](https://github.com/rollchains/spawn/assets/31943163/fefe8c7d-88b5-42d5-afd9-cb33cd22df16)
</details>



### Transaction

Also add interaction in `x/nameservice/autocli.go` to set the name of a wallet account.

```go title="x/nameservice/autocli.go"
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service: modulev1.Msg_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				// highlight-start
				{
					RpcMethod: "SetServiceName",
					Use:       "set [name]",
					Short:     "Set the mapping to your wallet address",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "name"},
					},
				},
				// highlight-end
                {
					// NOTE: this is already included in the current source
					RpcMethod: "UpdateParams",
					Skip:      false,
				},
			},
		},
```

<details>
  <summary>AutoCLI Tx</summary>

  ![AutoCLI Tx](https://github.com/rollchains/spawn/assets/31943163/e945c898-415c-4d22-8bb3-b8af34a44cee)
</details>


