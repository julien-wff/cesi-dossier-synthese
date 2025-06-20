<script lang="ts">
    import type { TelemetryUnit } from '$lib/state/telemetry.svelte';
    import CircleCheck from 'lucide-svelte/icons/circle-check';
    import CircleX from 'lucide-svelte/icons/circle-x';

    interface Props {
        unit: TelemetryUnit;
    }

    let { unit }: Props = $props();

    let pageCount = $derived(unit.timings.filter(t => t.name.startsWith('extract-page-'))?.length ?? 0);
    let duration = $derived(unit.timings.reduce((acc, t) => acc + t.duration, 0));
</script>

<div class="bg-slate-50 dark:bg-slate-800 p-2 sm:p-4 rounded-sm shadow-sm flex gap-2">
    {#if unit.success}
        <CircleCheck class="text-green-500 dark:text-green-400"/>
    {:else}
        <CircleX class="text-red-500 dark:text-red-400"/>
    {/if}

    <div>[{unit.timestamp.toLocaleString()}] {unit.source} ({unit.contentLengthKB} KB)</div>

    -

    {#if unit.success}
        <div>
            Read {pageCount} page{pageCount === 1 ? '' : 's'} in {duration.toFixed(2)} ms
        </div>
    {:else}
        <div>{unit.error}</div>
    {/if}
</div>
