<script lang="ts">
    import { type DebugLineLine, type DebugResponse, DrawMode } from '$lib/types/debug.js';
    import DebugPdfViewer from '$lib/components/debug/DebugPdfViewer.svelte';
    import DebugToggleRange from '$lib/components/debug/controls/DebugToggleRange.svelte';
    import DebugCheckBox from '$lib/components/debug/controls/DebugCheckBox.svelte';

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
    export let showNeighbours: boolean;

    const round = (v: number) => Math.round(v * 100) / 100;

    function handleViewClick(pageNumber: number, drawMode: DrawMode) {
        page = pageNumber;
        mode = drawMode;
    }

    function getNeighboursIndexes(line: DebugLineLine, page: number) {
        return data
            .lines[page]
            .lines
            .map((l, i) => ({ id: l.id, ind: i }))
            .filter(l => line.start_neighbours_ids.includes(l.id) || line.end_neighbours_ids.includes(l.id))
            .map(({ ind }) => ind);
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
<br/>
Squares: {data.squares.reduce((t, p) => t + p.squares.length, 0)}

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
    <DebugToggleRange label="Text picker"
                      bind:checked={displaySingleText}
                      bind:value={textIndex}
                      min={0}
                      max={data.pages[page].text.length - 1}/>
    {#if displaySingleText}
        {@const text = data.pages[page].text[textIndex]}
        <div class="text-sm p-1 rounded-s border border-indigo-400 mb-2 mt-1">
            Position: ({round(text.position.x)} ; {round(text.position.y)})
            <br/>
            Font size: {text.font_size}
            <br/>
            Text: {text.content}
        </div>
    {/if}

    <DebugToggleRange label="Line picker"
                      bind:checked={displaySingleLine}
                      bind:value={lineIndex}
                      min={0}
                      max={data.pages[page].lines.length - 1}/>
    {#if displaySingleLine}
        {@const line = data.pages[page].lines[lineIndex]}
        <div class="text-sm p-1 rounded-s border border-indigo-400 mb-2 mt-1">
            ({round(line.x1)} ; {round(line.y2)}) -> ({round(line.x2)} ; {round(line.y2)})
        </div>
    {/if}

    <DebugToggleRange label="Square picker"
                      bind:checked={displaySingleSquare}
                      bind:value={squareIndex}
                      min={0}
                      max={data.pages[page].rectangles.length - 1}/>
    {#if displaySingleSquare}
        {@const square = data.pages[page].rectangles[squareIndex]}
        <div class="text-sm p-1 rounded-s border border-indigo-400 mb-2 mt-1">
            Position: ({round(square.position.x)} ; {round(square.position.y)})
            -> ({round(square.position.x + square.size.width)} ; {round(square.position.y + square.size.height)})
            <br/>
            Size: {round(square.size.width)} x {round(square.size.height)}
        </div>
    {/if}
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

{#if mode === DrawMode.Line}
    <DebugToggleRange label="Line picker"
                      bind:checked={displaySingleLine}
                      bind:value={lineIndex}
                      min={0}
                      max={data.lines[page].lines.length - 1}/>
    {#if displaySingleLine}
        {@const line = data.lines[page].lines[lineIndex]}
        <div class="text-sm p-1 rounded-s border border-indigo-400 mb-2 mt-1">
            ({round(line.x1)} ; {round(line.y2)}) -> ({round(line.x2)} ; {round(line.y2)})
        </div>
        <DebugCheckBox label="Show neighbours" bind:checked={showNeighbours}/>
        {#if showNeighbours}
            <div class="text-sm p-1 rounded-s border border-indigo-400 mb-2 mt-1">
                {line.start_neighbours_ids.length + line.end_neighbours_ids.length} neighbours
                ({getNeighboursIndexes(line, page).join(', ')})
            </div>
        {/if}
    {/if}
{/if}

<div class="h-[1px] w-full bg-gray-400 my-2"/>
<h2 class="text-lg font-bold mb-2">Squares view</h2>
<div class="grid grid-cols-2 gap-2">
    {#each data.squares as { page: pageNumber }}
        <button class="aspect-square grid place-content-center rounded border-2 cursor-pointer"
                class:border-indigo-400={mode === DrawMode.Square && pageNumber === page}
                on:click={() => handleViewClick(pageNumber, DrawMode.Square)}>
            <DebugPdfViewer margin={0} {data} page={pageNumber} mode={DrawMode.Square} {debugColors}/>
        </button>
    {/each}
</div>

{#if mode === DrawMode.Square}
    <DebugToggleRange label="Square picker"
                      bind:checked={displaySingleSquare}
                      bind:value={squareIndex}
                      min={0}
                      max={data.squares[page].squares.length - 1}/>
    {#if displaySingleSquare}
        {@const square = data.squares[page].squares[squareIndex]}
        <div class="text-sm p-1 rounded-s border border-indigo-400 mb-2 mt-1">
            Position: ({round(square.x1)} ; {round(square.y1)})
            -> ({round(square.x2)} ; {round(square.y2)})
            <br/>
            Size: {round(square.x2 - square.x1)} x {round(square.y2 - square.y1)}
            <br/>
            Content: {square.content}
        </div>
    {/if}
{/if}
