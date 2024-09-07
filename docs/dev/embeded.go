package docs

import "embed"

//go:embed demo DEPLOYMENTS.md SYSTEM_SETUP.md
var Docs embed.FS
