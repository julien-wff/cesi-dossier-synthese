package parser

import (
	"github.com/julien-wff/cesi-dossier-synthese/internal/apierrors"
	"strconv"
	"strings"
)

type GradeLetter string

const (
	A GradeLetter = "A"
	B GradeLetter = "B"
	C GradeLetter = "C"
	D GradeLetter = "D"
)

// Grade struct definition
type Grade struct {
	Name          string       `json:"name"`
	Letter        *GradeLetter `json:"letter"`
	Coefficient   float64      `json:"coefficient"`
	PreviousGrade *GradeLetter `json:"previousGrade"`
}

// Category struct definition
type Category struct {
	Name   string  `json:"name"`
	Grades []Grade `json:"grades"`
}

// Section struct definition
type Section struct {
	Name       string     `json:"name"`
	Categories []Category `json:"categories"`
}

// ExtractGrades extracts the grades from the page squares from all the pages.
// Here are the different steps:
//
//   - Sections are the main categories of the grades. They take a whole square line.
//   - Categories are the subcategories of the grades. They are 3 squares wide.
//   - Grades are the actual grades. They are 6 squares wide.
func extractGrades(pageSquares []*PageSquares) ([]Section, *apierrors.APIError) {
	squares := make([][]*pageSquare, 0)

	for _, squareLine := range pageSquares {
		// Check if headers are correct
		if len(squareLine.Squares) < 2 || len(squareLine.Squares[0]) != 6 {
			continue
		}

		// Remove the first line
		currentSquareLines := squareLine.Squares[1:]

		// Append the current square lines to the squares array
		squares = append(squares, currentSquareLines...)
	}

	// Check if there are enough squares
	if len(squares) < 2 {
		return []Section{}, apierrors.NewGradesExtractionError("not enough squares")
	}

	// Get sections indexes (full width squares)
	sectionsIndexes := make([]int, 0)
	for i, square := range squares {
		if len(square) == 1 {
			sectionsIndexes = append(sectionsIndexes, i)
		}
	}

	// Check if there are enough sections
	if len(sectionsIndexes) == 0 {
		return []Section{}, apierrors.NewGradesExtractionError("no sections found")
	}

	// Parse sections
	sections := make([]Section, 0)
	for i, sectionIndex := range sectionsIndexes {
		var nextSectionIndex int
		if i == len(sectionsIndexes)-1 {
			nextSectionIndex = len(squares)
		} else {
			nextSectionIndex = sectionsIndexes[i+1]
		}
		sectionSquares := squares[sectionIndex:nextSectionIndex]
		section := parseSection(sectionSquares)
		if len(section.Categories) > 0 {
			sections = append(sections, section)
		}
	}

	return sections, nil
}

// parseSection parses a single section from the squares
func parseSection(squares [][]*pageSquare) Section {
	// Find the section name
	sectionName := squares[0][0].Content
	squares = squares[1:]

	// Find the categories indexes
	categoriesIndexes := make([]int, 0)
	for i, square := range squares {
		if len(square) == 3 {
			categoriesIndexes = append(categoriesIndexes, i)
		}
	}

	// Parse categories
	categories := make([]Category, 0)
	for i, categoryIndex := range categoriesIndexes {
		var nextCategoryIndex int
		if i == len(categoriesIndexes)-1 {
			nextCategoryIndex = len(squares)
		} else {
			nextCategoryIndex = categoriesIndexes[i+1]
		}
		categorySquares := squares[categoryIndex:nextCategoryIndex]
		category := parseCategory(categorySquares)
		if len(category.Grades) > 0 {
			categories = append(categories, category)
		}
	}

	return Section{
		Name:       sectionName,
		Categories: categories,
	}
}

// parseCategory parses a single category from the squares
func parseCategory(squares [][]*pageSquare) Category {
	// Find the category name
	categoryName := squares[0][0].Content
	squares = squares[1:]

	// Parse grades
	grades := make([]Grade, 0)
	for _, square := range squares {
		grade, ok := parseGrade(square)
		if ok {
			grades = append(grades, grade)
		}
	}

	return Category{
		Name:   categoryName,
		Grades: grades,
	}
}

// parseGrade parses a single grade from the squares
func parseGrade(squares []*pageSquare) (Grade, bool) {
	// Check if the square is a grade
	if len(squares) != 6 {
		return Grade{}, false
	}

	// Parse coefficient
	coefficient, err := strconv.ParseFloat(squares[3].Content, 64)
	if err != nil {
		return Grade{}, false
	}

	// Parse letter
	letter, previousLetter := parseLetter(squares[2].Content)

	return Grade{
		Name:          squares[0].Content,
		Letter:        letter,
		Coefficient:   coefficient,
		PreviousGrade: previousLetter,
	}, true
}

// parseLetter parses the letter from the grade string
func parseLetter(grade string) (*GradeLetter, *GradeLetter) {
	if grade == "" {
		return nil, nil
	}

	// If single letter, return it
	if len(grade) == 1 {
		return stringToLetter(grade), nil
	}

	// If two letters, split them by /
	sections := strings.SplitN(grade, "/", 2)
	if len(sections) != 2 {
		return nil, nil
	}

	previousLetter := stringToLetter(strings.TrimSpace(sections[0]))
	letter := stringToLetter(strings.TrimSpace(sections[1]))

	return letter, previousLetter
}

// stringToLetter converts a string to a GradeLetter
func stringToLetter(grade string) *GradeLetter {
	var letter GradeLetter
	switch grade {
	case "A":
		letter = A
	case "B":
		letter = B
	case "C":
		letter = C
	case "D":
		letter = D
	default:
		return nil
	}

	return &letter
}
