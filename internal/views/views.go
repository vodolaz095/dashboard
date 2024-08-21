package views

import "embed"

// Views holds templates used for site rendering
//
//go:embed *.html
var Views embed.FS
