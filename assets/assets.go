package assets

import "embed"

// Assets holds js and css used for site rendering
//
//go:embed *.css *.js favicon.ico
var Assets embed.FS
