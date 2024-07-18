package parser

import (
	"github.com/julien-wff/cesi-dossier-synthese/internal/utils"
	"io"
)

// PdfParseDebugResponse is the response of the ParsePdfDebug function,
// containing the full content of the PDF file through all the steps of the parsing.
// It will be rendered as JSON in the API response.
type PdfParseDebugResponse struct {
	Performance []utils.TimingElement `json:"performance"`
	Pages       *[]pdfPageContent     `json:"pages"`
	Lines       []*PageLines          `json:"lines"`
	Squares     []*PageSquares        `json:"squares"`
}

// ParsePdfDebug parses a PDF file returns the full retrieved content through all the steps of the parsing
func ParsePdfDebug(f *io.ReadSeeker) (PdfParseDebugResponse, *utils.ProcessTiming, error) {
	// Initialize performance counter
	pt := utils.NewProcessTiming()

	// Initialize response
	response := PdfParseDebugResponse{}

	// Extract the raw content of the PDF
	pages, err := extractRawPdfContent(f, pt, true)
	if err != nil {
		return response, pt, err
	}
	response.Pages = pages

	// Find lines
	var lines []*PageLines
	for _, page := range *pages {
		pageLines, err := findPageLines(&page, true)
		if err != nil {
			return response, pt, err
		}
		lines = append(lines, &pageLines)
	}
	response.Lines = lines
	pt.AddElement("find-lines", "Find pages lines")

	// Find squares
	var squares []*PageSquares
	for i, line := range lines {
		pageSquares := findPageSquares(line, &(*pages)[i])
		squares = append(squares, pageSquares)
	}
	response.Squares = squares
	pt.AddElement("find-squares", "Find pages squares")

	// Add performance counter
	response.Performance = pt.Elements

	return response, pt, nil
}
