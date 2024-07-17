<script lang="ts">
    import { type DebugResponse, DrawMode } from '$lib/types/debug.js';
    import DebugPdfViewer from '$lib/components/debug/DebugPdfViewer.svelte';
    import DebugToggleRange from '$lib/components/debug/controls/DebugToggleRange.svelte';

    export let data: DebugResponse;
    export let debugColors: boolean;
    export let page: number;
    export let mode: DrawMode;

    export let displaySingleText: boolean;
    export let textIndex: number;
    export let displaySingleLine: boolean;
    export let lineIndex: number;
    export let displaySingleSquare: boolean;
    export let squareIndex: number;

    function handleViewClick(pageNumber: number, drawMode: DrawMode) {
        page = pageNumber;
        mode = drawMode;
    }
</script>

<div class="h-[1px] w-full bg-gray-400 mb-2"/>
<h2 class="text-lg font-bold">Statistics</h2>
Duration: {Math.round(data.performance.reduce((t, v) => t + v.duration, 0) * 100) / 100} ms
<br/>
Pages: {data.pages.length}
<br/>
Text nodes: {data.pages.reduce((t, p) => t + p.text.length, 0)}
<br/>
Lines: {data.lines.reduce((t, p) => t + p.lines.length, 0)}

<div class="h-[1px] w-full bg-gray-400 my-2"/>
<h2 class="text-lg font-bold mb-2">Global parameters</h2>
<label class="block select-none">
    <input type="checkbox" bind:checked={debugColors} class="mr-2"/>
    Debug colors
</label>

<div class="h-[1px] w-full bg-gray-400 my-2"/>
<h2 class="text-lg font-bold mb-2">Page view</h2>
<div class="grid grid-cols-2 gap-2">
    {#each data.pages as { page: pageNumber }}
        <button class="aspect-square grid place-content-center rounded border-2 cursor-pointer mb-4"
                class:border-indigo-400={mode === DrawMode.Page && pageNumber === page}
                on:click={() => handleViewClick(pageNumber, DrawMode.Page)}>
            <DebugPdfViewer margin={0} {data} page={pageNumber} mode={DrawMode.Page} {debugColors}/>
        </button>
    {/each}
</div>

{#if mode === DrawMode.Page}
    <DebugToggleRange label="Text picker" bind:checked={displaySingleText} bind:value={textIndex} min={0}
                      max={data.pages[page].text.length - 1}/>
    <DebugToggleRange label="Line picker" bind:checked={displaySingleLine} bind:value={lineIndex} min={0}
                      max={data.pages[page].lines.length - 1}/>
    <DebugToggleRange label="Square picker" bind:checked={displaySingleSquare} bind:value={squareIndex} min={0}
                      max={data.pages[page].rectangles.length - 1}/>
{/if}

<div class="h-[1px] w-full bg-gray-400 my-2"/>
<h2 class="text-lg font-bold mb-2">Lines view</h2>
<div class="grid grid-cols-2 gap-2">
    {#each data.lines as { page: pageNumber }}
        <button class="aspect-square grid place-content-center rounded border-2 cursor-pointer"
                class:border-indigo-400={mode === DrawMode.Line && pageNumber === page}
                on:click={() => handleViewClick(pageNumber, DrawMode.Line)}>
            <DebugPdfViewer margin={0} {data} page={pageNumber} mode={DrawMode.Line} {debugColors}/>
        </button>
    {/each}
</div>