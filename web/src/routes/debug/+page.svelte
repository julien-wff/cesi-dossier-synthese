<script lang="ts">
    import '$lib/app.css';
    import Meta from '$lib/components/Meta.svelte';
    import { type DebugResponse, DrawMode } from '$lib/types/debug';
    import DebugPdfViewer from '$lib/components/debug/DebugPdfViewer.svelte';
    import DebugViewController from '$lib/components/debug/DebugViewController.svelte';

    let fileInput = $state<HTMLInputElement>();
    let error = $state<string | null>(null);
    let loading = $state(false);
    let data = $state<DebugResponse | null>(null);

    // Viewer
    let debugColors = $state(true);
    let mode = $state(DrawMode.Page);
    let page = $state(0);

    let displaySingleText = $state(false);
    let textIndex = $state(0);
    let displaySingleLine = $state(false);
    let lineIndex = $state(0);
    let displaySingleSquare = $state(false);
    let squareIndex = $state(0);
    let showNeighbours = $state(false);

    async function handlePDFSubmit(ev: SubmitEvent | Event) {
        if (ev instanceof SubmitEvent)
            ev.preventDefault();

        error = null;
        data = null;

        const file = fileInput?.files?.[0];
        if (!file) {
            if (ev instanceof SubmitEvent)
                error = 'No file selected';
            return;
        }

        const form = new FormData();
        form.append('file', file);

        try {
            loading = true;
            const res = await fetch('/api/parse/debug', {
                method: 'POST',
                body: form,
            });

            if (res.status !== 200) {
                throw new Error(`Failed to parse PDF: ${await res.text()}`);
            }

            data = await res.json();
        } catch (e) {
            console.error(e);
            error = (e as Error).message;
        } finally {
            loading = false;
        }
    }
</script>


<Meta/>

<main class="flex min-h-screen">
    <aside class="w-72 p-2 bg-slate-100 dark:bg-slate-700 shadow-md max-h-screen overflow-y-auto">
        <h1 class="text-xl font-bold mb-4">Debug viewer</h1>

        <form class="mb-4" onsubmit={handlePDFSubmit}>
            <input accept="application/pdf"
                   bind:this={fileInput}
                   class="block w-full text-sm disabled:opacity-50 disabled:cursor-not-allowed"
                   disabled={loading}
                   onchange={handlePDFSubmit}
                   type="file">
            <button class="w-full bg-indigo-500 text-white py-1 mt-2 rounded-sm disabled:opacity-50 disabled:cursor-not-allowed"
                    disabled={loading}
                    type="submit">
                Parse PDF
            </button>
            {#if error}
                <p class="text-red-500 text-sm mt-2">{error}</p>
            {/if}
        </form>

        {#if data}
            <DebugViewController bind:data
                                 bind:debugColors
                                 bind:mode
                                 bind:page
                                 bind:displaySingleText
                                 bind:textIndex
                                 bind:displaySingleLine
                                 bind:lineIndex
                                 bind:displaySingleSquare
                                 bind:squareIndex
                                 bind:showNeighbours/>
        {/if}
    </aside>

    <div class="flex-1 grid place-content-center">
        <DebugPdfViewer {data}
                        {debugColors}
                        {displaySingleLine}
                        {displaySingleSquare}
                        {displaySingleText}
                        {lineIndex} {mode}
                        {page} resolution={2}
                        {showNeighbours} {squareIndex}
                        {textIndex}/>
    </div>
</main>
