package simapp

import (
	"embed"
)

// !IMPORTANT: interchaintest/ has its own `InterchainTest` embed.FS that will need to be iterated on.

// TODO: proto/*.* *.*
//
//go:embed .github/* app/* chains/* scripts/* Makefile Dockerfile nginx/* contrib/* go.mod go.sum *.go
var SimAppFS embed.FS

// To embed the interchaintest/ directory, rename the go.mod file to `go.mod_`
//
// //go:embed interchaintest/*
var ICTestFS embed.FS

// //go:embed proto/example/* proto/ibcmiddleware/*
var ProtoModuleFS embed.FS

// //go:embed x/*
var ExtensionFS embed.FS
