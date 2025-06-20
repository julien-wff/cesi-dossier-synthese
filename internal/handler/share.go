package handler

import (
	"encoding/json"
	"github.com/julien-wff/cesi-dossier-synthese/internal/parser"
	"github.com/julien-wff/cesi-dossier-synthese/internal/utils"
	"io"
	"net/http"
	"net/url"
)

func redirectWithError(w http.ResponseWriter, r *http.Request, err *utils.APIError) {
	jsonErr, _ := json.Marshal(err)
	http.Redirect(w, r, "/?error="+url.QueryEscape(string(jsonErr)), http.StatusSeeOther)
	_ = utils.LogParseTelemetry(r, nil, err)
}

func ParseSharePdfHandler(w http.ResponseWriter, r *http.Request) {
	// Handle panic
	defer func() {
		if rec := recover(); rec != nil {
			redirectWithError(w, r, utils.NewGradesExtractionError("parsing error"))
		}
	}()

	// Extract PDF file from request
	file, err := extractPdf(w, r)
	if err != nil {
		redirectWithError(w, r, err)
		return
	}

	defer func() {
		_ = file.Close()
	}()

	// Type assertion to convert multipart.File to io.ReadSeeker
	readSeeker, ok := file.(io.ReadSeeker)
	if !ok {
		redirectWithError(w, r, utils.NewTypeAssertionError())
		return
	}

	// Parse PDF file
	result, pt, err := parser.ParsePdf(&readSeeker)
	if err != nil {
		redirectWithError(w, r, err)
		return
	}

	// Write headers
	w.Header().Set("Server-Timing", pt.ServerTiming())

	// Convert the result into a JSON string
	jsonResult, jsonErr := json.Marshal(result)
	if jsonErr != nil {
		redirectWithError(w, r, utils.NewJsonEncoderError(jsonErr))
		return
	}

	// URL encode the JSON result
	encodedJsonResult := url.QueryEscape(string(jsonResult))

	// Redirect to the root URL with the JSON result as a query parameter
	http.Redirect(w, r, "/?result="+string(encodedJsonResult), http.StatusSeeOther)

	// Log to stats
	_ = utils.LogParseTelemetry(r, pt, nil)
}
