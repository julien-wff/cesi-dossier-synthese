<script lang="ts">
    import '$lib/app.css';
    import Meta from '$lib/components/Meta.svelte';
    import Failure from '$lib/views/Failure.svelte';
    import Grades from '$lib/views/Grades.svelte';
    import Home from '$lib/views/Home.svelte';
    import { fade } from 'svelte/transition';
    import type { Section } from '$lib/types/grades';
    import { PUBLIC_API_ENDPOINT } from '$env/static/public';
    import { onMount } from 'svelte';

    enum AppState {
        Selection = 'selection',
        Loading = 'loading',
        Display = 'display',
        Error = 'error',
    }

    let appState = $state<AppState>(AppState.Selection);
    let selectedFile = $state<File | null>(null);
    let grades = $state<Section[]>([]);
    let error = $state<string | null>(null);

    onMount(() => {
        // Handle file_handler from PWA webmanifest
        if ('launchQueue' in window && window.launchQueue) {
            window.launchQueue.setConsumer(async (launchParams) => {
                for (const fileHandler of launchParams.files) {
                    selectedFile = await fileHandler.getFile();
                    await handlePDFSubmit();
                }
            });
        }

        // Check if the URL contains a parsing result
        const urlParams = new URLSearchParams(window.location.search);
        if (appState === AppState.Selection && urlParams.has('result')) {
            // Get and apply the result
            grades = JSON.parse(urlParams.get('result')!).data;
            appState = AppState.Display;

            // Clear the URL
            window.history.replaceState({}, document.title, window.location.pathname);
        }
    });

    async function handlePDFSubmit() {
        if (!selectedFile)
            return;

        appState = AppState.Loading;
        try {
            const form = new FormData();
            form.append('file', selectedFile);

            const res = await fetch(PUBLIC_API_ENDPOINT + '/parse', {
                method: 'POST',
                body: form,
            });
            const content = await res.json();

            if (!res.ok)
                throw new Error(content?.message?.fr || 'Une erreur inconnue est survenue');
            if (!('data' in content))
                throw new Error('Aucune donnée n\'a été trouvée dans le PDF');

            grades = content.data;
            appState = AppState.Display;
        } catch (e) {
            console.error(e);
            selectedFile = null;
            error = (e as Error).message;
            appState = AppState.Error;
        }
    }
</script>


<Meta/>

<main class="min-h-svh">
    {#if appState === AppState.Selection || appState === AppState.Loading}
        <div transition:fade class="absolute inset-0">
            <Home onsubmit={handlePDFSubmit} loading={appState === AppState.Loading} bind:selectedFile/>
        </div>
    {:else if appState === AppState.Display}
        <div transition:fade>
            <Grades bind:content={grades}/>
        </div>
    {:else if appState === AppState.Error && error}
        <div transition:fade class="absolute inset-0">
            <Failure {error} onback={() => (appState = AppState.Selection)}/>
        </div>
    {/if}
</main>
