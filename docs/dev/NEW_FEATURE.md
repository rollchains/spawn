# Adding a new Module

## Steps

1. Modify the simapp/ app.go, go.mod, etc to setup your module.
2. Use [`spawntag`'s](#spawntag) to signal where to remove extra content from in the app.
3. Add your module/feature to `SupportedFeatures` in [cmd/spawn/new_chain.go](../../cmd/spawn/new_chain.go)
4. Add the logic to [spawn/remove_features.go](../../spawn/remove_features.go) to remove your feature from the simapp on generate.

**note**: if your feature has a complex setup, reference removing `wasm` from the app for a good guide.

---

## SpawnTag

Anywhere that takes multiple lines to setup, you set `spawntag`'s to signal to the spawn app builder that this code should be removed if the module is not being used.

The following are supported:
- `// <spawntag:moduleName` - Start of the code block
- `// spawntag:moduleName>` - End of the code block
- `// spawntag:moduleName` - Remove single line
- `// ?spawntag:moduleName` - Uncomment line if module is not used (i.e. line swap)


### Multi-Line Removal

This example shows how to remove multiple lines if the wasm module is not being used. Notice, since the module name is `wasm` and the methods names are `wasm`, we do not have to wrap all code. Just the code that does not reference wasm directly. i.e. `if err != nil {` mentions no where it is for the wasmConfig err.

```go title="app.go"
wasmDir := filepath.Join(homePath, "wasm")
wasmConfig, err := wasm.ReadWasmConfig(appOpts)

// <spawntag:wasm
if err != nil {
    panic(fmt.Sprintf("error while reading wasm config: %s", err))
}
// spawntag:wasm>
```

### Single Line Removal

Odd namespaces or designs can cause the app to not be able to find content to automatically remove. Since the module name is `tokenfactory` it does not know to also look for `token_factory`. This edge case is easier to cover by using a spawntag to remove the line. Just make sure the line is on it's own and formatted in a similar design style. On the first save of the file, it will format properly.

```go title="app.go"

capabilities = strings.Join(
    []string{
        "iterator", "staking", "stargate",
        "cosmwasm_1_1", "cosmwasm_1_2", "cosmwasm_1_3", "cosmwasm_1_4",
        "token_factory", // spawntag:tokenfactory
    }, ",")
```

### Line Swap

If you wish to swap in a line (due to a fork of logic), you can use the `?` spawntag. This will uncomment the line if the module is not being used.

If the user does not select to use globalfee, the `globalfeeante` line would be removed as expected. However the `globalfeeante.NewFeeDecorator` is a fork of `ante.NewDeductFeeDecorator`. So we can't leave it as it, and use the ?spawntag: line to default to the base logic.

```go title="ante.go"
    // simapp ante.go

    globalfeeante.NewFeeDecorator(options.BypassMinFeeMsgTypes, options.GlobalFeeKeeper, options.StakingKeeper, 2_000_000),
    //ante.NewDeductFeeDecorator(options.AccountKeeper, options.BankKeeper, options.FeegrantKeeper, options.TxFeeChecker), // ?spawntag:globalfee
```
