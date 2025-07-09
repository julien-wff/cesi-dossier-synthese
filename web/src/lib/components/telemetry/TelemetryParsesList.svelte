<script lang="ts">
    import { telemetryState } from '$lib/state/telemetry.svelte.js';
    import TelemetryParse from '$lib/components/telemetry/TelemetryParse.svelte';

    let successfulParses = $derived(telemetryState.telemetry.filter(t => t.success));
    let failedParses = $derived(telemetryState.telemetry.filter(t => !t.success));
</script>


<div class="grid sm:gap-4 gap-2 2xl:grid-cols-2">
    <div class="h-fit flex items-center justify-center flex-col bg-slate-100 dark:bg-slate-700 p-2 sm:p-4 rounded-sm shadow-sm">
        <span class="text-center font-bold mb-2 sm:mb-4">Latest successful parses</span>

        <div class="flex flex-col gap-2 w-full">
            {#if successfulParses.length === 0}
                <div class="text-sm text-gray-500 dark:text-gray-400 w-full text-center">
                    No telemetry data available
                </div>
            {/if}

            {#each successfulParses.slice(-25).toReversed() as telemetry}
                <TelemetryParse unit={telemetry}/>
            {/each}
        </div>
    </div>

    <div class="h-fit flex items-center justify-center flex-col bg-slate-100 dark:bg-slate-700 p-2 sm:p-4 rounded-sm shadow-sm">
        <span class="text-center font-bold mb-2 sm:mb-4">Latest failed parses</span>

        <div class="flex flex-col gap-2 w-full">
            {#if failedParses.length === 0}
                <div class="text-sm text-gray-500 dark:text-gray-400 w-full text-center">
                    No telemetry data available
                </div>
            {/if}

            {#each failedParses.slice(-25).toReversed() as telemetry}
                <TelemetryParse unit={telemetry}/>
            {/each}
        </div>
    </div>
</div>
