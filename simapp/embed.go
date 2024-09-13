package simapp

import (
	"embed"
)

// !IMPORTANT: interchaintest/ has its own `InterchainTest` embed.FS that will need to be iterated on.

//go:embed .github/* app/* chains/* cmd/* contrib/* scripts/* Makefile Dockerfile proto/*.* *.*
var SimAppFS embed.FS

// To embed the interchaintest/ directory, rename the go.mod file to `go.mod_`
//
//go:embed interchaintest/*
var ICTestFS embed.FS

//go:embed proto/example/* proto/ibcmiddleware/* proto/ibcmodule/*
var ProtoModuleFS embed.FS

//go:embed x/*
var ExtensionFS embed.FS
