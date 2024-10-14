# Checklist

This document outlines the steps required to verify the functionality of spawn and the associated components.

## Docs

- [ ] ./scripts/bump_docs.sh
- [ ] ./scripts/bump_localic.sh
- [ ] Verify the README is up to date
- [ ] Ensure examples run in the [../versioned_docs/](../versioned_docs/) for this release as expected, across multiple machines
    - [ ] MacOS
    - [ ] Linux
    - [ ] Windows (WSL)


## Semi-Automatic Verification

- [ ] Run the [matrix generator](../../scripts/matrix_generator.py)
    - [ ] Local
    - [ ] Github CI

## Manual Verification

- [ ] `make template-*`: Must pass test and make sh-testnet.
    - [ ] Verify with 1 that make testnet works
- [ ] Ensure goreleaser runs for all instances like `goreleaser build --skip-validate --snapshot --clean -f .goreleaser.yaml`
- [ ] Validate the explorer
    - `make template-staking && cd myproject && make sh-testnet`
    - `make explorer`


## Future Verifications
- Simulate cobra commands and ensure they work as expected, validate input
- Validate user modules: https://github.com/rollchains/spawn/pull/172
- Local Chain Validation (like the github pushes, but entire local w/ ICT): https://github.com/rollchains/spawn/pull/171
- Run testnets with `make testnet` instead of just `sh-testnet`
- https://github.com/rollchains/spawn/pull/143#discussion_r1619583860
- Run other commands after new chain, before unit test (i.e. add module). Currently github CI does not like `make proto-gen` creation
