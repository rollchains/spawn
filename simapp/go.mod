module github.com/rollchains/spawn/simapp

go 1.22.3

toolchain go1.22.5

// Below are the long-lived replace of the SimApp
replace (
	// use cosmos fork of keyring
	github.com/99designs/keyring => github.com/cosmos/keyring v1.2.0
	// Simapp always use the latest version of the cosmos-sdk
	github.com/cosmos/cosmos-sdk => /home/reece/Desktop/Programming/Go/gordian/gcosmos/_cosmosvendor/cosmos-sdk
	// Fix upstream GHSA-h395-qcrw-5vmq and GHSA-3vp4-m3rf-835h vulnerabilities.
	// TODO Remove it: https://github.com/cosmos/cosmos-sdk/issues/10409
	github.com/gin-gonic/gin => github.com/gin-gonic/gin v1.9.1
	// replace broken goleveldb
	github.com/syndtr/goleveldb => github.com/syndtr/goleveldb v1.0.1-0.20210819022825-2ae1ddf74ef7
)

// SimApp on main always tests the latest extracted SDK modules importing the sdk
replace (
	cosmossdk.io/client/v2 => /home/reece/Desktop/Programming/Go/gordian/gcosmos/_cosmosvendor/cosmos-sdk/client/v2
	cosmossdk.io/collections => /home/reece/Desktop/Programming/Go/gordian/gcosmos/_cosmosvendor/cosmos-sdk/collections
	cosmossdk.io/core => /home/reece/Desktop/Programming/Go/gordian/gcosmos/_cosmosvendor/cosmos-sdk/core
	cosmossdk.io/tools/confix => /home/reece/Desktop/Programming/Go/gordian/gcosmos/_cosmosvendor/cosmos-sdk/tools/confix
	cosmossdk.io/x/accounts => /home/reece/Desktop/Programming/Go/gordian/gcosmos/_cosmosvendor/cosmos-sdk/x/accounts
	cosmossdk.io/x/accounts/defaults/lockup => /home/reece/Desktop/Programming/Go/gordian/gcosmos/_cosmosvendor/cosmos-sdk/x/accounts/defaults/lockup
	cosmossdk.io/x/accounts/defaults/multisig => /home/reece/Desktop/Programming/Go/gordian/gcosmos/_cosmosvendor/cosmos-sdk/x/accounts/defaults/multisig
	cosmossdk.io/x/auth => /home/reece/Desktop/Programming/Go/gordian/gcosmos/_cosmosvendor/cosmos-sdk/x/auth
	cosmossdk.io/x/authz => /home/reece/Desktop/Programming/Go/gordian/gcosmos/_cosmosvendor/cosmos-sdk/x/authz
	cosmossdk.io/x/bank => /home/reece/Desktop/Programming/Go/gordian/gcosmos/_cosmosvendor/cosmos-sdk/x/bank
	cosmossdk.io/x/circuit => /home/reece/Desktop/Programming/Go/gordian/gcosmos/_cosmosvendor/cosmos-sdk/x/circuit
	cosmossdk.io/x/consensus => /home/reece/Desktop/Programming/Go/gordian/gcosmos/_cosmosvendor/cosmos-sdk/x/consensus
	cosmossdk.io/x/distribution => /home/reece/Desktop/Programming/Go/gordian/gcosmos/_cosmosvendor/cosmos-sdk/x/distribution
	cosmossdk.io/x/evidence => /home/reece/Desktop/Programming/Go/gordian/gcosmos/_cosmosvendor/cosmos-sdk/x/evidence
	cosmossdk.io/x/feegrant => /home/reece/Desktop/Programming/Go/gordian/gcosmos/_cosmosvendor/cosmos-sdk/x/feegrant
	cosmossdk.io/x/gov => /home/reece/Desktop/Programming/Go/gordian/gcosmos/_cosmosvendor/cosmos-sdk/x/gov
	cosmossdk.io/x/group => /home/reece/Desktop/Programming/Go/gordian/gcosmos/_cosmosvendor/cosmos-sdk/x/group
	cosmossdk.io/x/mint => /home/reece/Desktop/Programming/Go/gordian/gcosmos/_cosmosvendor/cosmos-sdk/x/mint
	cosmossdk.io/x/nft => /home/reece/Desktop/Programming/Go/gordian/gcosmos/_cosmosvendor/cosmos-sdk/x/nft
	cosmossdk.io/x/params => /home/reece/Desktop/Programming/Go/gordian/gcosmos/_cosmosvendor/cosmos-sdk/x/params
	cosmossdk.io/x/protocolpool => /home/reece/Desktop/Programming/Go/gordian/gcosmos/_cosmosvendor/cosmos-sdk/x/protocolpool
	cosmossdk.io/x/slashing => /home/reece/Desktop/Programming/Go/gordian/gcosmos/_cosmosvendor/cosmos-sdk/x/slashing
	cosmossdk.io/x/staking => /home/reece/Desktop/Programming/Go/gordian/gcosmos/_cosmosvendor/cosmos-sdk/x/staking
	cosmossdk.io/x/tx => /home/reece/Desktop/Programming/Go/gordian/gcosmos/_cosmosvendor/cosmos-sdk/x/tx
	cosmossdk.io/x/upgrade => /home/reece/Desktop/Programming/Go/gordian/gcosmos/_cosmosvendor/cosmos-sdk/x/upgrade
)

// server v2 integration
replace (
	cosmossdk.io/api => /home/reece/Desktop/Programming/Go/gordian/gcosmos/_cosmosvendor/cosmos-sdk/api
	cosmossdk.io/core/testing => /home/reece/Desktop/Programming/Go/gordian/gcosmos/_cosmosvendor/cosmos-sdk/core/testing
	cosmossdk.io/runtime/v2 => /home/reece/Desktop/Programming/Go/gordian/gcosmos/_cosmosvendor/cosmos-sdk/runtime/v2
	cosmossdk.io/server/v2 => /home/reece/Desktop/Programming/Go/gordian/gcosmos/_cosmosvendor/cosmos-sdk/server/v2
	cosmossdk.io/server/v2/appmanager => /home/reece/Desktop/Programming/Go/gordian/gcosmos/_cosmosvendor/cosmos-sdk/server/v2/appmanager
	cosmossdk.io/server/v2/cometbft => /home/reece/Desktop/Programming/Go/gordian/gcosmos/_cosmosvendor/cosmos-sdk/server/v2/cometbft
	cosmossdk.io/server/v2/stf => /home/reece/Desktop/Programming/Go/gordian/gcosmos/_cosmosvendor/cosmos-sdk/server/v2/stf
	cosmossdk.io/store => /home/reece/Desktop/Programming/Go/gordian/gcosmos/_cosmosvendor/cosmos-sdk/store
	cosmossdk.io/store/v2 => /home/reece/Desktop/Programming/Go/gordian/gcosmos/_cosmosvendor/cosmos-sdk/store/v2
)

// TODO:
// Gordian Integration
replace (
	github.com/rollchains/gordian => /home/reece/Desktop/Programming/Go/gordian
	github.com/rollchains/gordian/gcosmos => /home/reece/Desktop/Programming/Go/gordian/gcosmos
)

require (
	cosmossdk.io/api v0.7.5
	cosmossdk.io/client/v2 v2.0.0-00010101000000-000000000000
	cosmossdk.io/core v0.12.1-0.20240725072823-6a2d039e1212
	cosmossdk.io/depinject v1.0.0
	cosmossdk.io/log v1.4.0
	cosmossdk.io/math v1.3.0 // indirect
	cosmossdk.io/runtime/v2 v2.0.0-00010101000000-000000000000
	cosmossdk.io/server/v2 v2.0.0-20240718121635-a877e3e8048a
	cosmossdk.io/store/v2 v2.0.0 // indirect
	cosmossdk.io/tools/confix v0.0.0-00010101000000-000000000000
	cosmossdk.io/x/accounts v0.0.0-20240226161501-23359a0b6d91
	cosmossdk.io/x/auth v0.0.0-00010101000000-000000000000
	cosmossdk.io/x/authz v0.0.0-00010101000000-000000000000
	cosmossdk.io/x/bank v0.0.0-20240226161501-23359a0b6d91
	cosmossdk.io/x/circuit v0.0.0-20230613133644-0a778132a60f
	cosmossdk.io/x/consensus v0.0.0-00010101000000-000000000000
	cosmossdk.io/x/distribution v0.0.0-20230925135524-a1bc045b3190
	cosmossdk.io/x/evidence v0.0.0-20230613133644-0a778132a60f
	cosmossdk.io/x/feegrant v0.0.0-20230613133644-0a778132a60f
	cosmossdk.io/x/gov v0.0.0-20231113122742-912390d5fc4a
	cosmossdk.io/x/group v0.0.0-00010101000000-000000000000
	cosmossdk.io/x/mint v0.0.0-00010101000000-000000000000
	cosmossdk.io/x/nft v0.0.0-20230613133644-0a778132a60f
	cosmossdk.io/x/protocolpool v0.0.0-20230925135524-a1bc045b3190
	cosmossdk.io/x/slashing v0.0.0-00010101000000-000000000000
	cosmossdk.io/x/staking v0.0.0-00010101000000-000000000000
	cosmossdk.io/x/upgrade v0.0.0-20230613133644-0a778132a60f
	github.com/cometbft/cometbft v1.0.0-rc1 // indirect
	github.com/cosmos/cosmos-db v1.0.2
	// this version is not used as it is always replaced by the latest Cosmos SDK version
	github.com/cosmos/cosmos-sdk v0.53.0
	github.com/spf13/cobra v1.8.1
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/spf13/viper v1.19.0
	github.com/stretchr/testify v1.9.0
	google.golang.org/protobuf v1.34.2
)

require (
	cosmossdk.io/server/v2/cometbft v0.0.0-00010101000000-000000000000
	github.com/jhump/protoreflect v1.16.0
	github.com/libp2p/go-libp2p v0.35.0
	github.com/rollchains/gordian/gcosmos v0.0.0-00010101000000-000000000000
)

require (
	buf.build/gen/go/cometbft/cometbft/protocolbuffers/go v1.34.2-20240701160653-fedbb9acfd2f.2 // indirect
	buf.build/gen/go/cosmos/gogo-proto/protocolbuffers/go v1.34.2-20240130113600-88ef6483f90f.2 // indirect
	cloud.google.com/go v0.115.0 // indirect
	cloud.google.com/go/auth v0.5.1 // indirect
	cloud.google.com/go/auth/oauth2adapt v0.2.2 // indirect
	cloud.google.com/go/compute/metadata v0.3.0 // indirect
	cloud.google.com/go/iam v1.1.8 // indirect
	cloud.google.com/go/storage v1.42.0 // indirect
	cosmossdk.io/collections v0.4.0 // indirect
	cosmossdk.io/core/testing v0.0.0-00010101000000-000000000000 // indirect
	cosmossdk.io/errors v1.0.1 // indirect
	cosmossdk.io/errors/v2 v2.0.0-20240731132947-df72853b3ca5 // indirect
	cosmossdk.io/schema v0.1.1 // indirect
	cosmossdk.io/server/v2/appmanager v0.0.0-00010101000000-000000000000 // indirect
	cosmossdk.io/server/v2/stf v0.0.0-00010101000000-000000000000 // indirect
	cosmossdk.io/store v1.1.1-0.20240418092142-896cdf1971bc // indirect
	cosmossdk.io/x/accounts/defaults/lockup v0.0.0-20240417181816-5e7aae0db1f5 // indirect
	cosmossdk.io/x/accounts/defaults/multisig v0.0.0-00010101000000-000000000000 // indirect
	cosmossdk.io/x/epochs v0.0.0-20240522060652-a1ae4c3e0337 // indirect
	cosmossdk.io/x/tx v0.13.3 // indirect
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/99designs/go-keychain v0.0.0-20191008050251-8e49817e8af4 // indirect
	github.com/99designs/keyring v1.2.2 // indirect
	github.com/DataDog/datadog-go v4.8.3+incompatible // indirect
	github.com/DataDog/zstd v1.5.5 // indirect
	github.com/Microsoft/go-winio v0.6.1 // indirect
	github.com/aws/aws-sdk-go v1.54.6 // indirect
	github.com/aymanbagabas/go-osc52/v2 v2.0.1 // indirect
	github.com/benbjohnson/clock v1.3.5 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bgentry/go-netrc v0.0.0-20140422174119-9fd32a8b3d3d // indirect
	github.com/bgentry/speakeasy v0.2.0 // indirect
	github.com/bits-and-blooms/bitset v1.13.0 // indirect
	github.com/btcsuite/btcd/btcec/v2 v2.3.3 // indirect
	github.com/bufbuild/protocompile v0.10.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/chzyer/readline v1.5.1 // indirect
	github.com/cockroachdb/apd/v2 v2.0.2 // indirect
	github.com/cockroachdb/errors v1.11.1 // indirect
	github.com/cockroachdb/logtags v0.0.0-20230118201751-21c54148d20b // indirect
	github.com/cockroachdb/pebble v1.1.0 // indirect
	github.com/cockroachdb/redact v1.1.5 // indirect
	github.com/cockroachdb/tokenbucket v0.0.0-20230807174530-cc333fc44b06 // indirect
	github.com/cometbft/cometbft-db v0.12.0 // indirect
	github.com/cometbft/cometbft/api v1.0.0-rc.1 // indirect
	github.com/containerd/cgroups v1.1.0 // indirect
	github.com/coreos/go-systemd/v22 v22.5.0 // indirect
	github.com/cosmos/btcutil v1.0.5 // indirect
	github.com/cosmos/cosmos-proto v1.0.0-beta.5 // indirect
	github.com/cosmos/crypto v0.1.2 // indirect
	github.com/cosmos/go-bip39 v1.0.0 // indirect
	github.com/cosmos/gogogateway v1.2.0 // indirect
	github.com/cosmos/gogoproto v1.6.1-0.20240809124342-d6a57064ada0 // indirect
	github.com/cosmos/iavl v1.2.1-0.20240731145221-594b181f427e // indirect
	github.com/cosmos/ics23/go v0.10.0 // indirect
	github.com/cosmos/ledger-cosmos-go v0.13.3 // indirect
	github.com/creachadair/atomicfile v0.3.4 // indirect
	github.com/creachadair/tomledit v0.0.26 // indirect
	github.com/danieljoos/wincred v1.2.1 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/davidlazar/go-crypto v0.0.0-20200604182044-b73af7476f6c // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.3.0 // indirect
	github.com/dgraph-io/badger/v4 v4.2.0 // indirect
	github.com/dgraph-io/ristretto v0.1.1 // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/dvsekhvalnov/jose2go v1.6.0 // indirect
	github.com/elastic/gosigar v0.14.3 // indirect
	github.com/emicklei/dot v1.6.2 // indirect
	github.com/fatih/color v1.17.0 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/flynn/noise v1.1.0 // indirect
	github.com/francoispqt/gojay v1.2.13 // indirect
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/getsentry/sentry-go v0.27.0 // indirect
	github.com/go-kit/kit v0.13.0 // indirect
	github.com/go-kit/log v0.2.1 // indirect
	github.com/go-logfmt/logfmt v0.6.0 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-task/slim-sprig v0.0.0-20230315185526-52ccab3ef572 // indirect
	github.com/godbus/dbus v0.0.0-20190726142602-4481cbc300e2 // indirect
	github.com/godbus/dbus/v5 v5.1.0 // indirect
	github.com/gofrs/uuid v4.4.0+incompatible // indirect
	github.com/gogo/googleapis v1.4.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/glog v1.2.1 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/mock v1.6.0 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/btree v1.1.2 // indirect
	github.com/google/flatbuffers v2.0.8+incompatible // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/google/gopacket v1.1.19 // indirect
	github.com/google/orderedcode v0.0.1 // indirect
	github.com/google/pprof v0.0.0-20240416155748-26353dc0451f // indirect
	github.com/google/s2a-go v0.1.7 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.3.2 // indirect
	github.com/googleapis/gax-go/v2 v2.12.5 // indirect
	github.com/gorilla/handlers v1.5.2 // indirect
	github.com/gorilla/mux v1.8.1 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.4.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.16.0 // indirect
	github.com/gsterjov/go-libsecret v0.0.0-20161001094733-a6f4afe4910c // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-getter v1.7.5 // indirect
	github.com/hashicorp/go-hclog v1.6.3 // indirect
	github.com/hashicorp/go-immutable-radix v1.3.1 // indirect
	github.com/hashicorp/go-metrics v0.5.3 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-plugin v1.6.1 // indirect
	github.com/hashicorp/go-safetemp v1.0.0 // indirect
	github.com/hashicorp/go-version v1.7.0 // indirect
	github.com/hashicorp/golang-lru v1.0.2 // indirect
	github.com/hashicorp/golang-lru/v2 v2.0.7 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/hashicorp/yamux v0.1.1 // indirect
	github.com/hdevalence/ed25519consensus v0.2.0 // indirect
	github.com/huandu/skiplist v1.2.0 // indirect
	github.com/huin/goupnp v1.3.0 // indirect
	github.com/iancoleman/strcase v0.3.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/ipfs/boxo v0.19.0 // indirect
	github.com/ipfs/go-cid v0.4.1 // indirect
	github.com/ipfs/go-datastore v0.6.0 // indirect
	github.com/ipfs/go-log v1.0.5 // indirect
	github.com/ipfs/go-log/v2 v2.5.1 // indirect
	github.com/ipld/go-ipld-prime v0.21.0 // indirect
	github.com/jackpal/go-nat-pmp v1.0.2 // indirect
	github.com/jbenet/go-temp-err-catcher v0.1.0 // indirect
	github.com/jbenet/goprocess v0.1.4 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/jmhodges/levigo v1.0.0 // indirect
	github.com/klauspost/compress v1.17.9 // indirect
	github.com/klauspost/cpuid/v2 v2.2.7 // indirect
	github.com/koron/go-ssdp v0.0.4 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/libp2p/go-buffer-pool v0.1.0 // indirect
	github.com/libp2p/go-cidranger v1.1.0 // indirect
	github.com/libp2p/go-flow-metrics v0.1.0 // indirect
	github.com/libp2p/go-libp2p-asn-util v0.4.1 // indirect
	github.com/libp2p/go-libp2p-kad-dht v0.25.2 // indirect
	github.com/libp2p/go-libp2p-kbucket v0.6.3 // indirect
	github.com/libp2p/go-libp2p-pubsub v0.11.0 // indirect
	github.com/libp2p/go-libp2p-record v0.2.0 // indirect
	github.com/libp2p/go-libp2p-routing-helpers v0.7.3 // indirect
	github.com/libp2p/go-msgio v0.3.0 // indirect
	github.com/libp2p/go-nat v0.2.0 // indirect
	github.com/libp2p/go-netroute v0.2.1 // indirect
	github.com/libp2p/go-reuseport v0.4.0 // indirect
	github.com/libp2p/go-yamux/v4 v4.0.1 // indirect
	github.com/linxGnu/grocksdb v1.8.14 // indirect
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/manifoldco/promptui v0.9.0 // indirect
	github.com/marten-seemann/tcp v0.0.0-20210406111302-dfbc87cc63fd // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.14 // indirect
	github.com/mattn/go-sqlite3 v1.14.22 // indirect
	github.com/mdp/qrterminal/v3 v3.2.0 // indirect
	github.com/miekg/dns v1.1.58 // indirect
	github.com/mikioh/tcpinfo v0.0.0-20190314235526-30a79bb1804b // indirect
	github.com/mikioh/tcpopt v0.0.0-20190314235656-172688c1accc // indirect
	github.com/minio/highwayhash v1.0.2 // indirect
	github.com/minio/sha256-simd v1.0.1 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/go-testing-interface v1.14.1 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/mr-tron/base58 v1.2.0 // indirect
	github.com/mtibben/percent v0.2.1 // indirect
	github.com/muesli/termenv v0.15.2 // indirect
	github.com/multiformats/go-base32 v0.1.0 // indirect
	github.com/multiformats/go-base36 v0.2.0 // indirect
	github.com/multiformats/go-multiaddr v0.12.4 // indirect
	github.com/multiformats/go-multiaddr-dns v0.3.1 // indirect
	github.com/multiformats/go-multiaddr-fmt v0.1.0 // indirect
	github.com/multiformats/go-multibase v0.2.0 // indirect
	github.com/multiformats/go-multicodec v0.9.0 // indirect
	github.com/multiformats/go-multihash v0.2.3 // indirect
	github.com/multiformats/go-multistream v0.5.0 // indirect
	github.com/multiformats/go-varint v0.0.7 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/oasisprotocol/curve25519-voi v0.0.0-20230904125328-1f23a7beb09a // indirect
	github.com/oklog/run v1.1.0 // indirect
	github.com/onsi/ginkgo/v2 v2.17.1 // indirect
	github.com/opencontainers/runtime-spec v1.2.0 // indirect
	github.com/opentracing/opentracing-go v1.2.0 // indirect
	github.com/pbnjay/memory v0.0.0-20210728143218-7b4eea64cf58 // indirect
	github.com/pelletier/go-toml/v2 v2.2.2 // indirect
	github.com/petermattis/goid v0.0.0-20240327183114-c42a807a84ba // indirect
	github.com/pion/datachannel v1.5.6 // indirect
	github.com/pion/dtls/v2 v2.2.11 // indirect
	github.com/pion/ice/v2 v2.3.24 // indirect
	github.com/pion/interceptor v0.1.29 // indirect
	github.com/pion/logging v0.2.2 // indirect
	github.com/pion/mdns v0.0.12 // indirect
	github.com/pion/randutil v0.1.0 // indirect
	github.com/pion/rtcp v1.2.14 // indirect
	github.com/pion/rtp v1.8.6 // indirect
	github.com/pion/sctp v1.8.16 // indirect
	github.com/pion/sdp/v3 v3.0.9 // indirect
	github.com/pion/srtp/v2 v2.0.18 // indirect
	github.com/pion/stun v0.6.1 // indirect
	github.com/pion/transport/v2 v2.2.5 // indirect
	github.com/pion/turn/v2 v2.1.6 // indirect
	github.com/pion/webrtc/v3 v3.2.40 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/polydawn/refmt v0.89.0 // indirect
	github.com/prometheus/client_golang v1.19.1 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common v0.55.0 // indirect
	github.com/prometheus/procfs v0.15.1 // indirect
	github.com/quic-go/qpack v0.4.0 // indirect
	github.com/quic-go/quic-go v0.44.0 // indirect
	github.com/quic-go/webtransport-go v0.8.0 // indirect
	github.com/raulk/go-watchdog v1.3.0 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20201227073835-cf1acfcdf475 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/rogpeppe/go-internal v1.12.0 // indirect
	github.com/rollchains/gordian v0.0.0-00010101000000-000000000000 // indirect
	github.com/rs/cors v1.11.0 // indirect
	github.com/rs/zerolog v1.33.0 // indirect
	github.com/sagikazarmark/locafero v0.4.0 // indirect
	github.com/sagikazarmark/slog-shim v0.1.0 // indirect
	github.com/sasha-s/go-deadlock v0.3.1 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spaolacci/murmur3 v1.1.0 // indirect
	github.com/spf13/afero v1.11.0 // indirect
	github.com/spf13/cast v1.7.0 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	github.com/supranational/blst v0.3.12 // indirect
	github.com/syndtr/goleveldb v1.0.1-0.20220721030215-126854af5e6d // indirect
	github.com/tendermint/go-amino v0.16.0 // indirect
	github.com/tidwall/btree v1.7.0 // indirect
	github.com/ulikunitz/xz v0.5.12 // indirect
	github.com/whyrusleeping/go-keyspace v0.0.0-20160322163242-5b898ac5add1 // indirect
	github.com/zondax/hid v0.9.2 // indirect
	github.com/zondax/ledger-go v0.14.3 // indirect
	gitlab.com/yawning/secp256k1-voi v0.0.0-20230925100816-f2616030848b // indirect
	gitlab.com/yawning/tuplehash v0.0.0-20230713102510-df83abbf9a02 // indirect
	go.etcd.io/bbolt v1.4.0-alpha.0.0.20240404170359-43604f3112c5 // indirect
	go.opencensus.io v0.24.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.52.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.52.0 // indirect
	go.opentelemetry.io/otel v1.27.0 // indirect
	go.opentelemetry.io/otel/metric v1.27.0 // indirect
	go.opentelemetry.io/otel/trace v1.27.0 // indirect
	go.uber.org/dig v1.17.1 // indirect
	go.uber.org/fx v1.21.1 // indirect
	go.uber.org/mock v0.4.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	golang.org/x/crypto v0.26.0 // indirect
	golang.org/x/exp v0.0.0-20240531132922-fd00a4e0eefc // indirect
	golang.org/x/mod v0.19.0 // indirect
	golang.org/x/net v0.28.0 // indirect
	golang.org/x/oauth2 v0.21.0 // indirect
	golang.org/x/sync v0.8.0 // indirect
	golang.org/x/sys v0.24.0 // indirect
	golang.org/x/term v0.23.0 // indirect
	golang.org/x/text v0.17.0 // indirect
	golang.org/x/time v0.5.0 // indirect
	golang.org/x/tools v0.23.0 // indirect
	gonum.org/v1/gonum v0.15.0 // indirect
	google.golang.org/api v0.185.0 // indirect
	google.golang.org/genproto v0.0.0-20240617180043-68d350f18fd4 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20240617180043-68d350f18fd4 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240808171019-573a1156607a // indirect
	google.golang.org/grpc v1.65.0 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	gotest.tools/v3 v3.5.1 // indirect
	lukechampine.com/blake3 v1.2.2 // indirect
	pgregory.net/rapid v1.1.0 // indirect
	rsc.io/qr v0.2.0 // indirect
	sigs.k8s.io/yaml v1.4.0 // indirect
)

// Here are the short-lived replace from the SimApp
// Replace here are pending PRs, or version to be tagged
// replace (
// 	<temporary replace>
// )
