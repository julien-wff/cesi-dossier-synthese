package parser

import (
	"strconv"
	"testing"
)

// Create a set of lines that looks like this:
//
//   -- l0  --   -- l1  -- l2  -- l3 --
//  |          |                       |
// l14        l4                      l5
//  |          |                       |
//                           -- l6  --
//  |          |           |           |
// l7         l8          l9          l10
//  |          |           |           |
//   -- l11 --   -- l12 --   -- l13 --

func setupPageLines(t *testing.T) PageLines {
	t.Helper()

	page := PageLines{}

	page.addLine(lineNode{X1: 0, Y1: 0, X2: 10, Y2: 0})    // l0
	page.addLine(lineNode{X1: 10, Y1: 0, X2: 20, Y2: 0})   // l1
	page.addLine(lineNode{X1: 10, Y1: 0, X2: 30, Y2: 0})   // l2
	page.addLine(lineNode{X1: 20, Y1: 0, X2: 30, Y2: 0})   // l3
	page.addLine(lineNode{X1: 10, Y1: 0, X2: 10, Y2: 10})  // l4
	page.addLine(lineNode{X1: 30, Y1: 0, X2: 30, Y2: 10})  // l5
	page.addLine(lineNode{X1: 20, Y1: 10, X2: 30, Y2: 10}) // l6
	page.addLine(lineNode{X1: 0, Y1: 10, X2: 0, Y2: 20})   // l7
	page.addLine(lineNode{X1: 10, Y1: 10, X2: 10, Y2: 20}) // l8
	page.addLine(lineNode{X1: 20, Y1: 10, X2: 20, Y2: 20}) // l9
	page.addLine(lineNode{X1: 30, Y1: 10, X2: 30, Y2: 20}) // l10
	page.addLine(lineNode{X1: 0, Y1: 20, X2: 10, Y2: 20})  // l11
	page.addLine(lineNode{X1: 10, Y1: 20, X2: 20, Y2: 20}) // l12
	page.addLine(lineNode{X1: 20, Y1: 20, X2: 30, Y2: 20}) // l13
	page.addLine(lineNode{X1: 0, Y1: 0, X2: 0, Y2: 10})    // l14

	return page
}

func findLineName(page *PageLines, line *pageLine) string {
	for i, l := range page.Lines {
		if l == line {
			return "l" + strconv.Itoa(i)
		}
	}
	return "<nil>"

}

func TestPageLine_continueStraight(t *testing.T) {
	page := setupPageLines(t)

	tests := []struct {
		name     string
		line     *pageLine
		inverted bool
		want     *pageLine
	}{
		{
			name:     "should go from l0 to l1",
			line:     page.Lines[0],
			inverted: false,
			want:     page.Lines[1],
		},
		{
			name:     "should go from l4 to l8",
			line:     page.Lines[4],
			inverted: false,
			want:     page.Lines[8],
		},
		{
			name:     "should go from l12 to l11",
			line:     page.Lines[12],
			inverted: true,
			want:     page.Lines[11],
		},
		{
			name:     "should go from l10 to l5",
			line:     page.Lines[10],
			inverted: true,
			want:     page.Lines[5],
		},
		{
			name:     "should not continue on l2",
			line:     page.Lines[2],
			inverted: false,
			want:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.line.continueStraight(tt.inverted); got != tt.want {
				t.Errorf("ContinueStraight() = %v, want %v", findLineName(&page, got), findLineName(&page, tt.want))
			}
		})
	}
}

func TestPageLine_continueLeft(t *testing.T) {
	page := setupPageLines(t)

	tests := []struct {
		name     string
		line     *pageLine
		inverted bool
		want     *pageLine
	}{
		{
			name:     "should go from l0 to l4",
			line:     page.Lines[0],
			inverted: false,
			want:     page.Lines[4],
		},
		{
			name:     "should go from l5 to l6",
			line:     page.Lines[5],
			inverted: false,
			want:     page.Lines[6],
		},
		{
			name:     "should go from l13 to l9",
			line:     page.Lines[13],
			inverted: true,
			want:     page.Lines[9],
		},
		{
			name:     "should go from l4 to l1",
			line:     page.Lines[4],
			inverted: true,
			want:     page.Lines[1],
		},
		{
			name:     "should go from l6 to l10",
			line:     page.Lines[6],
			inverted: false,
			want:     page.Lines[10],
		},
		{
			name:     "should not continue on l1",
			line:     page.Lines[1],
			inverted: false,
			want:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.line.continueLeft(tt.inverted); got != tt.want {
				t.Errorf("ContinueLeft() = %v, want %v", findLineName(&page, got), findLineName(&page, tt.want))
			}
		})
	}
}

func TestPageLine_getSmallestSquare(t *testing.T) {
	tests := []struct {
		name      string
		startLine int
		want      []int
	}{
		{
			name:      "should find square 6-10-13-9 from 6",
			startLine: 6,
			want:      []int{6, 10, 13, 9},
		},
		{
			name:      "should find square 0-4-8-11-7-14 from 0",
			startLine: 0,
			want:      []int{0, 4, 8, 11, 7, 14},
		},
		{
			name:      "should find square 2-5-10-13-12-8-4 from 2",
			startLine: 2,
			want:      []int{2, 5, 10, 13, 12, 8, 4},
		},
		{
			name:      "should find square 1-3-5-10-13-12-8-4 from 1",
			startLine: 1,
			want:      []int{1, 3, 5, 10, 13, 12, 8, 4},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			page := setupPageLines(t)
			startLine := page.Lines[tt.startLine]
			square := startLine.getSmallestSquare()

			expected := ""
			for _, l := range tt.want {
				expected += findLineName(&page, page.Lines[l]) + ", "
			}
			expected = expected[:len(expected)-2]

			result := ""
			for _, l := range square {
				result += findLineName(&page, l) + ", "
			}
			result = result[:len(result)-2]

			if expected != result {
				t.Errorf("GetSmallestSquare() = %v, want %v", result, expected)
			}
		})
	}
}
