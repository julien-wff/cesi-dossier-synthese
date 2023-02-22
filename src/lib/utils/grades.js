export const LETTER_VALUES = {
    A: 5,
    B: 4,
    C: 2,
    D: 1,
};

export const convertToLetterGrade = (grade) => LETTER_VALUES[grade] || 0;

export function calculateAverage(grades) {
    grades = grades.filter(grade => LETTER_VALUES[grade]);
    const total = grades.reduce((acc, grade) => acc + convertToLetterGrade(grade), 0);
    return total / grades.length;
}

export function validationLevel(value) {
    if (value > 4.5)
        return 2;
    if (value >= 3.6)
        return 1;
    if (value > 0)
        return 0;
    return null;
}

export function gradesWithCoefficientToList(grades) {
    const letters = [];
    for (const { grade, coefficient } of grades) {
        if (LETTER_VALUES[grade])
            letters.push(...new Array(coefficient).fill(grade));
    }
    return letters;
}
