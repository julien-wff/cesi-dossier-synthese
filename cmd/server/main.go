package main

import (
	"fmt"
	"github.com/julien-wff/cesi-dossier-synthese/internal/router"
	"github.com/julien-wff/cesi-dossier-synthese/internal/utils"
	"net/http"
)

func main() {
	cfg := utils.GetAppConfig()

	r := router.NewRouter()

	if cfg.Production {
		fmt.Println("Starting production server on port", cfg.Port)
	} else {
		fmt.Println("Starting development server on port", cfg.Port)
	}

	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		panic(err)
	}
}
