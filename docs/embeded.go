package docs

import "embed"

//go:embed versioned_docs/version-v0.50.x/**/*.md
var Docs embed.FS
