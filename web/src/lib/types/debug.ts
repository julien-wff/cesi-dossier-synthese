export interface DebugResponse {
    performance: DebugPerformanceElement[];
    pages: DebugPage[];
    lines: DebugLine[];
    squares: DebugSquare[];
}

interface DebugPerformanceElement {
    name: string;
    duration: number; // ms
    description: string;
}

interface DebugPage {
    page: number;
    size: DebugSize;
    text: DebugPageText[];
    rectangles: DebugPageRectangle[];
    lines: DebugPageLine[];
}

interface DebugSize {
    width: number;
    height: number;
}

interface DebugPosition {
    x: number;
    y: number;
}

interface DebugPageText {
    content: string;
    font_size: number;
    position: DebugPosition;
}

interface DebugPageRectangle {
    position: DebugPosition;
    size: DebugSize;
}

interface DebugPageLine {
    x1: number;
    y1: number;
    x2: number;
    y2: number;
}

interface DebugLine {
    page: number;
    lines: DebugLineLine[];
}

export enum DebugLineDirection {
    Horizontal = 'horizontal',
    Vertical = 'vertical',
}

export interface DebugLineLine {
    id: number;
    x1: number;
    y1: number;
    x2: number;
    y2: number;
    direction: DebugLineDirection;
    length: number;
    start_neighbours_ids: number[];
    end_neighbours_ids: number[];
}

interface DebugSquare {
    page: number;
    squares: DebugSquareSquare[][];
}

interface DebugSquareSquare {
    x1: number;
    y1: number;
    x2: number;
    y2: number;
    content: string;
}

export enum DrawMode {
    Page,
    Line,
    Square,
}
