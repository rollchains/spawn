package simapp

import (
	"embed"
)

// !IMPORTANT: interchaintest/ has its own `InterchainTest` embed.FS that will need to be iterated on.

//go:embed .github/* app/* chains/* cmd/* contrib/* scripts/* Makefile Dockerfile proto/*.lock proto/*.yaml *.*
var SimAppFS embed.FS

//go:embed interchaintest/*
var ICTestFS embed.FS

//go:embed proto/*
var ProtoModuleFS embed.FS

//go:embed x/*
var ExtensionFS embed.FS
