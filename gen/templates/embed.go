package templates

import "embed"

//go:embed object.go.tpl
var embFS embed.FS

//go:embed object.dart.tpl
var embFSdart embed.FS
