export type Letter = 'A' | 'B' | 'C' | 'D';

export interface Section {
    name: string;
    categories: Category[];
}

export interface Category {
    name: string;
    grades: Grade[];
}

export interface Grade {
    name: string;
    letter: Letter | null;
    coefficient: number;
    previousGrade?: Letter | null;
}
