module github.com/rollchains/spawn/simapp

go 1.23.1

// Below are the long-lived replace of the SimApp
replace (
	// use cosmos fork of keyring
	github.com/99designs/keyring => github.com/cosmos/keyring v1.2.0
	// Simapp always use the latest version of the cosmos-sdk
	github.com/cosmos/cosmos-sdk => /home/reece/Desktop/Programming/Go/gcosmos/_cosmosvendor/cosmos-sdk
	// Fix upstream GHSA-h395-qcrw-5vmq and GHSA-3vp4-m3rf-835h vulnerabilities.
	// TODO Remove it: https://github.com/cosmos/cosmos-sdk/issues/10409
	github.com/gin-gonic/gin => github.com/gin-gonic/gin v1.9.1
	// replace broken goleveldb
	github.com/syndtr/goleveldb => github.com/syndtr/goleveldb v1.0.1-0.20210819022825-2ae1ddf74ef7
)

//TODO: remove everything below after tags are created
replace (
	cosmossdk.io/api => cosmossdk.io/api v0.7.3-0.20240815194237-858ec2fcb897 // main
	cosmossdk.io/client/v2 => cosmossdk.io/client/v2 v2.0.0-20240905174638-8ce77cbb2450
	cosmossdk.io/core => cosmossdk.io/core v0.12.1-0.20240906083041-6033330182c7 // main
	cosmossdk.io/store => cosmossdk.io/store v1.0.0-rc.0.0.20240906090851-36d9b25e8981 // main
	cosmossdk.io/x/accounts => cosmossdk.io/x/accounts v0.0.0-20240911130545-9e7848985491
	cosmossdk.io/x/accounts/defaults/lockup => cosmossdk.io/x/accounts/defaults/lockup v0.0.0-20240911130545-9e7848985491
	cosmossdk.io/x/accounts/defaults/multisig => cosmossdk.io/x/accounts/defaults/multisig v0.0.0-20240911130545-9e7848985491
	cosmossdk.io/x/authz => cosmossdk.io/x/authz v0.0.0-20240911130545-9e7848985491
	cosmossdk.io/x/bank => cosmossdk.io/x/bank v0.0.0-20240911130545-9e7848985491
	cosmossdk.io/x/circuit => cosmossdk.io/x/circuit v0.0.0-20240911130545-9e7848985491
	cosmossdk.io/x/consensus => cosmossdk.io/x/consensus v0.0.0-20240911130545-9e7848985491
	cosmossdk.io/x/epochs => cosmossdk.io/x/epochs v0.0.0-20240911130545-9e7848985491
	cosmossdk.io/x/evidence => cosmossdk.io/x/evidence v0.0.0-20240911130545-9e7848985491
	cosmossdk.io/x/feegrant => cosmossdk.io/x/feegrant v0.0.0-20240911130545-9e7848985491
	cosmossdk.io/x/gov => cosmossdk.io/x/gov v0.0.0-20240911130545-9e7848985491
	cosmossdk.io/x/group => cosmossdk.io/x/group v0.0.0-20240911130545-9e7848985491
	cosmossdk.io/x/mint => cosmossdk.io/x/mint v0.0.0-20240909082436-01c0e9ba3581
	cosmossdk.io/x/nft => cosmossdk.io/x/nft v0.0.0-20240911130545-9e7848985491
	cosmossdk.io/x/params => cosmossdk.io/x/params v0.0.0-20240911130545-9e7848985491
	cosmossdk.io/x/protocolpool => cosmossdk.io/x/protocolpool v0.0.0-20240911130545-9e7848985491
	cosmossdk.io/x/slashing => cosmossdk.io/x/slashing v0.0.0-20240911130545-9e7848985491
	cosmossdk.io/x/staking => cosmossdk.io/x/staking v0.0.0-20240911130545-9e7848985491
	cosmossdk.io/x/tx => cosmossdk.io/x/tx v0.13.4-0.20240815194237-858ec2fcb897 // main
	cosmossdk.io/x/upgrade => cosmossdk.io/x/upgrade v0.0.0-20240911130545-9e7848985491
	github.com/cometbft/cometbft => github.com/cometbft/cometbft v1.0.0-rc1.0.20240908111210-ab0be101882f

	// github.com/cosmos/cosmos-sdk => github.com/cosmos/cosmos-sdk v0.52.0-beta.1
	github.com/decred/dcrd/dcrec/secp256k1/v4 => github.com/decred/dcrd/dcrec/secp256k1/v4 v4.0.1
)

// server v2 integration
// replace (
// 	cosmossdk.io/api => /home/reece/Desktop/Programming/Go/gcosmos/_cosmosvendor/cosmos-sdk/api
// 	cosmossdk.io/core/testing => /home/reece/Desktop/Programming/Go/gcosmos/_cosmosvendor/cosmos-sdk/core/testing
// 	cosmossdk.io/runtime/v2 => /home/reece/Desktop/Programming/Go/gcosmos/_cosmosvendor/cosmos-sdk/runtime/v2
// 	cosmossdk.io/server/v2 => /home/reece/Desktop/Programming/Go/gcosmos/_cosmosvendor/cosmos-sdk/server/v2
// 	cosmossdk.io/server/v2/appmanager => /home/reece/Desktop/Programming/Go/gcosmos/_cosmosvendor/cosmos-sdk/server/v2/appmanager
// 	cosmossdk.io/server/v2/cometbft => /home/reece/Desktop/Programming/Go/gcosmos/_cosmosvendor/cosmos-sdk/server/v2/cometbft
// 	cosmossdk.io/server/v2/stf => /home/reece/Desktop/Programming/Go/gcosmos/_cosmosvendor/cosmos-sdk/server/v2/stf
// 	cosmossdk.io/store => /home/reece/Desktop/Programming/Go/gcosmos/_cosmosvendor/cosmos-sdk/store
// 	cosmossdk.io/store/v2 => /home/reece/Desktop/Programming/Go/gcosmos/_cosmosvendor/cosmos-sdk/store/v2
// )

// TODO:
// Gordian Integration
replace (
	github.com/rollchains/gordian => /home/reece/Desktop/Programming/Go/gordian
// github.com/rollchains/gordian/gcosmos => /home/reece/Desktop/Programming/Go/gcosmos
)

require (
	cosmossdk.io/api v0.8.0
	cosmossdk.io/client/v2 v2.0.0-beta.5
	cosmossdk.io/collections v0.4.1-0.20240802064046-23fac2f1b8ab
	cosmossdk.io/core v1.0.0
	cosmossdk.io/log v1.4.1
	cosmossdk.io/store v1.1.1
	cosmossdk.io/tools/confix v0.1.2
	cosmossdk.io/x/accounts v0.0.0-20240226161501-23359a0b6d91
	cosmossdk.io/x/authz v0.0.0-00010101000000-000000000000
	cosmossdk.io/x/bank v0.0.0-20240226161501-23359a0b6d91
	cosmossdk.io/x/circuit v0.1.1
	cosmossdk.io/x/consensus v0.0.0-20240911130545-9e7848985491
	cosmossdk.io/x/distribution v0.0.0-20240906090851-36d9b25e8981
	cosmossdk.io/x/evidence v0.1.1
	cosmossdk.io/x/feegrant v0.1.1
	cosmossdk.io/x/gov v0.0.0-20231113122742-912390d5fc4a
	cosmossdk.io/x/group v0.0.0-00010101000000-000000000000
	cosmossdk.io/x/mint v0.0.0-20240909082436-01c0e9ba3581
	cosmossdk.io/x/params v0.0.0-00010101000000-000000000000
	cosmossdk.io/x/protocolpool v0.0.0-20230925135524-a1bc045b3190
	cosmossdk.io/x/slashing v0.0.0-00010101000000-000000000000
	cosmossdk.io/x/staking v0.0.0-00010101000000-000000000000
	cosmossdk.io/x/tx v1.0.0-alpha.1
	cosmossdk.io/x/upgrade v0.1.4
	github.com/cometbft/cometbft v1.0.0-rc1.0.20240908111210-ab0be101882f
	github.com/cometbft/cometbft/api v1.0.0-rc.1.0.20240909080621-90c99d9658b0
	github.com/cosmos/cosmos-db v1.0.3-0.20240911104526-ddc3f09bfc22
	github.com/cosmos/cosmos-sdk v0.53.0
	github.com/cosmos/gogoproto v1.7.0
	github.com/cosmos/ibc-go/v9 v9.0.0
	github.com/spf13/cast v1.7.0
	github.com/spf13/cobra v1.8.1
	github.com/spf13/viper v1.19.0
	github.com/stretchr/testify v1.9.0
)

require (
	buf.build/gen/go/cometbft/cometbft/protocolbuffers/go v1.35.1-20240701160653-fedbb9acfd2f.1 // indirect
	buf.build/gen/go/cosmos/gogo-proto/protocolbuffers/go v1.35.1-20240130113600-88ef6483f90f.1 // indirect
	cosmossdk.io/depinject v1.0.0 // indirect
	cosmossdk.io/errors v1.0.1 // indirect
	cosmossdk.io/errors/v2 v2.0.0-20240731132947-df72853b3ca5 // indirect
	cosmossdk.io/math v1.3.0 // indirect
	cosmossdk.io/store/v2 v2.0.0-20241028185935-7b9a89b5689c // indirect
	cosmossdk.io/x/nft v0.0.0-20241028185935-7b9a89b5689c // indirect
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/DataDog/zstd v1.5.6 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/btcsuite/btcd/btcec/v2 v2.3.4 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/cockroachdb/errors v1.11.3 // indirect
	github.com/cockroachdb/fifo v0.0.0-20240816210425-c5d0cb0b6fc0 // indirect
	github.com/cockroachdb/logtags v0.0.0-20230118201751-21c54148d20b // indirect
	github.com/cockroachdb/pebble v1.1.2 // indirect
	github.com/cockroachdb/redact v1.1.5 // indirect
	github.com/cockroachdb/tokenbucket v0.0.0-20230807174530-cc333fc44b06 // indirect
	github.com/cometbft/cometbft-db v0.15.0 // indirect
	github.com/cosmos/btcutil v1.0.5 // indirect
	github.com/cosmos/cosmos-proto v1.0.0-beta.5 // indirect
	github.com/cosmos/crypto v0.1.2 // indirect
	github.com/cosmos/ics23/go v0.11.0 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.3.0 // indirect
	github.com/dgraph-io/badger/v4 v4.3.0 // indirect
	github.com/dgraph-io/ristretto v0.1.2-0.20240116140435-c67e07994f91 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/getsentry/sentry-go v0.28.1 // indirect
	github.com/go-kit/kit v0.13.0 // indirect
	github.com/go-kit/log v0.2.1 // indirect
	github.com/go-logfmt/logfmt v0.6.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/btree v1.1.3 // indirect
	github.com/google/flatbuffers v24.3.25+incompatible // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.16.0 // indirect
	github.com/hashicorp/go-immutable-radix v1.3.1 // indirect
	github.com/hashicorp/go-metrics v0.5.3 // indirect
	github.com/hashicorp/golang-lru v1.0.2 // indirect
	github.com/hdevalence/ed25519consensus v0.2.0 // indirect
	github.com/iancoleman/strcase v0.3.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jmhodges/levigo v1.0.0 // indirect
	github.com/klauspost/compress v1.17.9 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/linxGnu/grocksdb v1.9.3 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/oasisprotocol/curve25519-voi v0.0.0-20230904125328-1f23a7beb09a // indirect
	github.com/petermattis/goid v0.0.0-20240813172612-4fcff4a6cae7 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/prometheus/client_golang v1.20.5 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common v0.60.1 // indirect
	github.com/prometheus/procfs v0.15.1 // indirect
	github.com/rogpeppe/go-internal v1.12.0 // indirect
	github.com/rs/zerolog v1.33.0 // indirect
	github.com/sasha-s/go-deadlock v0.3.5 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/supranational/blst v0.3.13 // indirect
	github.com/syndtr/goleveldb v1.0.1-0.20220721030215-126854af5e6d // indirect
	github.com/tendermint/go-amino v0.16.0 // indirect
	gitlab.com/yawning/secp256k1-voi v0.0.0-20230925100816-f2616030848b // indirect
	gitlab.com/yawning/tuplehash v0.0.0-20230713102510-df83abbf9a02 // indirect
	go.etcd.io/bbolt v1.4.0-alpha.1 // indirect
	go.opencensus.io v0.24.0 // indirect
	golang.org/x/crypto v0.28.0 // indirect
	golang.org/x/exp v0.0.0-20240909161429-701f63a606c0 // indirect
	golang.org/x/net v0.30.0 // indirect
	golang.org/x/sys v0.26.0 // indirect
	golang.org/x/text v0.19.0 // indirect
	google.golang.org/genproto v0.0.0-20240701130421-f6361c86f094 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20240814211410-ddb44dafa142 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20241021214115-324edc3d5d38 // indirect
	google.golang.org/grpc v1.67.1 // indirect
	google.golang.org/protobuf v1.35.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	sigs.k8s.io/yaml v1.4.0 // indirect
)
