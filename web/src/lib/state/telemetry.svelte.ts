interface TimingElement {
    name: string;
    description: string;
    duration: number; // in milliseconds
}

interface UserAgent {
    os: string;
    browser: string;
    platform: string;
}

interface TelemetryStats {
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
    UAOSs: Record<string, number> | null;
    UABrowsers: Record<string, number> | null;
    UAPlatforms: Record<string, number> | null;
    latestSuccessfulParses: TelemetryUnit[] | null;
    latestFailedParses: TelemetryUnit[] | null;
    timingAverageDuration: Record<string, number> | null;
    timingKeys: { name: string, description: string }[] | null;
}

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

export interface TelemetryState {
    loaded: boolean;
    stats: TelemetryStats | null;
}

export const telemetryState = $state<TelemetryState>({
    loaded: false,
    stats: null,
});
