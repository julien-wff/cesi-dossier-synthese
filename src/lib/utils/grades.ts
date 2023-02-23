import type { Grade } from '../types/grades';

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
