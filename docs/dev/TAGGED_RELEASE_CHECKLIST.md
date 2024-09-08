# Checklist

This document outlines the steps required to verify the functionality of spawn and the associated components.

## Readme

- Run ./scripts/bump_docs.sh and ./scripts/bump_localic.sh to update spawn and local-ic versions in all markdown files, docker, and makefiles.
- Verify the README is up to date
- Ensure examples run in the [../versioned_docs/](../versioned_docs/) for this release as expected.


## Semi-Automatic Verification

- The [matrix generator](../../scripts/matrix_generator.py) builds up a variety of test cases to verify chains build, test, and push to github without issue. Use this to verify that many different consensus, features, bech32s, and denominations generate and work as expected.
- Push some of these chains up to ensure all CI work as expected and end to end passes.

## Manual Verification

- `make template-*` has default chains configured for you to test on with different consensus values.
- Create new modules and add some proto code to be auto generated on the next `make proto-gen`.
- Ensure goreleaser runs for all instances like `goreleaser build --skip-validate --snapshot --clean -f .goreleaser.yaml`
- Validate the explorer
    - `make template-staking && cd myproject && make sh-testnet`
    - `make explorer`


## Future Verifications
- Simulate cobra commands and ensure they work as expected, validate input
- Validate user modules: https://github.com/rollchains/spawn/pull/172
- Local Chain Validation (like the github pushes, but entire local w/ ICT): https://github.com/rollchains/spawn/pull/171
- Run testnets with `make testnet` instead of just `sh-testnet`
- https://github.com/rollchains/spawn/pull/143#discussion_r1619583860
- Run other commands after new chain, before unit test (i.e. add module). Currently github CI does not like `make proto-gen` creation
