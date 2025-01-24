<script lang="ts">
    import TelemetryStatCard from '$lib/components/telemetry/TelemetryStatCard.svelte';
    import { telemetryState } from '$lib/state/telemetry.svelte.js';

    const isAWeekOld = (timestamp: Date) => timestamp.getTime() > Date.now() - 7 * 24 * 60 * 60 * 1000;

    let totalRecords = $derived(telemetryState.telemetry.length);
    let lastWeekRecords = $derived(
        telemetryState.telemetry.filter(t => isAWeekOld(t.timestamp)).length,
    );

    let totalErrors = $derived(
        telemetryState.telemetry.filter((t) => !t.success).length,
    );
    let errorRate = $derived(totalErrors / totalRecords * 100);
    let errorsOverLastWeek = $derived(
        telemetryState.telemetry.filter((t) => !t.success && isAWeekOld(t.timestamp)).length,
    );

    let averagePdfSize = $derived(
        telemetryState.telemetry.reduce((acc, t) => acc + t.contentLengthKB, 0) / totalRecords,
    );
    let maxPdfSize = $derived(
        telemetryState.telemetry.reduce((acc, t) => Math.max(acc, t.contentLengthKB), 0),
    );

    let averageParseDuration = $derived(
        telemetryState
            .telemetry
            .filter(t => t.success)
            .reduce(
                (acc, t) => acc + t.timings.reduce((acc, t) => acc + t.duration, 0),
                0,
            ) / (totalRecords - totalErrors),
    );
    let percentile95ParseValues = $derived(
        telemetryState.telemetry
            .filter(t => t.success)
            .map(t => t.timings.reduce((acc, t) => acc + t.duration, 0))
            .sort((a, b) => a - b)
            .slice(Math.floor((totalRecords - totalErrors) * 0.95)),
    );
    let percentile95ParseDuration = $derived(
        percentile95ParseValues.reduce((acc, t) => acc + t, 0) / percentile95ParseValues.length,
    );
</script>


<div class="p-2 sm:p-4">
    <div class="grid grid-cols-2 lg:grid-cols-4 sm:gap-4 gap-2">
        <TelemetryStatCard description="{lastWeekRecords} parse over last week"
                           label="Total Records"
                           value={totalRecords}/>
        <TelemetryStatCard description="{errorsOverLastWeek} errors over last week"
                           label="Error Rate"
                           value={errorRate + ' %'}/>
        <TelemetryStatCard description="Max: {maxPdfSize} KB"
                           label="Average PDF Size"
                           value="{Math.round(averagePdfSize)} KB"/>
        <TelemetryStatCard description="95th percentile: {Math.round(percentile95ParseDuration)} ms"
                           label="Average Duration"
                           value="{Math.round(averageParseDuration)} ms"/>
    </div>
</div>
