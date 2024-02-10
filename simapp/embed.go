package simapp

import (
	"embed"
)

// !IMPORTANT: interchaintest/ has its own `InterchainTest` embed.FS that will need to be iterated on.

// TODO: x/ and proto/ in the future
//
//go:embed .github/* app/* chains/* cmd/* configs/* contrib/* scripts/* Makefile Dockerfile *.*
var SimAppFS embed.FS

//go:embed interchaintest/*
var ICTestFS embed.FS
