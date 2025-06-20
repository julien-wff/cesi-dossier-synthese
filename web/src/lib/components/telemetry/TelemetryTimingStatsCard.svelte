<script lang="ts">
    import { telemetryState } from '$lib/state/telemetry.svelte.js';

    // Group the durations of the timings by their key
    const accumulatedStats = $derived(telemetryState.telemetry.reduce(
        function (accParse, parse) {
            if (!parse.success || !parse.timings) {
                return accParse;
            }

            return parse.timings.reduce((acc, timing) => {
                if (acc[timing.name] === undefined) {
                    acc[timing.name] = [];
                }
                acc[timing.name].push(timing.duration);
                return acc;
            }, accParse);
        },
        {} as Record<string, number[]>,
    ));

    // Make the average of the durations for each key
    const accumulatedStatsAverage = $derived(Object.entries(accumulatedStats).reduce(
        (acc, [ key, values ]) => {
            acc[key] = values.reduce((sum, value) => sum + value, 0) / values.length;
            return acc;
        },
        {} as Record<string, number>,
    ));

    // Total average duration across all timings
    const totalAverageDuration = $derived(
        Object.values(accumulatedStatsAverage).reduce((sum, value) => sum + value, 0),
    );

    // Take the element with the most keys to get the order of the keys, along with their description
    const timingKeys = $derived(telemetryState
        .telemetry
        .filter(t => t.success && t.timings)
        .sort((a, b) => b.timings.length - a.timings.length)
        [0]
        ?.timings ?? [],
    );
</script>


<div class="flex items-center justify-center flex-col bg-slate-100 dark:bg-slate-700 p-2 sm:p-4 rounded-sm shadow-sm">
    <div class="text-center font-bold mb-2 md:mb-4">Timing Statistics</div>

    <div class="flex flex-1 w-full gap-0.5">
        {#if timingKeys.length === 0}
            <div class="text-sm text-gray-500 dark:text-gray-400 w-full text-center">
                No timing data available
            </div>
        {/if}

        {#each timingKeys as { name, description }, idx (name)}
            {@const duration = `${accumulatedStatsAverage[name].toFixed(2)}ms`}

            <div class="text-sm py-4 text-center"
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
                 style="width: {(accumulatedStatsAverage[name] / totalAverageDuration * 100).toFixed(2)}%">
                <div class="line-clamp-1 mb-1">{description}</div>
                <div class="line-clamp-1">({duration})</div>
            </div>
        {/each}
    </div>
</div>
