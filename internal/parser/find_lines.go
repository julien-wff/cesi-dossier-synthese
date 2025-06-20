package parser

import (
	"fmt"
	"github.com/julien-wff/cesi-dossier-synthese/internal/utils"
	"math"
	"sort"
)

// lineDirection represents the direction of a line (horizontal or vertical)
type lineDirection string

const (
	horizontal lineDirection = "horizontal"
	vertical   lineDirection = "vertical"
)

// pageLine represents a line in a PDF page, with its start and end coordinates, direction, length and neighbours
type pageLine struct {
	Id                 int           `json:"id"`
	X1                 float64       `json:"x1"`
	Y1                 float64       `json:"y1"`
	X2                 float64       `json:"x2"`
	Y2                 float64       `json:"y2"`
	Direction          lineDirection `json:"direction"`
	Length             float64       `json:"length"`
	Neighbours         []*pageLine   `json:"-"`
	StartNeighbours    []*pageLine   `json:"-"`
	EndNeighbours      []*pageLine   `json:"-"`
	StartNeighboursIds []int         `json:"start_neighbours_ids"`
	EndNeighboursIds   []int         `json:"end_neighbours_ids"`
}

// atStart checks if the beginning or the end of l1 is the same as l's start
func (l *pageLine) atStart(l1 *pageLine) bool {
	return (l.X1 == l1.X1 && l.Y1 == l1.Y1) || (l.X1 == l1.X2 && l.Y1 == l1.Y2)
}

// atEnd checks if the beginning or the end of l1 is the same as l's end
func (l *pageLine) atEnd(l1 *pageLine) bool {
	return (l.X2 == l1.X1 && l.Y2 == l1.Y1) || (l.X2 == l1.X2 && l.Y2 == l1.Y2)
}

// gravity returns the gravity center of l, on the X and Y axis
func (l *pageLine) gravity() (float64, float64) {
	return (l.X1 + l.X2) / 2, (l.Y1 + l.Y2) / 2
}

// isOutside checks if l is outside the zone
func (l *pageLine) isOutside(zone size, margin float64) bool {
	if l.X1 < margin || l.Y1 < margin || l.X2 < margin || l.Y2 < margin {
		return true
	}

	if l.X1 > zone.Width-margin || l.Y1 > zone.Height-margin || l.X2 > zone.Width-margin || l.Y2 > zone.Height-margin {
		return true
	}

	return false
}

// addNeighbour adds a neighbour to l and vice versa
func (l *pageLine) addNeighbour(neighbour *pageLine) {
	l.Neighbours = append(l.Neighbours, neighbour)
	neighbour.Neighbours = append(neighbour.Neighbours, l)

	if l.atStart(neighbour) {
		l.StartNeighbours = append(l.StartNeighbours, neighbour)
	}

	if l.atEnd(neighbour) {
		l.EndNeighbours = append(l.EndNeighbours, neighbour)
	}

	if neighbour.atStart(l) {
		neighbour.StartNeighbours = append(neighbour.StartNeighbours, l)
	}

	if neighbour.atEnd(l) {
		neighbour.EndNeighbours = append(neighbour.EndNeighbours, l)
	}
}

// removeNeighbour removes a neighbour from l and vice versa
func (l *pageLine) removeNeighbour(neighbour *pageLine) {
	for i, n := range l.Neighbours {
		if n == neighbour {
			l.Neighbours = append(l.Neighbours[:i], l.Neighbours[i+1:]...)
			break
		}
	}

	for i, n := range l.StartNeighbours {
		if n == neighbour {
			l.StartNeighbours = append(l.StartNeighbours[:i], l.StartNeighbours[i+1:]...)
			break
		}
	}

	for i, n := range l.EndNeighbours {
		if n == neighbour {
			l.EndNeighbours = append(l.EndNeighbours[:i], l.EndNeighbours[i+1:]...)
			break
		}
	}
}

// hasSameCoordinates checks if l and l1 have the same coordinates (works if the coordinates are inverted)
func (l *pageLine) hasSameCoordinates(l1 *pageLine) bool {
	return (l.X1 == l1.X1 && l.Y1 == l1.Y1 && l.X2 == l1.X2 && l.Y2 == l1.Y2) || (l.X1 == l1.X2 && l.Y1 == l1.Y2 && l.X2 == l1.X1 && l.Y2 == l1.Y1)
}

// String formats the coordinates of l
func (l *pageLine) String() string {
	return fmt.Sprintf("(%f, %f) ; (%f, %f)", l.X1, l.Y1, l.X2, l.Y2)
}

// PageLines represents a list of all the lines in a PDF page
type PageLines struct {
	Page  int         `json:"page"`
	Lines []*pageLine `json:"lines"`
}

// addLine adds a line to the list.
// Also calculates the information and checks if the line is already in the list.
func (pls *PageLines) addLine(line lineNode) {
	var direction lineDirection
	if line.X1 == line.X2 {
		direction = vertical
	} else {
		direction = horizontal
	}

	x1 := math.Round(line.X1)
	y1 := math.Round(line.Y1)
	x2 := math.Round(line.X2)
	y2 := math.Round(line.Y2)

	x1, y1, x2, y2 = reorderCoordinates(x1, y1, x2, y2)

	var length = math.Sqrt(math.Pow(x2-x1, 2) + math.Pow(y2-y1, 2))
	if length == 0 {
		return
	}

	pl := pageLine{
		Id:              len(pls.Lines),
		X1:              x1,
		Y1:              y1,
		X2:              x2,
		Y2:              y2,
		Direction:       direction,
		Length:          length,
		Neighbours:      []*pageLine{},
		StartNeighbours: []*pageLine{},
		EndNeighbours:   []*pageLine{},
	}

	// Check if the line is already in the list
	for _, l := range pls.Lines {
		if l.hasSameCoordinates(&pl) {
			return
		}
	}

	// Add as neighbour if it has one point in common
	for _, l := range pls.Lines {
		if (l.X1 == x1 && l.Y1 == y1) || (l.X2 == x2 && l.Y2 == y2) || (l.X1 == x2 && l.Y1 == y2) || (l.X2 == x1 && l.Y2 == y1) {
			l.addNeighbour(&pl)
		}
	}

	pls.Lines = append(pls.Lines, &pl)
}

// removeLine removes a line from the list and its neighbours
func (pls *PageLines) removeLine(line *pageLine) {
	// Remove line from its neighbours
	for _, n := range line.Neighbours {
		n.removeNeighbour(line)
	}

	// Remove line from the list
	for i, l := range pls.Lines {
		if l == line {
			pls.Lines = append(pls.Lines[:i], pls.Lines[i+1:]...)
			break
		}
	}
}

// mergeLines merges two lines that have one point in common.
// It removes l2 and adds its neighbours to l1.
// Once merged, if the new l1 has the same coordinates as another line in pls, it is removed.
func (pls *PageLines) mergeLines(l1 *pageLine, l2 *pageLine) {
	l1.X2, l1.Y2 = l2.X2, l2.Y2

	// Save l2's neighbours
	neighbours := l2.Neighbours

	// Remove l2
	pls.removeLine(l2)

	// Add l2's neighbours to l1
	for _, n := range neighbours {
		if n != l1 {
			l1.addNeighbour(n)
		}
	}

	// Check if a line in pls has the same coordinates as l1
	// If so, remove l1
	for _, l := range pls.Lines {
		if l != l1 && l.hasSameCoordinates(l1) {
			pls.removeLine(l1)
			break
		}
	}
}

// optimize merges lines that have one point in common, and reorders by coordinates.
func (pls *PageLines) optimize(pageContent *pdfPageContent) {
	// Merge lines that have one point in common
	for _, l := range pls.Lines {
		// Only merge lines that have one neighbour. Else, it's an intersection
		if len(l.EndNeighbours) != 1 {
			continue
		}

		// Get the end neighbour. It must have the same direction
		endNeighbour := l.EndNeighbours[0]
		if endNeighbour.Direction != l.Direction {
			continue
		}

		// Merge the two lines
		pls.mergeLines(l, endNeighbour)
	}

	// Remove lines that are too close to the page borders
	for _, l := range pls.Lines {
		if l.isOutside(pageContent.Size, 10) {
			pls.removeLine(l)
		}
	}

	// Reorder by coordinates to start from the top left corner, then like reading a book
	sort.Slice(pls.Lines, func(i, j int) bool {
		// Calculate x and y center of the lines
		xi, yi := pls.Lines[i].gravity()
		xj, yj := pls.Lines[j].gravity()

		// Compare the two lines
		if yi == yj {
			return xi < xj
		}

		return yi < yj
	})
}

// verifyNeighbours checks if the neighbours count is correct.
func (pls *PageLines) verifyNeighbours() error {
	for _, l := range pls.Lines {
		neighboursCount := len(l.Neighbours)
		startNeighboursCount := len(l.StartNeighbours)
		endNeighboursCount := len(l.EndNeighbours)

		// There should not be neighbours at both the start and the end
		if neighboursCount != (startNeighboursCount + endNeighboursCount) {
			return fmt.Errorf("neighbours count is not correct: %d %d %d (%.f, %.f), (%.f, %.f)",
				neighboursCount,
				startNeighboursCount,
				endNeighboursCount,
				l.X1,
				l.Y1,
				l.X2,
				l.Y2,
			)
		}
	}

	return nil
}

// addDebugInfos adds debug information to the lines, so they can be displayed in the frontend
func (pls *PageLines) addDebugInfos() {
	for _, l := range pls.Lines {
		l.StartNeighboursIds = make([]int, len(l.StartNeighbours))
		for i, n := range l.StartNeighbours {
			l.StartNeighboursIds[i] = n.Id
		}

		l.EndNeighboursIds = make([]int, len(l.EndNeighbours))
		for i, n := range l.EndNeighbours {
			l.EndNeighboursIds[i] = n.Id
		}
	}
}

// reorderCoordinates reorders the coordinates of a line, so that the first point is the top left corner
func reorderCoordinates(x1, y1, x2, y2 float64) (float64, float64, float64, float64) {
	if x1 > x2 || y1 > y2 {
		return x2, y2, x1, y1
	}

	return x1, y1, x2, y2
}

// findPageLines finds all the lines in a PDF page, based on the rectangles and lines found in the content.
// If debug is true, it will calculate and add debug information for the frontend.
func findPageLines(content *pdfPageContent, debug bool) (PageLines, *utils.APIError) {
	pageLines := PageLines{
		Page: content.Page,
	}

	// Rectangles
	for _, rect := range content.Rectangles {
		pageLines.addLine(lineNode{
			X1: rect.Position.X,
			Y1: rect.Position.Y,
			X2: rect.Position.X + rect.Size.Width,
			Y2: rect.Position.Y,
		})
		pageLines.addLine(lineNode{
			X1: rect.Position.X,
			Y1: rect.Position.Y,
			X2: rect.Position.X,
			Y2: rect.Position.Y + rect.Size.Height,
		})
		pageLines.addLine(lineNode{
			X1: rect.Position.X + rect.Size.Width,
			Y1: rect.Position.Y,
			X2: rect.Position.X + rect.Size.Width,
			Y2: rect.Position.Y + rect.Size.Height,
		})
		pageLines.addLine(lineNode{
			X1: rect.Position.X,
			Y1: rect.Position.Y + rect.Size.Height,
			X2: rect.Position.X + rect.Size.Width,
			Y2: rect.Position.Y + rect.Size.Height,
		})
	}

	// Lines
	for _, line := range content.Lines {
		pageLines.addLine(line)
	}

	// Optimizations
	pageLines.optimize(content)

	// Verify neighbours count
	err := pageLines.verifyNeighbours()
	if err != nil {
		return pageLines, utils.NewPdfLineParsingError(err)
	}

	// Debug information
	if debug {
		pageLines.addDebugInfos()
	}

	return pageLines, nil
}
