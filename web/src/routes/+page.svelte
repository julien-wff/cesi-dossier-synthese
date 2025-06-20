<script lang="ts">
    import '$lib/app.css';
    import Meta from '$lib/components/Meta.svelte';
    import Failure from '$lib/views/Failure.svelte';
    import Grades from '$lib/views/Grades.svelte';
    import Home from '$lib/views/Home.svelte';
    import { fade } from 'svelte/transition';
    import { onMount } from 'svelte';
    import { appState, AppView, handlePDFSubmit } from '$lib/state/app.svelte.js';
    import { page } from '$app/state';
    import { afterNavigate, goto, pushState } from '$app/navigation';

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
    });

    afterNavigate(async (nav) => {
        // Wait for the svelte router to be ready
        await nav.complete;

        // Check if the URL contains a parsing result
        if ((!page.state.view || page.state.view === AppView.SELECTION) && page.url.searchParams.has('result')) {
            const result = JSON.parse(page.url.searchParams.get('result')!);
            await goto('.'); // Clear the URL params
            pushState('', {
                view: AppView.DISPLAY,
                grades: result.data,
            });
        }

        // Check if the URL contains an error
        if ((!page.state.view || page.state.view === AppView.SELECTION) && page.url.searchParams.has('error')) {
            await goto('.'); // Clear the URL params
            pushState('', {
                view: AppView.ERROR,
                error: JSON.parse(page.url.searchParams.get('error')!).message.fr,
            });
        }
    });
</script>


<Meta/>


<main class="min-h-svh" class:cursor-progress={appState.loading}>
    {#if page.state.view === AppView.SELECTION}
        <div transition:fade class="absolute inset-0">
            <Home/>
        </div>
    {:else if page.state.view === AppView.DISPLAY}
        <div transition:fade>
            <Grades/>
        </div>
    {:else if page.state.view === AppView.ERROR}
        <div transition:fade class="absolute inset-0">
            <Failure/>
        </div>
    {/if}
</main>
