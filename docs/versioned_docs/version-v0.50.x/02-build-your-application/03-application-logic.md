---
title: "Name Service"
sidebar_label: "Application Logic"
sidebar_position: 3
slug: /build/name-service-application
---

# Save Storage Structure

You now need to set the data structure in the keeper to store the wallet to name pair. Keepers are where the data is stored for future use and handle the business logic.
Think of it like a box where you store your data and the methods used to modify this data. It is self contained and only gets further access if you allow the keeper to do so (shown in [part 2](../02-build-your-application/08-ibc-module.md)).

```go title="x/nameservice/keeper/keeper.go"

type Keeper struct {
	...
	// highlight-next-line
	NameMapping collections.Map[string, string]
}

...

func NewKeeper() Keeper {
  ...

  k := Keeper{
    ...
	// highlight-start
    NameMapping: collections.NewMap(sb, collections.NewPrefix(1),
		"name_mapping", collections.StringKey, collections.StringValue,
	),
	// highlight-end
  }

}
```

---

## Application Logic

Update the msg_server logic to set the name upon request from a user.

```go title="x/nameservice/keeper/msg_server.go"
func (ms msgServer) SetServiceName(ctx context.Context, msg *types.MsgSetServiceName) (*types.MsgSetServiceNameResponse, error) {
	// highlight-start
	if err := ms.k.NameMapping.Set(ctx, msg.Sender, msg.Name); err != nil {
		return nil, err
	}

	return &types.MsgSetServiceNameResponse{}, nil
	// highlight-end
}
```

and also for the query_server to retrieve the name.

```go title="x/nameservice/keeper/query_server.go"
func (k Querier) ResolveName(goCtx context.Context, req *types.QueryResolveNameRequest) (*types.QueryResolveNameResponse, error) {
	// highlight-start
	v, err := k.Keeper.NameMapping.Get(goCtx, req.Wallet)
	if err != nil {
		return nil, err
	}

	return &types.QueryResolveNameResponse{
		Name: v,
	}, nil
	// highlight-end
}
```

---
