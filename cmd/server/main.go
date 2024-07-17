package main

import (
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/julien-wff/cesi-dossier-synthese/internal/router"
	"github.com/julien-wff/cesi-dossier-synthese/internal/utils"
	"net/http"
	"os"
)

func main() {
	cfg := utils.GetAppConfig()

	r := router.NewRouter()

	if cfg.Production {
		fmt.Println("Starting production server on port", cfg.Port)
	} else {
		fmt.Println("Starting development server on port", cfg.Port)
		r = handlers.CORS()(r)
		r = handlers.LoggingHandler(os.Stdout, r)
	}

	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		panic(err)
	}
}
