<script lang="ts">
    import '$lib/app.css';
    import Meta from '$lib/components/Meta.svelte';
    import Failure from '$lib/views/Failure.svelte';
    import Grades from '$lib/views/Grades.svelte';
    import Home from '$lib/views/Home.svelte';
    import Loading from '$lib/views/Loading.svelte';
    import { fade } from 'svelte/transition';
    import type { Section } from '$lib/types/grades';

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
            const res = await fetch('/api/parse', {
                method: 'POST',
                body: ev.detail,
            });
            const content = await res.json();

            if (content.error)
                throw new Error(content.error);
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
    {#if state === AppState.Selection}
        <div transition:fade>
            <Home on:submit={handlePDFSubmit}/>
        </div>
    {:else if state === AppState.Loading}
        <div transition:fade>
            <Loading/>
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
