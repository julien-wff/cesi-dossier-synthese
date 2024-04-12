import type { Letter, Section } from '../types/grades';
import type { ExtractionTable } from '../types/tabula';

const formatString = (str: string) => str.trim().replace(/[\n\r]+/g, ' ');

export function parseTabulaResult(input: ExtractionTable[]): Section[] {

    // Filter out tables that don't have the correct header
    const tables = input.filter(table => {
        const headers = table.data[0].map(cell => formatString(cell.text));
        return headers.includes('AXE / U.E. / E.E.') && headers.includes('Evaluation') && headers.includes('Coeff.');
    });

    // Convert the tables to a temporary format
    const rawTables = tables
        .map(table => table.data.slice(1))
        .map(table => table.map(row => ({
            'AXE / U.E. / E.E.': formatString(row[0].text),
            Evaluation: formatString(row[2].text),
            Coeff: formatString(row[3].text),
        }) satisfies RawGradeTable));

    // Merge the tables into a single array of grades
    const grades = rawTables.reduce((acc, table) => [ ...acc, ...table ], [] as RawGradeTable[]);

    const result: Section[] = [];

    for (const row of grades) {
        const title = row['AXE / U.E. / E.E.'];
        const coefficient = row.Coeff || null;

        let letter: Letter | null = null, previousGrade: Letter | null = null;
        if ([ 'A', 'B', 'C', 'D' ].includes(row.Evaluation)) {
            letter = row.Evaluation as Letter;
        } else if (RegExp(/^[A-D]( \/ [A-D])?$/).exec(row.Evaluation)) {
            letter = row.Evaluation.split(' / ')[1] as Letter;
            previousGrade = row.Evaluation.split(' / ')[0] as Letter;
        }

        if (!coefficient) {
            if (!title.startsWith('[')) {
                result.push({
                    name: title,
                    categories: [],
                });
            } else {
                result.at(-1)!.categories.push({
                    name: title,
                    grades: [],
                });
            }
        } else {
            result.at(-1)!.categories.at(-1)!.grades.push({
                name: title,
                letter: letter,
                previousGrade,
                coefficient: Number(coefficient),
            });
        }
    }

    return result;
}


interface RawGradeTable {
    'AXE / U.E. / E.E.': string;
    Evaluation: string;
    Coeff: string;
}
