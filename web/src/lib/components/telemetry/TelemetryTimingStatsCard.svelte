<script lang="ts">
    import { telemetryState } from '$lib/state/telemetry.svelte.js';

    // Total average duration across all timings
    const totalAverageDuration = $derived(
        Object.values(telemetryState.stats?.timingAverageDuration ?? {}).reduce((sum, value) => sum + value, 0),
    );
</script>


<div class="flex items-center justify-center flex-col bg-slate-100 dark:bg-slate-700 p-2 sm:p-4 rounded-sm shadow-sm">
    <div class="text-center font-bold mb-2 md:mb-4">Timing Statistics</div>

    <div class="flex flex-col md:flex-row flex-1 w-full gap-0.5">
        {#if !telemetryState.stats || !telemetryState.stats.timingKeys || telemetryState.stats.timingKeys.length === 0}
            <div class="text-sm text-gray-500 dark:text-gray-400 w-full text-center">
                No timing data available
            </div>
        {:else}
            {#each telemetryState.stats.timingKeys as { name, description }, idx (name)}
                {@const stats = telemetryState.stats.timingAverageDuration![name]}
                {@const duration = `${stats.toFixed(2)}ms`}
                {@const fraction = (stats / totalAverageDuration * 100).toFixed(2)}

                <div class="text-sm py-4 text-center md:w-(--fraction) h-(--fraction) md:h-auto flex md:flex-col gap-1 items-center justify-center overflow-clip"
                     class:bg-red-300={idx === 0}
                     class:bg-orange-300={idx === 1}
                     class:bg-yellow-300={idx === 2}
                     class:bg-green-300={idx === 3}
                     class:bg-blue-300={idx === 4}
                     class:bg-purple-300={idx === 5}
                     class:bg-pink-300={idx === 6}
                     class:bg-gray-300={idx >= 7}
                     class:dark:bg-red-500={idx === 0}
                     class:dark:bg-orange-500={idx === 1}
                     class:dark:bg-yellow-500={idx === 2}
                     class:dark:bg-green-500={idx === 3}
                     class:dark:bg-blue-500={idx === 4}
                     class:dark:bg-purple-500={idx === 5}
                     class:dark:bg-pink-500={idx === 6}
                     class:dark:bg-gray-500={idx >= 7}
                     title="{description} ({duration})"
                     style="--fraction: {fraction}%">
                    <div class="line-clamp-1">{description}</div>
                    <div class="line-clamp-1">({duration})</div>
                </div>
            {/each}
        {/if}
    </div>
</div>
