package handler

import (
	"encoding/json"
	"github.com/julien-wff/cesi-dossier-synthese/internal/parser"
	"io"
	"mime/multipart"
	"net/http"
)

const (
	pdfMaxSize = 100 << 10 // 100KB
	pdfFormKey = "file"
)

// extractPdf extracts the PDF file from the request
func extractPdf(w http.ResponseWriter, r *http.Request) (multipart.File, error) {
	r.Body = http.MaxBytesReader(w, r.Body, pdfMaxSize)

	if err := r.ParseMultipartForm(pdfMaxSize); err != nil {
		return nil, err
	}

	file, _, err := r.FormFile(pdfFormKey)
	if err != nil {
		return nil, err
	}

	return file, nil
}

// ParsePdfDebugHandler handles the parsing of the request PDF file, and returns debug information about all
// the steps of the parsing process as JSON
func ParsePdfDebugHandler(w http.ResponseWriter, r *http.Request) {
	// Extract PDF file from request
	file, err := extractPdf(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer func() {
		_ = file.Close()
	}()

	// Type assertion to convert multipart.File to io.ReadSeeker
	readSeeker, ok := file.(io.ReadSeeker)
	if !ok {
		http.Error(w, "Failed to assert file type to io.ReadSeeker", http.StatusInternalServerError)
		return
	}

	// Parse PDF file
	result, pt, err := parser.ParsePdfDebug(&readSeeker)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write headers
	w.Header().Set("Server-Timing", pt.ServerTiming())
	w.Header().Set("Content-Type", "application/json")

	// Write response body as JSON
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
