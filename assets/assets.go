// Package assets provides access to embedded assets.
package assets

import "embed"

//go:embed translations/*.json
var Translations embed.FS
