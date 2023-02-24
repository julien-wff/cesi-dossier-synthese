<script>
    import '$lib/app.css';
    import Meta from '$lib/components/Meta.svelte';
    import Failure from '$lib/views/Failure.svelte';
    import Grades from '$lib/views/Grades.svelte';
    import Home from '$lib/views/Home.svelte';
    import Loading from '$lib/views/Loading.svelte';
    import { fade } from 'svelte/transition';

    let state = 'selection';
    let grades;
    let error;

    async function handlePDFSubmit(ev) {
        state = 'loading';
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
            state = 'display';
        } catch (e) {
            console.error(e);
            state = 'error';
            error = e.message;
        }
    }
</script>


<Meta/>

<main>
    {#if state === 'selection'}
        <div transition:fade>
            <Home on:submit={handlePDFSubmit}/>
        </div>
    {:else if state === 'loading'}
        <div transition:fade>
            <Loading/>
        </div>
    {:else if state === 'display'}
        <div transition:fade>
            <Grades content={grades}/>
        </div>
    {:else if state === 'error'}
        <div transition:fade>
            <Failure {error} on:back={() => (state = 'selection')}/>
        </div>
    {/if}
</main>
