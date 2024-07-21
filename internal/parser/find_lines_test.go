package parser

import "testing"

func TestPageLine_atStart(t *testing.T) {
	var tests = []struct {
		name     string
		l        pageLine
		l1       pageLine
		expected bool
	}{
		{
			name:     "should detect at start",
			l:        pageLine{X1: 3, Y1: 7, X2: 11, Y2: 13},
			l1:       pageLine{X1: 3, Y1: 7, X2: 17, Y2: 23},
			expected: true,
		},
		{
			name:     "should detect at end",
			l:        pageLine{X1: 3, Y1: 7, X2: 11, Y2: 13},
			l1:       pageLine{X1: 17, Y1: 23, X2: 3, Y2: 7},
			expected: true,
		},
		{
			name:     "should not detect at end",
			l:        pageLine{X1: 3, Y1: 7, X2: 0, Y2: 0},
			l1:       pageLine{X1: 0, Y1: 0, X2: 0, Y2: 0},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.l.atStart(&tt.l1) != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, tt.l.atStart(&tt.l1))
			}
		})
	}
}

func TestPageLine_atEnd(t *testing.T) {
	var tests = []struct {
		name     string
		l        pageLine
		l1       pageLine
		expected bool
	}{
		{
			name:     "should detect at end",
			l:        pageLine{X1: 3, Y1: 7, X2: 11, Y2: 13},
			l1:       pageLine{X1: 17, Y1: 23, X2: 11, Y2: 13},
			expected: true,
		},
		{
			name:     "should detect at start",
			l:        pageLine{X1: 3, Y1: 7, X2: 11, Y2: 13},
			l1:       pageLine{X1: 11, Y1: 13, X2: 17, Y2: 23},
			expected: true,
		},
		{
			name:     "should not detect at start",
			l:        pageLine{X1: 0, Y1: 0, X2: 3, Y2: 7},
			l1:       pageLine{X1: 0, Y1: 0, X2: 0, Y2: 0},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.l.atEnd(&tt.l1) != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, tt.l.atEnd(&tt.l1))
			}
		})
	}
}

func TestPageLine_addNeighbour(t *testing.T) {
	var tests = []struct {
		name      string
		l         pageLine
		neighbour pageLine
		isStart   bool
		isEnd     bool
	}{
		{
			name:      "should add start neighbour",
			l:         pageLine{X1: 3, Y1: 7, X2: 11, Y2: 13},
			neighbour: pageLine{X1: 3, Y1: 7, X2: 17, Y2: 23},
			isStart:   true,
			isEnd:     false,
		},
		{
			name:      "should add end neighbour",
			l:         pageLine{X1: 3, Y1: 7, X2: 11, Y2: 13},
			neighbour: pageLine{X1: 11, Y1: 13, X2: 17, Y2: 23},
			isStart:   false,
			isEnd:     true,
		},
		{
			name:      "should add start and end neighbour",
			l:         pageLine{X1: 3, Y1: 7, X2: 11, Y2: 13},
			neighbour: pageLine{X1: 3, Y1: 7, X2: 11, Y2: 13},
			isStart:   true,
			isEnd:     true,
		},
		{
			name:      "should not add start or end neighbour",
			l:         pageLine{X1: 3, Y1: 7, X2: 11, Y2: 13},
			neighbour: pageLine{X1: 17, Y1: 23, X2: 23, Y2: 29},
			isStart:   false,
			isEnd:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.l.addNeighbour(&tt.neighbour)

			if tt.isStart && len(tt.l.StartNeighbours) == 0 {
				t.Errorf("expected start neighbour to be added")
			}

			if tt.isEnd && len(tt.l.EndNeighbours) == 0 {
				t.Errorf("expected end neighbour to be added")
			}
		})
	}
}
