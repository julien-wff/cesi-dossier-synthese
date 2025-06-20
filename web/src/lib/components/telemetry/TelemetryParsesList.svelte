<script lang="ts">
    import { telemetryState } from '$lib/state/telemetry.svelte.js';
    import TelemetryParse from '$lib/components/telemetry/TelemetryParse.svelte';
</script>


<div class="flex items-center justify-center flex-col bg-slate-100 dark:bg-slate-700 p-2 sm:p-4 rounded-sm shadow-sm">
    <span class="text-center font-bold mb-2 sm:mb-4">Latest parses</span>

    <div class="flex flex-col gap-2 w-full">
    {#each telemetryState.telemetry.toReversed().slice(0, 25) as telemetry}
        <TelemetryParse unit={telemetry}/>
    {/each}
    </div>
</div>
