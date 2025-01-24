package handler

import (
	"encoding/json"
	"fmt"
	"github.com/julien-wff/cesi-dossier-synthese/internal/parser"
	"github.com/julien-wff/cesi-dossier-synthese/internal/utils"
	"io"
	"mime/multipart"
	"net/http"
	"runtime/debug"
)

const (
	pdfMaxSize = 100 << 10 // 100KB
	pdfFormKey = "file"
)

// extractPdf extracts the PDF file from the request
func extractPdf(w http.ResponseWriter, r *http.Request) (multipart.File, *utils.APIError) {
	r.Body = http.MaxBytesReader(w, r.Body, pdfMaxSize)

	if err := r.ParseMultipartForm(pdfMaxSize); err != nil {
		return nil, utils.NewFileTooBigError(err)
	}

	file, _, err := r.FormFile(pdfFormKey)
	if err != nil {
		return nil, utils.NewFileNotFoundError(err)
	}

	return file, nil
}

// handleParsingPanic handles panics that may occur during the parsing process.
// It writes a GradesExtractionError to the response writer, and prints the stack trace.
func handleParsingPanic(w http.ResponseWriter) {
	if r := recover(); r != nil {
		utils.NewGradesExtractionError("parsing error").Write(w)
		fmt.Printf("Recovered from panic: %v\n", r)
		debug.PrintStack()
	}
}

// ParsePdfDebugHandler handles the parsing of the request PDF file, and returns debug information about all
// the steps of the parsing process as JSON
func ParsePdfDebugHandler(w http.ResponseWriter, r *http.Request) {
	// Handle panics
	defer handleParsingPanic(w)

	// Extract PDF file from request
	file, err := extractPdf(w, r)
	if err != nil {
		err.Write(w)
		return
	}

	defer func() {
		_ = file.Close()
	}()

	// Type assertion to convert multipart.File to io.ReadSeeker
	readSeeker, ok := file.(io.ReadSeeker)
	if !ok {
		utils.NewTypeAssertionError().Write(w)
		return
	}

	// Parse PDF file
	result, pt, err := parser.ParsePdfDebug(&readSeeker)
	if err != nil {
		err.Write(w)
		return
	}

	// Write headers
	w.Header().Set("Server-Timing", pt.ServerTiming())
	w.Header().Set("Content-Type", "application/json")

	// Write response body as JSON
	if err := json.NewEncoder(w).Encode(result); err != nil {
		utils.NewJsonEncoderError(err).Write(w)
	}
}

func ParsePdfHandler(w http.ResponseWriter, r *http.Request) {
	// Handle panic
	defer handleParsingPanic(w)

	// Extract PDF file from request
	file, err := extractPdf(w, r)
	if err != nil {
		err.Write(w)
		return
	}

	defer func() {
		_ = file.Close()
	}()

	// Type assertion to convert multipart.File to io.ReadSeeker
	readSeeker, ok := file.(io.ReadSeeker)
	if !ok {
		utils.NewTypeAssertionError().Write(w)
		return
	}

	// Parse PDF file
	result, pt, err := parser.ParsePdf(&readSeeker)
	if err != nil {
		err.Write(w)
		return
	}

	// Write headers
	w.Header().Set("Server-Timing", pt.ServerTiming())
	w.Header().Set("Content-Type", "application/json")

	// Write response body as JSON
	if err := json.NewEncoder(w).Encode(result); err != nil {
		utils.NewJsonEncoderError(err).Write(w)
	}
}
