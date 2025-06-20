export interface TelemetryUnit {
    success: boolean;
    timestamp: Date;
    clientIP: string;
    source: string;
    contentLengthKB: number;
    timings: TimingElement[];
    error: string | null;
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

export const telemetryState = $state<TelemetryState>({
    loaded: false,
    telemetry: [],
});
