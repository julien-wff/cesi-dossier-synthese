package parser

import (
	"sort"
	"strings"
)

// continueStraight returns the neighbour pageLine that continues straight from the line.
// Inverted is used to determine the direction to look for the neighbour:
//
//	false: continues to the right / bottom of the line.
//	true: continues to the left / top of the line.
func (l *pageLine) continueStraight(inverted bool) *pageLine {
	if l == nil {
		return nil
	}

	neighbours := make([]*pageLine, 0)

	if !inverted {
		neighbours = l.EndNeighbours
	} else {
		neighbours = l.StartNeighbours
	}

	// Filter out neighbours that are not in the same direction or does not continue straight (depends on inverted)
	filteredNeighbours := make([]*pageLine, 0)
	for _, neighbour := range neighbours {
		lineEnd, neighbourEnd := 0.0, 0.0

		if l.Direction == horizontal {
			lineEnd = l.X2
			neighbourEnd = neighbour.X2
		} else {
			lineEnd = l.Y2
			neighbourEnd = neighbour.Y2
		}

		// If neighbour is not in the same direction, skip
		if neighbour.Direction != l.Direction {
			continue
		}

		// If neighbour does not continue straight, skip
		if (!inverted && neighbourEnd <= lineEnd) || (inverted && neighbourEnd >= lineEnd) {
			continue
		}

		filteredNeighbours = append(filteredNeighbours, neighbour)
	}

	// Sort neighbours by shortest length first
	sort.Slice(filteredNeighbours, func(i, j int) bool {
		return filteredNeighbours[i].Length < filteredNeighbours[j].Length
	})

	if len(filteredNeighbours) == 0 {
		return nil
	} else {
		return filteredNeighbours[0]
	}
}

// continueLeft returns the neighbour PageLine that starts where the line ends, at an angle of 90°.
// Inverted is used to determine the direction to look for the neighbour:
//
//	false: continues to the right / bottom of the line.
//	true: continues to the left / top of the line.
func (l *pageLine) continueLeft(inverted bool) *pageLine {
	if l == nil {
		return nil
	}

	neighbours := make([]*pageLine, 0)
	if !inverted {
		neighbours = l.EndNeighbours
	} else {
		neighbours = l.StartNeighbours
	}

	var lineDirInverted = (l.Direction == "vertical" && !inverted) || (l.Direction == "horizontal" && inverted)

	// Filter out neighbours that are in the same direction or does not continue left (depends on inverted)
	filteredNeighbours := make([]*pageLine, 0)
	for i, neighbour := range neighbours {
		lineEnd, neighbourEnd := 0.0, 0.0
		if l.Direction == horizontal {
			if !lineDirInverted {
				lineEnd = l.Y2
				neighbourEnd = neighbour.Y2
			} else {
				lineEnd = l.Y1
				neighbourEnd = neighbour.Y1
			}
		} else {
			if !lineDirInverted {
				lineEnd = l.X2
				neighbourEnd = neighbour.X2
			} else {
				lineEnd = l.X1
				neighbourEnd = neighbour.X1
			}
		}

		// If neighbour is in the same direction, skip
		if neighbour.Direction == l.Direction {
			continue
		}

		// If neighbour does not continue left, skip
		if (!lineDirInverted && neighbourEnd <= lineEnd) || (lineDirInverted && neighbourEnd >= lineEnd) {
			continue
		}

		filteredNeighbours = append(filteredNeighbours, neighbours[i])
	}

	// Sort neighbours by shortest length first
	sort.Slice(filteredNeighbours, func(i, j int) bool {
		return filteredNeighbours[i].Length < filteredNeighbours[j].Length
	})

	if len(filteredNeighbours) == 0 {
		return nil
	} else {
		return filteredNeighbours[0]
	}
}

// getSmallestSquare returns the smallest square, starting from the pageLine l.
// It returns two values:
//
//   - The square, as a 2D array of pageLines, where the first array is the top line, going clockwise.
//   - The path, as a 1D array of pageLines, where the first element is the starting line, and the last element is the
//     line just before.
func (l *pageLine) getSmallestSquare() ([][]*pageLine, []*pageLine) {
	path := make([]*pageLine, 0)
	path = append(path, l)
	result := make([][]*pageLine, 4)
	result[0] = append(result[0], l)

	// lines direction: 1 = top, 2 = right, 3 = bottom, 4 = left
	direction := 1

	for {
		inverted := direction > 2
		lastLine := path[len(path)-1]

		// Try to continue left
		next := lastLine.continueLeft(inverted)

		// If no neighbour was found, try to continue straight
		if next == nil {
			next = lastLine.continueStraight(inverted)
		} else {
			direction++
		}

		// If the circle is closed, or no neighbour was found, break
		if direction > 4 || next == nil {
			break
		}

		path = append(path, next)
		result[direction-1] = append(result[direction-1], next)
	}

	return result, path
}

// pageSquare represents a square on a page, defined by its top-left and bottom-right coordinates.
// It also contains the lines that form the square.
type pageSquare struct {
	X1      float64 `json:"x1"`
	Y1      float64 `json:"y1"`
	X2      float64 `json:"x2"`
	Y2      float64 `json:"y2"`
	Content string  `json:"content"`
	lines   [][]*pageLine
}

// newPageSquare creates a new pageSquare from a 2D array of pageLines (returned from getSmallestSquare).
// Returns nil if the square is not valid.
func newPageSquare(lineSquare [][]*pageLine) *pageSquare {
	if len(lineSquare) < 4 || len(lineSquare[0]) == 0 || len(lineSquare[2]) == 0 {
		return nil
	}

	return &pageSquare{
		X1:    lineSquare[0][0].X1,
		Y1:    lineSquare[0][0].Y1,
		X2:    lineSquare[2][0].X2,
		Y2:    lineSquare[2][0].Y2,
		lines: lineSquare,
	}
}

// findLeft returns the square that is on the left of the current square.
// Returns nil if no square is found.
func (ps *pageSquare) findLeft() *pageSquare {
	// Find the left line
	if len(ps.lines[0]) == 0 {
		return nil
	}

	// Find the starting line of the left square (its first top line)
	leftLine := ps.lines[0][len(ps.lines[0])-1].continueStraight(false)
	if leftLine == nil {
		return nil
	}

	// Find the smallest square starting from the left line
	square, _ := leftLine.getSmallestSquare()
	return newPageSquare(square)
}

// findBottom returns the square that is at the bottom of the current square.
// Returns nil if no square is found.
func (ps *pageSquare) findBottom() *pageSquare {
	// Find the bottom line
	if len(ps.lines[3]) == 0 {
		return nil
	}

	// Find the starting line of the bottom square (its first top line)
	bottomLine := ps.lines[3][0].continueStraight(false).continueLeft(true)

	// Find the smallest square starting from the bottom line
	square, _ := bottomLine.getSmallestSquare()
	return newPageSquare(square)
}

func (ps *pageSquare) addTextContent(pageContent *pdfPageContent) {
	for _, text := range pageContent.Text {
		if text.Position.X >= ps.X1 && text.Position.X <= ps.X2 && text.Position.Y >= ps.Y1 && text.Position.Y <= ps.Y2 {
			text := strings.TrimSpace(text.Content)
			if text != "" {
				ps.Content += text + " "
			}
		}
	}

	ps.Content = strings.TrimSpace(ps.Content)
}

// PageSquares represents all the squares found in a page.
type PageSquares struct {
	Squares [][]*pageSquare `json:"squares"`
	Page    int             `json:"page"`
}

// findPageSquares returns all the squares found in the page.
func findPageSquares(lines *PageLines, pageContent *pdfPageContent) *PageSquares {
	result := PageSquares{
		Page:    lines.Page,
		Squares: make([][]*pageSquare, 0),
	}

	// Find the starting square, based on the first page line
	lineSquare, _ := lines.Lines[0].getSmallestSquare()
	bottomSquare := newPageSquare(lineSquare)
	result.Squares = append(result.Squares, []*pageSquare{bottomSquare})

	lineIndex := 0

	for {
		// Find left squares
		leftSquare := bottomSquare
		for {
			leftSquare = leftSquare.findLeft()
			if leftSquare == nil {
				break
			}

			result.Squares[lineIndex] = append(result.Squares[lineIndex], leftSquare)
		}

		// Find bottom squares
		bottomSquare = bottomSquare.findBottom()
		if bottomSquare == nil {
			break
		}

		lineIndex++
		result.Squares = append(result.Squares, []*pageSquare{bottomSquare})
	}

	// Add text content to each square
	for _, line := range result.Squares {
		for _, square := range line {
			square.addTextContent(pageContent)
		}
	}

	return &result
}
