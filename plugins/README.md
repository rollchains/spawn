# Plugins

Spawn plugins allow you to build custom functionality on top of Spawn. To accomplish this, you build a cobra CLI binary using the `github.com/rollchains/spawn` import, then build a bianry off of it. Saving this to `$HOME/.spawn/plugins` will allow you to use the binary as a plugin with spawn. Opening the opertunity to closed source plugins and add on features across the stack.

## Getting Started

Reference the [example spawn plugin](./example/example-plugin.go) to get started.

## Running a Plugin

Note that to use flags, you must use a `--` before flags for the child command context. Flags before the `--` apply to the root of the plugin command.

- `spawn plugin <name> [arguments] -- [--flags]`