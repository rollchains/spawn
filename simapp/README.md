# Generated With [Spawn](https://github.com/rollchains/spawn)

## Module Scaffolding

- `spawn module new <name>` *Generates a Cosmos module template*

## Content Generation

- `make proto-gen` *Generates go code from proto files, stubs interfaces*

## Testnet

- `make testnet` *IBC testnet from chain <-> local cosmos-hub*
- `make sh-testnet` *Single node, no IBC. quick iteration*
- `local-ic chains` *See available testnets from the chains/ directory*
- `local-ic start <name>` *Starts a local chain with the given name*

## Local Images

- `make install`      *Builds the chain's binary*
- `make local-image`  *Builds the chain's docker image*

## Testing

- `go test ./... -v` *Unit test*
- `make ictest-*`  *E2E testing*

## Webapp Template

Generate the template base with spawn. Requires [npm](https://nodejs.org/en/download/package-manager) and [yarn](https://classic.yarnpkg.com/lang/en/docs/install) to be installed.

- `make generate-webapp` *[Cosmology Webapp Template](https://github.com/cosmology-tech/create-cosmos-app)*

Start the testnet with `make testnet`, and open the webapp `cd ./web && yarn dev`

## Typescript Client

<https://cosmology.zone/learn/telescope/overview-of-telescope>

```bash

# Install Telescope

telescope generate --name chain-js # TODO: how do I default set other information? like the __CHAINNAME__ ?
    ```
    ? [__CHAINNAME__] Enter chain name in all lowercase, e.g. osmosis localchain
    ? [__USERFULLNAME__] Enter author full name reece williams
    ? [__USEREMAIL__] Enter author email reecepbcups@gmail.com
    ? [__MODULENAME__] Enter the module name chain-js
    ? [__MODULEDESC__] Enter the module description desc of chain-js
    ? [__USERNAME__] Enter your github username reecepbcups
    ? [__ACCESS__] Module access? public
    ? [__LICENSE__] Which license? MIT
    ? [scoped] use npm scopes? Yes
    ```

cd chain-js

telescope install @protobufs/tendermint @protobufs/ibc @protobufs/google @protobufs/gogoproto @protobufs/cosmos_proto @protobufs/cosmos

cp -r ../proto/* ./proto # I dislike this step, why? Can't I use the parent package for my custom modules, but then relative for the installed ones?

telescope transpile --protoDirs=./proto --outPath=./src/codegen --config ../.telescope.json

yarn add @cosmology/lcd

# install package
yarn

# publish
npm publish
```


## Connect to a TS locally without publish

```bash
cd chain-js
yarn link

mkdir t # testing
cd t

npm i typescript --save-dev
npx tsc --init

touch index.ts

    ```ts
    import {nameservice} from '@reecepbcups/chain-js';

    console.log(`Hello world`);

    const client = nameservice.ClientFactory.createRPCQueryClient({ rpcEndpoint: 'http://localhost:26657' });

    client.then((client) => {
        client.nameservice.v1.params().then((res) => {
            console.log(res);
        });
    });
    ```

yarn link "@reecepbcups/chain-js" # TODO: just simlink directly with relative paths

npx ts-node index.ts

```
