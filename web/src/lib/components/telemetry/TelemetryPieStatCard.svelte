<script lang="ts">
    import { onMount } from 'svelte';

    interface Props {
        label: string;
        values: Record<string, number>;
    }

    let { label, values }: Props = $props();

    let sortedValues: Record<string, number> = $derived((() => {
        // Sort the keys in descending order of their counts
        const sortedKeys = Object.keys(values).sort((a, b) => values[b] - values[a]);
        const sortedCounts: Record<string, number> = {};
        for (const key of sortedKeys) {
            sortedCounts[key] = values[key];
        }

        return sortedCounts;
    })());

    let totalDeviceCount = $derived(Object.values(values).reduce((a, b) => a + b, 0));

    const colors = [
        'red',
        'orange',
        'yellow',
        'green',
        'blue',
        'purple',
        'pink',
        'gray',
    ];

    const getTheme = () => window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
    let theme = $state<'light' | 'dark'>(getTheme());

    function updateTheme() {
        theme = getTheme();
    }

    onMount(() => {
        window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', updateTheme);
        return () => {
            window.matchMedia('(prefers-color-scheme: dark)').removeEventListener('change', updateTheme);
        };
    });

    const colorsThemes = $derived(
        colors.map(color => `var(--color-${color}-${theme === 'dark' ? '500' : '300'})`),
    );
</script>

<div class="flex items-center justify-center flex-col bg-slate-100 dark:bg-slate-700 p-2 sm:p-4 rounded-sm shadow-sm gap-2 sm:gap-4">
    <p class="text-center font-bold">{label}</p>
    {#if totalDeviceCount === 0}
        <p class="text-center text-sm text-gray-500">No data available</p>
    {:else}
        <div class="w-full max-w-sm grid grid-cols-2">
            <div class="justify-self-center relative aspect-square w-32">
                {#each Object.entries(sortedValues).slice(0, 8) as [ _, count ], idx}
                    {@const countUntilLast = Object.values(sortedValues).slice(0, idx).reduce((a, b) => a + b, 0)}
                    <div class="absolute w-full h-full inset-0 rounded-full"
                         style="
                             background: conic-gradient({colorsThemes[idx]} calc({count / totalDeviceCount} * 360deg), transparent 0);
                             rotate: {(countUntilLast / totalDeviceCount) * 360}deg;
                         "></div>
                {/each}
            </div>

            <div class="justify-self-center flex flex-col justify-center gap-1">
                {#each Object.entries(sortedValues).slice(0, 8) as [ value, count ], idx}
                    <div class="flex items-center gap-2">
                        <div class="w-4 h-4 rounded-full" style="background-color: {colorsThemes[idx]}"></div>
                        <span class="text-sm">{value || 'Unknown'} ({count})</span>
                    </div>
                {/each}
            </div>
        </div>
    {/if}
</div>
