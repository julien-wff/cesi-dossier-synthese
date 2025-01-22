<script lang="ts">
    import '$lib/app.css';
    import Meta from '$lib/components/Meta.svelte';
    import Failure from '$lib/views/Failure.svelte';
    import Grades from '$lib/views/Grades.svelte';
    import Home from '$lib/views/Home.svelte';
    import { fade } from 'svelte/transition';
    import { PUBLIC_API_ENDPOINT } from '$env/static/public';
    import { onMount } from 'svelte';
    import { appState, State } from '$lib/state/app.svelte.js';

    onMount(() => {
        // Handle file_handler from PWA webmanifest
        if ('launchQueue' in window && window.launchQueue) {
            window.launchQueue.setConsumer(async (launchParams) => {
                for (const fileHandler of launchParams.files) {
                    appState.file = await fileHandler.getFile();
                    await handlePDFSubmit();
                }
            });
        }

        // Check if the URL contains a parsing result
        const urlParams = new URLSearchParams(window.location.search);
        if (appState.state === State.Selection && urlParams.has('result')) {
            // Get and apply the result
            appState.grades = JSON.parse(urlParams.get('result')!).data;
            appState.state = State.Display;

            // Clear the URL
            window.history.replaceState({}, document.title, window.location.pathname);
        }

        // Check if the URL contains an error
        if (appState.state === State.Selection && urlParams.has('error')) {
            // Get and apply the error
            appState.error = JSON.parse(urlParams.get('error')!).message.fr;
            appState.state = State.Error;

            // Clear the URL
            window.history.replaceState({}, document.title, window.location.pathname);
        }
    });

    $effect(() => {
        if (appState.file && !appState.grades.length)
            handlePDFSubmit();
    });

    async function handlePDFSubmit() {
        if (!appState.file)
            return;

        appState.state = State.Loading;
        try {
            const form = new FormData();
            form.append('file', appState.file);

            const res = await fetch(PUBLIC_API_ENDPOINT + '/parse', {
                method: 'POST',
                body: form,
            });
            const content = await res.json();

            if (!res.ok)
                throw new Error(content?.message?.fr || 'Une erreur inconnue est survenue');
            if (!('data' in content))
                throw new Error('Aucune donnée n\'a été trouvée dans le PDF');

            appState.grades = content.data;
            appState.state = State.Display;
        } catch (e) {
            console.error(e);
            appState.file = null;
            appState.error = (e as Error).message;
            appState.state = State.Error;
        }
    }
</script>


<Meta/>

<main class="min-h-svh" class:cursor-progress={appState.state === State.Loading}>
    {#if appState.state === State.Selection || appState.state === State.Loading}
        <div transition:fade class="absolute inset-0">
            <Home/>
        </div>
    {:else if appState.state === State.Display}
        <div transition:fade>
            <Grades/>
        </div>
    {:else if appState.state === State.Error}
        <div transition:fade class="absolute inset-0">
            <Failure/>
        </div>
    {/if}
</main>
