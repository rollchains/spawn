---
title: "Name Service"
sidebar_label: "Set Structure"
sidebar_position: 2
slug: /build/name-service-structure
---

# Set Data Structure

Extend the template module and add how to store and interact with data. Specifically, you need to set and retrieve a name.

### Set Name

Open the `proto/nameservice/v1` directory. Edit `tx.proto` to add the transaction setter message.

```protobuf title="proto/nameservice/v1/tx.proto"

  // SetServiceName allows a user to set their accounts name.
  rpc SetServiceName(MsgSetServiceName) returns (MsgSetServiceNameResponse);
}

// MsgSetServiceName defines the structure for setting a name.
message MsgSetServiceName {
  option (cosmos.msg.v1.signer) = "sender";

  string sender = 1 [(cosmos_proto.scalar) = "cosmos.AddressString"];

  string name = 2;
}

// MsgSetServiceNameResponse is an empty reply.
message MsgSetServiceNameResponse {}
```

<details>
proto/nameservice/v1/tx.proto file
![proto/nameservice/v1/tx.proto file](https://github.com/rollchains/spawn/assets/31943163/73a583e2-9edd-471f-ada6-1010d0dbf072)
</details>


### Get Name

Find `query.proto` and add the following

```protobuf title="proto/nameservice/v1/query.proto"

  // ResolveName allows a user to resolve the name of an account.
  rpc ResolveName(QueryResolveNameRequest) returns (QueryResolveNameResponse) {
    option (google.api.http).get = "/nameservice/v1/name/{wallet}";
  }
}

// QueryResolveNameRequest grabs the name of a wallet.
message QueryResolveNameRequest {
  string wallet = 1;
}

// QueryResolveNameResponse grabs the wallet linked to a name.
message QueryResolveNameResponse {
  string name = 1;
}
```

<details>
proto/nameservice/v1/query.proto
![proto/nameservice/v1/query.proto file](https://github.com/rollchains/spawn/assets/31943163/234a13d7-be62-492d-961c-63e92d7543d9)
</details>


## Generate Code

These .proto file templates will be converted into Golang source code for you to use. Build the Go source code using the command:

```bash
make proto-gen
```

<details>
make proto-gen expected output
![make proto-gen](https://github.com/rollchains/spawn/assets/31943163/c51bf57c-e83a-4004-8041-9b1f3d3a24f4)
</details>



