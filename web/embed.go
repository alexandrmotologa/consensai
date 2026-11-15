package web

import "embed"

// DistFS contains the embedded compiled React 19 web cockpit assets.
//
//go:embed all:dist
var DistFS embed.FS
