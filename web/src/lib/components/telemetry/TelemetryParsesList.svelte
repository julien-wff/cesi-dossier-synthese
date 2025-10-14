<script lang="ts">
    import { telemetryState, type TelemetryUnit } from '$lib/state/telemetry.svelte.js';
    import TelemetryParse from '$lib/components/telemetry/TelemetryParse.svelte';
</script>


{#snippet parsesList(title: string, telemetry: TelemetryUnit[])}
    <div class="h-fit flex items-center justify-center flex-col bg-slate-100 dark:bg-slate-700 p-2 sm:p-4 rounded-sm shadow-sm">
        <span class="text-center font-bold mb-2 sm:mb-4">{title}</span>

        <div class="flex flex-col gap-2 w-full">
            {#if telemetry.length === 0}
                <div class="text-sm text-gray-500 dark:text-gray-400 w-full text-center">
                    No telemetry data available
                </div>
            {/if}

            {#each telemetry as unit}
                <TelemetryParse {unit}/>
            {/each}
        </div>
    </div>
{/snippet}


<div class="grid sm:gap-4 gap-2 2xl:grid-cols-2">
    {@render parsesList('Latest successful parses', telemetryState.stats?.latestSuccessfulParses || [])}
    {@render parsesList('Latest failed parses', telemetryState.stats?.latestFailedParses || [])}
</div>
