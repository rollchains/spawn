package simapp

import (
	"embed"
)

//go:embed **/* go.mod go.sum Makefile chains.yaml
var SimApp embed.FS
