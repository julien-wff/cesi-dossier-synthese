package parser

import (
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"testing"
)

func TestParsePdf(t *testing.T) {
	tests := []struct {
		name                   string
		file                   string
		sectionNames           []string
		sectionCategoriesCount []int
		categoriesNames        []string
		categoriesGradesCount  []int
	}{
		{
			name: "should parse S6 empty PDF",
			file: "julien-s6-empty.PDF",
			sectionNames: []string{
				"Sciences de Base",
				"Sciences et Méthodes de l'Ingénieur",
				"Sciences et Techniques de la Spécialité",
				"Sciences Humaines, Economiques, Juridiques et Sociales",
				"Langues",
				"Missions en Entreprise",
			},
			sectionCategoriesCount: []int{2, 1, 1, 1, 1, 1},
			categoriesNames: []string{
				"[S6] 6.1 Mathématiques pour l'ingénieur (Semestre 6)",
				"[S6] 6.2 Sciences pour l'ingénieur en informatique",
				"[S6] 6.3 Méthodes et outils de l'ingénieur (Semestre 6)",
				"[S6] 6.4 Stockage et traitement des données",
				"[S6] 6.5 Engagement de l'ingénieur et communication",
				"[S6] 6.6 Anglais (Semestre 6)",
				"[S6] 6.7 Stage en entreprise (Semestre 6)",
			},
			categoriesGradesCount: []int{4, 3, 5, 4, 7, 4, 4},
		},
		{
			name: "should parse S5 partial PDF",
			file: "julien-s5-partial.PDF",
			sectionNames: []string{
				"Sciences de Base",
				"Sciences et Méthodes de l'Ingénieur",
				"Sciences et Techniques de la Spécialité",
				"Sciences Humaines, Economiques, Juridiques et Sociales",
				"Langues",
			},
			sectionCategoriesCount: []int{2, 1, 1, 1, 1},
			categoriesNames: []string{
				"[S5] 5.1 Mathématiques pour l'ingénieur (Semestre 5)",
				"[S5] 5.2 Architecture des systèmes",
				"[S5] 5.4 Génie logiciel",
				"[S5] 5.3 Sécurité et Gestion du Système d'Information",
				"[S5] 5.5 Engagement de l'ingénieur et communication",
				"[S5] 5.6 Anglais (Semestre 5)",
			},
			categoriesGradesCount: []int{5, 5, 5, 4, 6, 5},
		},
		{
			name: "should parse S3 full PDF",
			file: "julien-s3-full.PDF",
			sectionNames: []string{
				"Sciences de Base",
				"Sciences et Méthodes de l'Ingénieur",
				"Sciences et Techniques de la Spécialité",
				"Sciences Humaines, Economiques, Juridiques et Sociales",
			},
			sectionCategoriesCount: []int{1, 2, 3, 2},
			categoriesNames: []string{
				"[S3] Algorithmes et langages",
				"[S3] Génie logiciel",
				"[S3] Méthodes et outils pour l'entreprise",
				"[S3] Bases en programmation",
				"[S3] Architecture du réseau local",
				"[S3] Bases pour l'administration système et réseau",
				"[S3] Comportement de l'ingénieur (Semestre 3)",
				"[S3] Anglais (Semestre 3)",
			},
			categoriesGradesCount: []int{5, 5, 5, 5, 5, 5, 7, 4},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := os.Open("../../testdata/" + tt.file)
			assert.Nil(t, err)

			// convert *File to io.ReadSeeker
			var r io.ReadSeeker
			r = f

			res, _, apiErr := ParsePdf(&r)
			assert.Nil(t, apiErr)

			// Sections
			assert.Len(t, res.Data, len(tt.sectionNames))
			for i, name := range tt.sectionNames {
				assert.Equal(t, name, res.Data[i].Name)
			}

			// Categories
			for i, count := range tt.sectionCategoriesCount {
				assert.Len(t, res.Data[i].Categories, count)
			}

			categories := make([]Category, 0)
			for _, section := range res.Data {
				categories = append(categories, section.Categories...)
			}

			for i, name := range tt.categoriesNames {
				assert.Equal(t, name, categories[i].Name)
			}

			// Grades
			for i, count := range tt.categoriesGradesCount {
				assert.Len(t, categories[i].Grades, count)
			}
		})
	}
}
