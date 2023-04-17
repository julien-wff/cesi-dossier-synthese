export interface ExtractionTable {
    extraction_method: string;
    top: number;
    left: number;
    width: number;
    height: number;
    right: number;
    bottom: number;
    data: ExtractionRow[][];
}

interface ExtractionRow {
    top: number;
    left: number;
    width: number;
    height: number;
    text: string;
}
