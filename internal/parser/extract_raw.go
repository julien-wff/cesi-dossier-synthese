package parser

import (
	"bytes"
	"github.com/julien-wff/cesi-dossier-synthese/internal/utils"
	"github.com/ledongthuc/pdf"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"io"
	"strconv"
	"strings"
)

// position represents a 2D position
type position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// size represents a 2D size
type size struct {
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

// textNode represents a text element in a PDF
type textNode struct {
	Content  string   `json:"content"`
	FontSize float64  `json:"font_size"`
	Position position `json:"position"`
}

// rectangleNode represents a rectangle element in a PDF
type rectangleNode struct {
	Position position `json:"position"`
	Size     size     `json:"size"`
}

// lineNode represents a line element in a PDF
type lineNode struct {
	X1 float64 `json:"x1"`
	X2 float64 `json:"x2"`
	Y1 float64 `json:"y1"`
	Y2 float64 `json:"y2"`
}

// pdfPageContent represents all the content (text, lines, rectangles) of a PDF page
type pdfPageContent struct {
	Page       int             `json:"page"` // Start from 0
	Size       size            `json:"size"`
	Text       []textNode      `json:"text"`
	Rectangles []rectangleNode `json:"rectangles"`
	Lines      []lineNode      `json:"lines"`
}

// extractRawPdfContent extracts the raw content of a PDF file page by page
// This includes text, rectangles and lines
func extractRawPdfContent(f *io.ReadSeeker, pt *utils.ProcessTiming) (*[]pdfPageContent, error) {
	// Read the content of the original file into a buffer
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, *f); err != nil {
		return nil, err
	}

	// Create a writer to the buffer to apply modifications
	var modifiedBuf bytes.Buffer
	writer := io.Writer(&modifiedBuf)

	// Create configuration for the PDF operation
	conf := model.NewDefaultConfiguration()
	conf.Cmd = model.OPTIMIZE

	// Check if the file is a valid PDF
	ctx, err := api.ReadValidateAndOptimize(*f, conf)
	if err != nil {
		return nil, err
	}

	pt.AddElement("read-validate-optimize", "Read, validate and optimize")

	// Apply the optimization to the buffer
	if err := api.Optimize(bytes.NewReader(buf.Bytes()), writer, conf); err != nil {
		return nil, err
	}

	pt.AddElement("optimize", "Optimize")

	// Reinitialize the reader on the modified content
	modifiedReader := bytes.NewReader(modifiedBuf.Bytes())
	sectionReader := &utils.ReadSeekerAt{ReadSeeker: modifiedReader}

	// Obtain the size of the modified content
	size := int64(modifiedBuf.Len())

	reader, err := pdf.NewReader(sectionReader, size)
	if err != nil {
		return nil, err
	}

	pagesCount := reader.NumPage()

	conf.Cmd = model.EXTRACTCONTENT

	pagesStringIndexes := make([]string, pagesCount)
	for i := 0; i < pagesCount; i++ {
		pagesStringIndexes[i] = strconv.Itoa(i + 1)
	}

	_, err = api.PagesForPageSelection(pagesCount, pagesStringIndexes, true, false)
	if err != nil {
		return nil, err
	}

	pages := make([]pdfPageContent, pagesCount)

	// Extract the content of each page
	for i := 0; i < pagesCount; i++ {
		p := reader.Page(i + 1)
		pages[i].Page = i
		pages[i].Size = getPageSize(&p)
		pages[i].Text = getTextContent(&p, pages[i].Size.Height)
		pages[i].Rectangles = getRectangleContent(&p, pages[i].Size.Height)

		pageReader, err := pdfcpu.ExtractPageContent(ctx, i+1)
		if err != nil {
			return nil, err
		}

		lines, err := getLineContent(&pageReader, pages[i].Size.Height)
		if err != nil {
			return nil, err
		}
		pages[i].Lines = lines

		pt.AddElement("extract-page-"+strconv.Itoa(i+1)+"-content", "Extract page "+strconv.Itoa(i+1)+" content")
	}

	return &pages, nil
}

// getPageSize returns the size of a PDF page from the MediaBox
func getPageSize(p *pdf.Page) size {
	mediaBox := p.V.Key("MediaBox")
	llx := mediaBox.Index(0).Float64()
	lly := mediaBox.Index(1).Float64()
	urx := mediaBox.Index(2).Float64()
	ury := mediaBox.Index(3).Float64()
	width := urx - llx
	height := ury - lly

	return size{
		Width:  width,
		Height: height,
	}
}

// getTextContent returns all the text elements of a PDF page
func getTextContent(p *pdf.Page, pageHeight float64) []textNode {
	var content = make(map[position]textNode)

	rows := p.Content().Text
	for _, row := range rows {
		pos := position{X: row.X, Y: row.Y}
		var text = ""
		if content[pos] == (textNode{}) {
			text = row.S
		} else {
			text = content[pos].Content + row.S
		}

		content[pos] = textNode{
			Content:  text,
			FontSize: row.FontSize,
			Position: position{},
		}
	}

	var textNodes []textNode
	for pos, text := range content {
		text.Position = position{
			X: pos.X,
			Y: pageHeight - pos.Y,
		}
		textNodes = append(textNodes, text)
	}

	return textNodes
}

// getRectangleContent returns all the rectangles of a PDF page
func getRectangleContent(p *pdf.Page, pageHeight float64) []rectangleNode {
	rectangles := p.Content().Rect
	rectangleNodes := make([]rectangleNode, len(rectangles))

	for _, rect := range rectangles {
		rectangleNodes = append(rectangleNodes, rectangleNode{
			Position: position{
				X: rect.Min.X,
				Y: pageHeight - rect.Max.Y,
			},
			Size: size{
				Width:  rect.Max.X - rect.Min.X,
				Height: rect.Max.Y - rect.Min.Y,
			},
		})
	}

	return rectangleNodes
}

// getLineContent returns all the lines of a PDF page
func getLineContent(reader *io.Reader, pageHeight float64) ([]lineNode, error) {
	linesNode := make([]lineNode, 0)

	buf := new(strings.Builder)
	_, err := io.Copy(buf, *reader)
	if err != nil {
		return nil, err
	}

	result := buf.String()
	commands := strings.Split(result, "\n")

	lastX, lastY := 0.0, 0.0

	for _, command := range commands {
		if strings.HasSuffix(command, " m") {
			coordinates := strings.Split(command, " ")
			x, _ := strconv.ParseFloat(coordinates[0], 64)
			y, _ := strconv.ParseFloat(coordinates[1], 64)
			lastX, lastY = x, y
		} else if strings.HasSuffix(command, " l") {
			coordinates := strings.Split(command, " ")
			x, _ := strconv.ParseFloat(coordinates[0], 64)
			y, _ := strconv.ParseFloat(coordinates[1], 64)
			linesNode = append(linesNode, lineNode{
				X1: lastX,
				Y1: pageHeight - lastY,
				X2: x,
				Y2: pageHeight - y,
			})
		}
	}

	return linesNode, nil
}