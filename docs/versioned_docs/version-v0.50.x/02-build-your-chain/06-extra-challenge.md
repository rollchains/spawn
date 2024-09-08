---
title: "Name Service"
sidebar_label: "Bonus Challenge"
sidebar_position: 6
slug: /build/name-service-bonus
---

# Extra Challenge

It seems the nameservice will let you set any name length you want. Add a validation check in `SetServiceName` to ensure the name is less than 32 characters long.

<details>
<summary>Hint #1</summary>
<p>The `SetServiceName` in the msg_server.go looks like an interesting place to start. It should return an error if the name is too long.</p>
</details>

<details>
<summary>Solution</summary>

If a user attempts to submit a name longer than 32 characters, it will return an error that is not allowed.
```go
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
