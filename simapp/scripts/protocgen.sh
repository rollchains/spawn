#!/usr/bin/env bash

set -e

GO_MOD_PACKAGE="github.com/rollchains/spawn/simapp"

echo "Generating gogo proto code"
cd proto
proto_dirs=$(find . -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)

# Docker
protoImageName=ghcr.io/cosmos/proto-builder:0.14.0
# protoImage="$(which docker) run --rm -v $(pwd):/workspace --workdir /workspace $(protoImageName)"


OS_PLATFORM=`awk -F= '/^NAME/{print $2}' /etc/os-release`
OS_PLATFORM=$(echo $OS_PLATFORM | sed 's/\"//g') # remove the \" from the string


USE_DOCKER=true
if [ "$OS_PLATFORM" = "Ubuntu" ]; then
  if ! [ -x "$(command -v buf)" ]; then
    echo "Installing buf as it is not found"
    # https://buf.build/docs/installation
    BIN="/usr/local/bin" && \
    VERSION="1.45.0" && \
    curl -sSL \
    "https://github.com/bufbuild/buf/releases/download/v${VERSION}/buf-$(uname -s)-$(uname -m)" \
    -o "${BIN}/buf" && \
    chmod +x "${BIN}/buf"
  fi

  if ! [ -x "$(command -v protoc-gen-gocosmos)" ]; then
    docker create --name cosmos-proto-builder $protoImageName
    docker cp cosmos-proto-builder:/go/bin/. $(go env GOPATH)/bin/
    docker rm -f cosmos-proto-builder
    # https://github.com/cosmos/cosmos-sdk/blob/e84c0eb86b20dc95be413b21b0da7377a9bbedc6/contrib/devtools/Dockerfile#L30
    go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway@v1.16.0
    go install cosmossdk.io/orm/cmd/protoc-gen-go-cosmos-orm@v1.0.0-beta.3
  fi

  USE_DOCKER=false
fi

# All other platforms will use the proto builder image w/ docker

echo "Using Docker: $USE_DOCKER"

if [ "$USE_DOCKER" = false ]; then
  buf generate
else
  for dir in $proto_dirs; do
    for file in $(find "${dir}" -maxdepth 1 -name '*.proto'); do
      # this regex checks if a proto file has its go_package set to github.com/strangelove-ventures/poa/...
      # gogo proto files SHOULD ONLY be generated if this is false
      # we don't want gogo proto to run for proto files which are natively built for google.golang.org/protobuf
      if grep -q "option go_package" "$file" && grep -H -o -c "option go_package.*$GO_MOD_PACKAGE/api" "$file" | grep -q ':0$'; then
        echo "Generating gogo proto code for $file"
        docker run --rm -v $(pwd):/workspace --workdir /workspace ${protoImageName} buf generate $file
      fi
    done
  done
fi

echo "Generating pulsar proto code"
buf generate --template buf.gen.pulsar.yaml

# cd ..

sleep 1

# make sure GO_MOD_PACKAGE is not empty.
if [ -z "$GO_MOD_PACKAGE" ]; then
  echo "GO_MOD_PACKAGE is empty. Please set the GO_MOD_PACKAGE variable in the Makefile"
  exit 1
fi

([ -z "$GO_MOD_PACKAGE" ] && exit 1) && cp -r $GO_MOD_PACKAGE/* ./
rm -rf github.com

# # Copy files over for dep injection
rm -rf api && mkdir api
custom_modules=$(find . -name 'module' -type d -not -path "./proto/*")

# # get the 1 up directory (so ./cosmos/mint/module becomes ./cosmos/mint)
# # remove the relative path starter from base namespaces. so ./cosmos/mint becomes cosmos/mint
base_namespace=$(echo $custom_modules | sed -e 's|/module||g' | sed -e 's|\./||g')

# # echo "Base namespace: $base_namespace"
for module in $base_namespace; do
  echo " [+] Moving: ./$module to ./api/$module"

  mkdir -p api/$module

  mv $module/* ./api/$module/

  # # incorrect reference to the module for coins
  find api/$module -type f -name '*.go' -exec sed -i -e 's|types "github.com/cosmos/cosmos-sdk/types"|types "cosmossdk.io/api/cosmos/base/v1beta1"|g' {} \;
  find api/$module -type f -name '*.go' -exec sed -i -e 's|types1 "github.com/cosmos/cosmos-sdk/x/bank/types"|types1 "cosmossdk.io/api/cosmos/bank/v1beta1"|g' {} \;

  rm -rf $module
done
