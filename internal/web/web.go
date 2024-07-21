package web

import (
	"embed"
	"io/fs"
)

//go:embed build/*
var staticFiles embed.FS

var StaticFS = fs.FS(staticFiles)
