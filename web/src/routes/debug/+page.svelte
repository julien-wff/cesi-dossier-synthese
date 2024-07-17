<script lang="ts">
    import '$lib/app.css';
    import Meta from '$lib/components/Meta.svelte';
    import { type DebugResponse, DrawMode } from '$lib/types/debug';
    import DebugPdfViewer from '$lib/components/debug/DebugPdfViewer.svelte';
    import DebugViewController from '$lib/components/debug/DebugViewController.svelte';

    let fileInput: HTMLInputElement;
    let error: string | null = null;
    let loading = false;
    let data: DebugResponse | null = null;

    // Viewer
    let debugColors = false;
    let mode: DrawMode = DrawMode.Page;
    let page: number = 0;

    let displaySingleText = false;
    let textIndex = 0;
    let displaySingleLine = false;
    let lineIndex = 0;
    let displaySingleSquare = false;
    let squareIndex = 0;

    async function handlePDFSubmit(ev: SubmitEvent | Event) {
        error = null;
        data = null;

        const file = fileInput.files?.[0];
        if (!file) {
            if (ev instanceof SubmitEvent)
                error = 'No file selected';
            return;
        }

        const form = new FormData();
        form.append('file', file);

        try {
            loading = true;
            const res = await fetch('http://localhost:8080/api/parse/debug', {
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
    <aside class="w-72 p-2 bg-indigo-50 shadow-md max-h-screen overflow-y-auto">
        <h1 class="text-xl font-bold mb-4">Debug viewer</h1>

        <form on:submit|preventDefault={handlePDFSubmit} class="mb-4">
            <input type="file"
                   accept="application/pdf"
                   class="block w-full text-sm disabled:opacity-50 disabled:cursor-not-allowed"
                   disabled={loading}
                   on:change={handlePDFSubmit}
                   bind:this={fileInput}>
            <button type="submit"
                    class="w-full bg-indigo-500 text-white py-1 mt-2 rounded disabled:opacity-50 disabled:cursor-not-allowed"
                    disabled={loading}>
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
                                 bind:squareIndex/>
        {/if}
    </aside>

    <div class="flex-1 grid place-content-center">
        <DebugPdfViewer {data}
                        {mode}
                        {page}
                        {debugColors}
                        {displaySingleText} {textIndex}
                        {displaySingleLine} {lineIndex}
                        {displaySingleSquare} {squareIndex}/>
    </div>
</main>
