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

export interface TimingElement {
    name: string;
    description: string;
    duration: number; // in milliseconds
}

export interface TelemetryState {
    loaded: boolean;
    telemetry: TelemetryUnit[];
}

export interface UserAgent {
    os: string;
    browser: string;
    platform: string;
}

export const telemetryState = $state<TelemetryState>({
    loaded: false,
    telemetry: [],
});
