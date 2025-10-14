<script lang="ts">
    import TelemetryStatCard from '$lib/components/telemetry/TelemetryStatCard.svelte';
    import { telemetryState } from '$lib/state/telemetry.svelte.js';
    import TelemetryTimingStatsCard from '$lib/components/telemetry/TelemetryTimingStatsCard.svelte';
    import TelemetryParsesList from '$lib/components/telemetry/TelemetryParsesList.svelte';
    import TelemetryPieStatCard from '$lib/components/telemetry/TelemetryPieStatCard.svelte';

    const round = (value: number | null | undefined, decimals: number = 0) => {
        const factor = Math.pow(10, decimals);
        return Math.round((value ?? 0) * factor) / factor;
    };
</script>


<div class="grid sm:gap-4 gap-2 p-2 sm:p-4">
    <div class="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-5 sm:gap-4 gap-2">
        <TelemetryStatCard description="{telemetryState.stats?.totalParsesOverLastWeek ?? 0} parse over last week"
                           label="Total Records"
                           value={telemetryState.stats?.totalParses ?? 0}/>
        <TelemetryStatCard description="{telemetryState.stats?.uniqueUsersOverLastWeek ?? 0} over last week"
                           label="Unique Users"
                           value={telemetryState.stats?.uniqueUsers ?? 0}/>
        <TelemetryStatCard description="{telemetryState.stats?.errorsOverLastWeek ?? 0} errors over last week"
                           label="Error Rate"
                           value="{round(telemetryState.stats?.errorRate, 1)} %"/>
        <TelemetryStatCard description="Max: {telemetryState.stats?.maxPdfSizeKb} KB"
                           label="Average PDF Size"
                           value="{round(telemetryState.stats?.averagePdfSizeKb)} KB"/>
        <TelemetryStatCard description="95th percentile: {round(telemetryState.stats?.averageParseTime95th)} ms"
                           label="Average Duration"
                           value="{round(telemetryState.stats?.averageParseTime)} ms"/>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-3 sm:gap-4 gap-2">
        <TelemetryPieStatCard label="Operating System" values={telemetryState.stats?.UAOSs ?? {}}/>
        <TelemetryPieStatCard label="Browser" values={telemetryState.stats?.UABrowsers ?? {}}/>
        <TelemetryPieStatCard label="Platform" values={telemetryState.stats?.UAPlatforms ?? {}}/>
    </div>

    <TelemetryTimingStatsCard/>

    <TelemetryParsesList/>
</div>
