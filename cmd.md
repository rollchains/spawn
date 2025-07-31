
```bash
# docci run docs/versioned_docs/version-v0.50.x/01-setup/02-install-spawn.md,docs/versioned_docs/version-v0.50.x/02-build-your-application/01-nameservice.md --hide-background-logs

# TODO: we have to support JSON files to have an array of relative file paths in the area

# generate chain (outside of the rollchain dir since it is not generated yet)
docci run docs/versioned_docs/version-v0.50.x/02-build-your-application/01-nameservice.md --hide-background-logs

# interaction (02 - 05)
docci run docs/versioned_docs/version-v0.50.x/02-build-your-application/docci_config.json --hide-background-logs --working-dir rollchain

# cleanup
rm -rf rollchain
```
