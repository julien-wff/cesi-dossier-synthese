package parser

import (
	"sort"
)

// continueStraight returns the neighbour pageLine that continues straight from the line.
// Inverted is used to determine the direction to look for the neighbour:
//
//	false: continues to the right / bottom of the line.
//	true: continues to the left / top of the line.
func (l *pageLine) continueStraight(inverted bool) *pageLine {
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

// getSmallestSquare returns the smallest square, starting from the line.
func (l *pageLine) getSmallestSquare() []*pageLine {
	path := make([]*pageLine, 0)
	path = append(path, l)

	direction := 1
	straightMoves := 0
	skipTurn := false

	for {
		inverted := direction > 3
		lastLine := path[len(path)-1]

		var next *pageLine
		if !skipTurn {
			next = lastLine.continueLeft(inverted)
		}
		skipTurn = false

		if next == nil {
			next = lastLine.continueStraight(inverted)
			straightMoves++
		} else {
			direction = direction%4 + 1
			straightMoves = 0
		}

		// Rollback straight moves
		if next == nil && straightMoves > 1 {
			path = path[:len(path)-straightMoves]
			straightMoves = 0
			skipTurn = true
			break
		}

		path = append(path, next)

		if direction == 4 || next == nil {
			break
		}
	}

	return path
}

// findPageSquares returns all the squares found in the page.
func findPageSquares(lines *PageLines) []*pageLine {
	line := lines.Lines[0]
	result := make([]*pageLine, 0)
	result = append(result, line.getSmallestSquare()...)

	return result
}
