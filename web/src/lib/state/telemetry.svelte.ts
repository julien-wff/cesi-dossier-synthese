export interface TelemetryUnit {
    success: boolean;
    timestamp: string;
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
    totalParses: number;
    totalParsesOverLastWeek: number;
    uniqueUsers: number;
    uniqueUsersOverLastWeek: number;
    errorRate: number;
    errorsOverLastWeek: number;
    averagePdfSizeKb: number;
    maxPdfSizeKb: number;
    averageParseTime: number;
    averageParseTime95th: number;
    UAOSs: Record<string, number>;
    UABrowsers: Record<string, number>;
    UAPlatforms: Record<string, number>;
    latestSuccessfulParses: TelemetryUnit[];
    latestFailedParses: TelemetryUnit[];
    timingAverageDuration: Record<string, number>;
    timingKeys: { name: string, description: string }[];
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
