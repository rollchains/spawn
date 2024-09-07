package docs

import "embed"

//go:embed dev/DEPLOYMENTS.md dev/SYSTEM_SETUP.md tutorials/*.md
var Docs embed.FS
