package main

import (
	"fmt"
	"github.com/julien-wff/cesi-dossier-synthese/internal/router"
	"net/http"
)

func main() {
	r := router.NewRouter()

	fmt.Println("Starting server on port 8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		panic(err)
	}
}
