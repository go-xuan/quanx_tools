package template

import "embed"

//go:embed *
var FS embed.FS

const (
	Suffix = ".tmpl"
)
