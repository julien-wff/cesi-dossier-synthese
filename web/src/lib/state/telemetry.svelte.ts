export interface TelemetryUnit {
    success: boolean;
    timestamp: Date;
    clientIP: string;
    source: string;
    contentLengthKB: number;
    timings: TimingElement[];
    error: string | null;
    userAgent?: UserAgent;
}

export interface UserAgent {
    os: string;
    browser: string;
    platform: string;
}

export interface TelemetryStats {
    total_parses: number;
    total_parses_over_last_week: number;
    unique_users: number;
    unique_users_over_last_week: number;
    error_rate: number;
    errors_over_last_week: number;
    average_pdf_size: number;
    max_pdf_size_kb: number;
    average_parse_time: number;
    average_parse_time_95th: number;
}

export interface TimingElement {
    name: string;
    description: string;
    duration: number; // in milliseconds
}

export interface TelemetryState {
    loaded: boolean;
    telemetry: TelemetryUnit[];
    stats: TelemetryStats | null;
}

export const telemetryState = $state<TelemetryState>({
    loaded: false,
    telemetry: [],
    stats: null,
});
