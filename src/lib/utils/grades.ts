import type { Grade, Section } from '../types/grades';

export const LETTER_VALUES = {
    A: 5,
    B: 4,
    C: 2,
    D: 1,
};

export const convertToLetterGrade = (grade: string) => LETTER_VALUES[grade as keyof typeof LETTER_VALUES] || 0;

export function calculateAverage(grades: string[]) {
    grades = grades.filter(convertToLetterGrade);
    const total = grades.reduce((acc, grade) => acc + convertToLetterGrade(grade), 0);
    return total / grades.length;
}

export function validationLevel(value: number) {
    if (value > 4.5)
        return 2;
    if (value >= 3.6)
        return 1;
    if (value > 0)
        return 0;
    return null;
}

export function gradesWithCoefficientToList(grades: Grade[]) {
    const letters: string[] = [];
    for (const { grade, coefficient } of grades) {
        if (grade && convertToLetterGrade(grade))
            letters.push(...new Array(coefficient).fill(grade));
    }
    return letters;
}

export function sectionsToGrades(sections: Section[]) {
    return sections.map(({ categories }) =>
        categories.map(({ grades }) =>
            gradesWithCoefficientToList(grades),
        ),
    );
}

export function calculateGradesAverage(sections: Section[]) {
    const grades = sectionsToGrades(sections).flat(2);
    return calculateAverage(grades);
}

export function calculateCategoriesAverage(sections: Section[]) {
    const categoryGrades = sectionsToGrades(sections).map(category => category.flat());
    const categoryAverages = categoryGrades
        .map(calculateAverage)
        .filter(average => !isNaN(average));
    return categoryAverages.reduce((acc, cat) => acc + cat, 0) / categoryAverages.length;
}
