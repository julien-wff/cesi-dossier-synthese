import type { Grade, Letter, Section } from '../types/grades';

export const LETTER_VALUES = {
    A: 5,
    B: 4,
    C: 2,
    D: 1,
} satisfies {[key in Letter]: number};

export const convertToLetterGrade = (letter: Letter) => LETTER_VALUES[letter] || 0;

export function calculateAverage(letters: Letter[]) {
    letters = letters.filter(convertToLetterGrade);
    const total = letters.reduce((acc, grade) => acc + convertToLetterGrade(grade), 0);
    return total / letters.length;
}

export function validationLevel(grade: number) {
    if (grade > 4.5)
        return 2;
    if (grade >= 3.6)
        return 1;
    if (grade > 0)
        return 0;
    return null;
}

export function gradesWithCoefficientToList(grades: Grade[]) {
    const letters: Letter[] = [];
    for (const { letter, coefficient } of grades) {
        if (letter && convertToLetterGrade(letter))
            letters.push(...new Array(coefficient).fill(letter));
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

export function countLetterWithCoeff(sections: Section[], letter: Letter) {
    const grades = sectionsToGrades(sections).flat(2);
    return grades.filter(grade => grade === letter).length;
}

export function countLetterWithoutCoeff(sections: Section[], letter: Letter) {
    const categoryGrades = sections.map(({ categories }) =>
        categories.map(({ grades }) =>
            grades.map(({ letter }) => letter)))
        .flat(2)
        .filter(grade => grade === letter);
    return categoryGrades.length;
}
