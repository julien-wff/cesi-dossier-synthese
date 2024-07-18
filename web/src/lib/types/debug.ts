export interface DebugResponse {
    performance: DebugPerformanceElement[];
    pages: DebugPage[];
    lines: DebugLine[];
    squares: DebugSquare[];
}

export interface DebugPerformanceElement {
    name: string;
    duration: number; // ms
    description: string;
}

export interface DebugPage {
    page: number;
    size: DebugSize;
    text: DebugPageText[];
    rectangles: DebugPageRectangle[];
    lines: DebugPageLine[];
}

export interface DebugSize {
    width: number;
    height: number;
}

export interface DebugPosition {
    x: number;
    y: number;
}

export interface DebugPageText {
    content: string;
    font_size: number;
    position: DebugPosition;
}

export interface DebugPageRectangle {
    position: DebugPosition;
    size: DebugSize;
}

export interface DebugPageLine {
    x1: number;
    y1: number;
    x2: number;
    y2: number;
}

export interface DebugLine {
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

export interface DebugSquare {
    page: number;
    squares: DebugSquareSquare[];
}

export interface DebugSquareSquare {
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
