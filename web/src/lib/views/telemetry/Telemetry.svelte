<script lang="ts">
    import TelemetryStatCard from '$lib/components/telemetry/TelemetryStatCard.svelte';
    import { telemetryState } from '$lib/state/telemetry.svelte.js';
    import TelemetryTimingStatsCard from '$lib/components/telemetry/TelemetryTimingStatsCard.svelte';
    import TelemetryParsesList from '$lib/components/telemetry/TelemetryParsesList.svelte';
    import TelemetryPieStatCard from '$lib/components/telemetry/TelemetryPieStatCard.svelte';

    let userAgents = $derived(
        telemetryState.telemetry
            .map(t => t.userAgent)
            .filter(u => !!u),
    );

    const round = (value: number | null | undefined, decimals: number = 0) => {
        const factor = Math.pow(10, decimals);
        return Math.round((value ?? 0) * factor) / factor;
    };
</script>


<div class="grid sm:gap-4 gap-2 p-2 sm:p-4">
    <div class="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-5 sm:gap-4 gap-2">
        <TelemetryStatCard description="{telemetryState.stats?.total_parses_over_last_week ?? 0} parse over last week"
                           label="Total Records"
                           value={telemetryState.stats?.total_parses ?? 0}/>
        <TelemetryStatCard description="{telemetryState.stats?.errors_over_last_week ?? 0} over last week"
                           label="Unique Users"
                           value={telemetryState.stats?.unique_users ?? 0}/>
        <TelemetryStatCard description="{telemetryState.stats?.unique_users_over_last_week ?? 0} errors over last week"
                           label="Error Rate"
                           value="{round(telemetryState.stats?.error_rate, 1)} %"/>
        <TelemetryStatCard description="Max: {telemetryState.stats?.max_pdf_size_kb} KB"
                           label="Average PDF Size"
                           value="{round(telemetryState.stats?.average_pdf_size)} KB"/>
        <TelemetryStatCard description="95th percentile: {round(telemetryState.stats?.average_parse_time_95th)} ms"
                           label="Average Duration"
                           value="{round(telemetryState.stats?.average_parse_time)} ms"/>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-3 sm:gap-4 gap-2">
        <TelemetryPieStatCard label="Operating System" values={userAgents.map(u => u.os)}/>
        <TelemetryPieStatCard label="Browser" values={userAgents.map(u => u.browser)}/>
        <TelemetryPieStatCard label="Platform" values={userAgents.map(u => u.platform)}/>
    </div>

    <TelemetryTimingStatsCard/>

    <TelemetryParsesList/>
</div>
