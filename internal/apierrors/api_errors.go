package apierrors

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// APIError represents a structured API error
type APIError struct {
	Code    string          `json:"code"`
	Status  int             `json:"status"`
	Message APIErrorMessage `json:"message"`
	Err     error           `json:"-"`
}

// APIErrorMessage contains localized error messages
type APIErrorMessage struct {
	Fr string `json:"fr"`
	En string `json:"en"`
}

// Error implements the error interface for APIError
func (e *APIError) Error() string {
	return fmt.Sprintf("code: %s, message: %s, status: %d, details: %v", e.Code, e.Message.En, e.Status, e.Err)
}

// Write writes the error to the response writer as a JSON object
func (e *APIError) Write(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(e.Status)
	_ = json.NewEncoder(w).Encode(e)
}

// Unwrap returns the wrapped error
func (e *APIError) Unwrap() error {
	return e.Err
}

// Common error messages
var (
	JsonEncodingErrorMessage     = APIErrorMessage{Fr: "Erreur d'encodage JSON", En: "JSON encoding error"}
	FileTooBigErrorMessage       = APIErrorMessage{Fr: "Fichier trop volumineux", En: "File too big"}
	InvalidFileErrorMessage      = APIErrorMessage{Fr: "Fichier introuvable dans la requête", En: "File not found in request"}
	TypeAssertionErrorMessage    = APIErrorMessage{Fr: "Impossible de charger le fichier en mémoire", En: "Failed to load file in memory"}
	PdfReadingErrorMessage       = APIErrorMessage{Fr: "Erreur lors de la lecture du fichier PDF", En: "Error while reading the PDF file"}
	PdfLineParsingErrorMessage   = APIErrorMessage{Fr: "Erreur lors de l'extraction des lignes du PDF", En: "Error while extracting lines from the PDF"}
	GradesExtractionErrorMessage = APIErrorMessage{Fr: "Erreur lors de l'extraction des notes", En: "Error while extracting grades"}
)

// NewAPIError creates a new APIError instance
func NewAPIError(code string, status int, message APIErrorMessage, err error) *APIError {
	return &APIError{
		Code:    code,
		Status:  status,
		Message: message,
		Err:     err,
	}
}

func NewJsonEncoderError(err error) *APIError {
	return NewAPIError("json_encoding_error", http.StatusInternalServerError, JsonEncodingErrorMessage, err)
}

func NewFileTooBigError(err error) *APIError {
	return NewAPIError("file_too_big", http.StatusRequestEntityTooLarge, FileTooBigErrorMessage, err)
}

func NewInvalidFileError(err error) *APIError {
	return NewAPIError("invalid_file", http.StatusBadRequest, InvalidFileErrorMessage, err)
}

func NewTypeAssertionError() *APIError {
	return NewAPIError("type_assertion_error", http.StatusInternalServerError, TypeAssertionErrorMessage, nil)
}

func NewPdfReadingError(err error) *APIError {
	return NewAPIError("pdf_reading_error", http.StatusInternalServerError, PdfReadingErrorMessage, err)
}

func NewPdfLineParsingError(err error) *APIError {
	return NewAPIError("pdf_line_parsing_error", http.StatusInternalServerError, PdfLineParsingErrorMessage, err)
}

func NewGradesExtractionError(details string) *APIError {
	return NewAPIError("grades_extraction_error", http.StatusInternalServerError, GradesExtractionErrorMessage, fmt.Errorf(details))
}
