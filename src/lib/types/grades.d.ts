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
    grade: string | null;
    coefficient: number;
    previousGrade?: string | null;
}
