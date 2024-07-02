# Checklist

This document outlines the steps required to verify the functionality of spawn and the associated components.

## Readme

- Update the README document to be the next spawn tag version (i.e. v0.50.X) that the next release will be.
- Verify the README running docs work as expected for the example chain
- Verify docs/demo works

## Semi-Automatic Verification

- The [matrix generator](../scripts/matrix_generator.py) builds up a variety of test cases to verify chains build, test, and push to github without issue. Use this to verify that many different consensus, features, bech32s, and denominations generate and work as expected.
- Push some of these chains up to ensure all CI work as expected and end to end passes.

## Manual Verification

- `make template-*` has default chains configured for you to test on with different consensus values.
- Create new modules and add some proto code to be auto generated on the next `make proto-gen`.
- Ensure goreleaser runs for all instances like `goreleaser build --skip-validate --snapshot --clean -f .goreleaser.yaml`