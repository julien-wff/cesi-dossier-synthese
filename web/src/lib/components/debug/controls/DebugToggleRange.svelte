<script lang="ts">
    import DebugCheckBox from '$lib/components/debug/controls/DebugCheckBox.svelte';

    export let checked: boolean;
    export let label: string;
    export let min = 0;
    export let max = 100;
    export let value = 0;
</script>

<DebugCheckBox {label} bind:checked/>
{#if checked}
    <div class="w-full grid grid-cols-[3rem_1fr] gap-2">
        <input type="number" {min} {max} step="1" bind:value/>
        <input type="range" {min} {max} step="1" bind:value/>
    </div>
{/if}
