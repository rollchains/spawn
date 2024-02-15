package simapp

import (
	"embed"
)

// !IMPORTANT: interchaintest/ has its own `InterchainTest` embed.FS that will need to be iterated on.

// TODO: x/ and proto/ in the future
//
//go:embed .github/* app/* chains/* cmd/* contrib/* scripts/* Makefile Dockerfile *.*
var SimAppFS embed.FS

//go:embed interchaintest/*
var ICTestFS embed.FS

// We only need to copy over the proto/ here, since the x/ will be generated automatically
//
//go:embed proto/*
var ProtoModule embed.FS
