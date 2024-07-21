package handler

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
)

//go:embed build/*
var staticFiles embed.FS

var staticFS = fs.FS(staticFiles)

// StaticFilesHandler serves files from the /build directory inside embed.FS. It also serves the /build/index.html file
func StaticFilesHandler() http.Handler {
	staticBuildFS, err := fs.Sub(staticFS, "build")
	if err != nil {
		log.Fatal(err)
	}

	return http.FileServer(http.FS(staticBuildFS))
}

// StaticHtmlHandler serves the /build/<name>.html file from embed.FS
func StaticHtmlHandler(name string) http.Handler {
	debugHtmlFile, err := fs.ReadFile(staticFS, "build/"+name+".html")
	if err != nil {
		log.Fatal("StaticHtmlHandler() error: ", err)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(debugHtmlFile)
	})
}
