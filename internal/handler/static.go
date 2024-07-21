package handler

import (
	"github.com/julien-wff/cesi-dossier-synthese/internal/web"
	"io/fs"
	"log"
	"net/http"
)

// StaticFilesHandler serves files from the /build directory inside embed.FS. It also serves the /build/index.html file
func StaticFilesHandler() http.Handler {
	staticBuildFS, err := fs.Sub(web.StaticFS, "build")
	if err != nil {
		log.Fatal(err)
	}

	return http.FileServer(http.FS(staticBuildFS))
}

// StaticHtmlHandler serves the /build/<name>.html file from embed.FS
func StaticHtmlHandler(name string) http.Handler {
	debugHtmlFile, err := fs.ReadFile(web.StaticFS, "build/"+name+".html")
	if err != nil {
		log.Fatal("StaticHtmlHandler() error: ", err)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(debugHtmlFile)
	})
}
