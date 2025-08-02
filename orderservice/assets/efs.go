package assets

import "embed"

//go:embed "migrations" "config.yml"
var EmbeddedFiles embed.FS
