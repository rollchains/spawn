#!/bin/bash

# Spawn a new instance
make install && spawn new my-project --debug --bech32=cosmos --bin=appd --disable=tokenfactory

# Specific instance with e2e tests
make install && spawn new rollchains --no-git --bin=rolld --bech32=roll --denom=uroll --disable=globalfee,poa
cd rollchains && make local-image && make ictest-basic