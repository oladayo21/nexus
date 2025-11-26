package nexus

import (
	"embed"
	"io/fs"
)

//go:embed web/dist/*
var webFS embed.FS

// WebFS returns the embedded frontend filesystem
func WebFS() (fs.FS, error) {
	return fs.Sub(webFS, "web/dist")
}
