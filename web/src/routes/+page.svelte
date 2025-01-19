<script lang="ts">
    import '$lib/app.css';
    import Meta from '$lib/components/Meta.svelte';
    import Failure from '$lib/views/Failure.svelte';
    import Grades from '$lib/views/Grades.svelte';
    import Home from '$lib/views/Home.svelte';
    import { fade } from 'svelte/transition';
    import type { Section } from '$lib/types/grades';
    import { PUBLIC_API_ENDPOINT } from '$env/static/public';

    enum AppState {
        Selection = 'selection',
        Loading = 'loading',
        Display = 'display',
        Error = 'error',
    }

    let appState = $state<AppState>(AppState.Selection);
    let grades = $state<Section[]>([]);
    let error = $state<string | null>(null);

    async function handlePDFSubmit(form: FormData) {
        appState = AppState.Loading;
        try {
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
            appState = AppState.Error;
            error = (e as Error).message;
        }
    }
</script>


<Meta/>

<main class="min-h-svh">
    {#if appState === AppState.Selection || appState === AppState.Loading}
        <div transition:fade class="absolute inset-0">
            <Home onsubmit={handlePDFSubmit} loading={appState === AppState.Loading}/>
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
