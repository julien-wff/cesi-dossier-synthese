<script lang="ts">
    import DebugCheckBox from '$lib/components/debug/controls/DebugCheckBox.svelte';

    interface Props {
        checked: boolean;
        label: string;
        min?: number;
        max?: number;
        value?: number;
    }

    let {
        checked = $bindable(),
        label,
        min = 0,
        max = 100,
        value = $bindable(0)
    }: Props = $props();
</script>

<DebugCheckBox {label} bind:checked/>
{#if checked}
    <div class="w-full grid grid-cols-[3rem_1fr] gap-2">
        <input type="number" {min} {max} step="1" bind:value class="dark:text-gray-900 px-1"/>
        <input type="range" {min} {max} step="1" bind:value class="accent-indigo-500"/>
    </div>
{/if}
