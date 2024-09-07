## IBc Integration WIP

grpcurl -plaintext localhost:9092 list


grpcurl -plaintext localhost:9092 gordian.server.v1.GordianGRPC/GetBlocksWatermark



## Steps
- Get a Tx hash from the node
- Get Relayer to use gRPC client with gordian
- Get Interchaintest to also use the gRPC client
- Connect gordian to a cometbft chain
- Profit


```bash
go build -o gcosmos . && rm -rf ~/.simappv2/data/application.db/ && ./gcosmos start --g-http-addr 127.0.0.1:26657 --g-grpc-addr 127.0.0.1:9092
```
