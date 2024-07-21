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

    let state: AppState = AppState.Selection;
    let grades: Section[];
    let error: string;

    async function handlePDFSubmit(ev: CustomEvent<FormData>) {
        state = AppState.Loading;
        try {
            const res = await fetch(PUBLIC_API_ENDPOINT + '/parse', {
                method: 'POST',
                body: ev.detail,
            });
            const content = await res.json();

            if (!res.ok)
                throw new Error(content?.message?.fr || 'Une erreur inconnue est survenue');
            if (!('data' in content))
                throw new Error('Aucune donnée n\'a été trouvée dans le PDF');

            grades = content.data;
            state = AppState.Display;
        } catch (e) {
            console.error(e);
            state = AppState.Error;
            error = (e as Error).message;
        }
    }
</script>


<Meta/>

<main>
    {#if state === AppState.Selection || state === AppState.Loading}
        <div transition:fade>
            <Home on:submit={handlePDFSubmit} loading={state === AppState.Loading}/>
        </div>
    {:else if state === AppState.Display}
        <div transition:fade>
            <Grades content={grades}/>
        </div>
    {:else if state === AppState.Error}
        <div transition:fade>
            <Failure {error} on:back={() => (state = AppState.Selection)}/>
        </div>
    {/if}
</main>
