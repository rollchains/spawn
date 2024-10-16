FROM golang:1.22-alpine AS go-builder

SHELL ["/bin/sh", "-ecuxo", "pipefail"]

RUN apk add --no-cache ca-certificates build-base git musl-dev

WORKDIR /code

ARG TARGETARCH=amd64
ARG BUILDARCH=amd64
RUN if [ "${TARGETARCH}" = "arm64" ] && [ "${BUILDARCH}" != "arm64" ]; then\
        wget -c https://musl.cc/aarch64-linux-musl-cross.tgz -O - | tar -xzvv --strip-components 1 -C /usr;\
    elif [ "${TARGETARCH}" = "amd64" ] && [ "${BUILDARCH}" != "amd64" ]; then\
        wget -c https://musl.cc/x86_64-linux-musl-cross.tgz -O - | tar -xzvv --strip-components 1 -C /usr;\
    fi

ADD go.mod go.sum ./
RUN set -eux; \
    export ARCH=$(uname -m); \
    LIBDIR=/lib;\
    if [ "${TARGETARCH}" = "arm64" ]; then\
      export ARCH=aarch64;\
      if [ "${BUILDARCH}" != "arm64" ]; then\
        LIBDIR=/usr/aarch64-linux-musl/lib;\
        mkdir -p $LIBDIR;\
        export CC=aarch64-linux-musl-gcc CXX=aarch64-linux-musl-g++;\
      fi;\
    elif [ "${TARGETARCH}" = "amd64" ]; then\
      export ARCH=x86_64;\
      if [ "${BUILDARCH}" != "amd64" ]; then\
        LIBDIR=/usr/x86_64-linux-musl/lib;\
        mkdir -p $LIBDIR;\
        export CC=x86_64-linux-musl-gcc CXX=x86_64-linux-musl-g++;\
      fi;\
    fi;\
    WASM_VERSION=$(go list -m all | grep github.com/CosmWasm/wasmvm/v2 || true); \
    if [ ! -z "${WASM_VERSION}" ]; then \
      WASMVM_REPO=$(echo $WASM_VERSION | awk '{print $1}');\
      WASMVM_REPO=${WASMVM_REPO%/v2};\
      WASMVM_VERS=$(echo $WASM_VERSION | awk '{print $2}');\
      wget -O $LIBDIR/libwasmvm_muslc.a https://${WASMVM_REPO}/releases/download/${WASMVM_VERS}/libwasmvm_muslc.${ARCH}.a;\
      ln $LIBDIR/libwasmvm_muslc.a $LIBDIR/libwasmvm.x86_64.a;\
      ln $LIBDIR/libwasmvm_muslc.a $LIBDIR/libwasmvm_muslc.x86_64.a;\
      ln $LIBDIR/libwasmvm_muslc.a $LIBDIR/libwasmvm.aarch64.a;\
      ln $LIBDIR/libwasmvm_muslc.a $LIBDIR/libwasmvm_muslc.aarch64.a;\
    fi;\
    go mod download;

# Copy over code
COPY . /code

# force it to use static lib (from above) not standard libgo_cosmwasm.so file
# then log output of file /code/bin/wasmd
# then ensure static linking
RUN LEDGER_ENABLED=false BUILD_TAGS=muslc LINK_STATICALLY=true make build \
  && file /code/build/wasmd \
  && echo "Ensuring binary is statically linked ..." \
  && (file /code/build/wasmd | grep "statically linked")

# --------------------------------------------------------
FROM alpine:3.16

COPY --from=go-builder /code/build/wasmd /usr/bin/wasmd

# Install dependencies used for Starship
RUN apk add --no-cache curl make bash jq sed

WORKDIR /opt

# rest server, tendermint p2p, tendermint rpc
EXPOSE 1317 26656 26657

CMD ["/usr/bin/wasmd", "version"]
