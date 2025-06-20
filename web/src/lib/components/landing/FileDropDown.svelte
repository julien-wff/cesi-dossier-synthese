<script lang="ts">
    import FileText from 'lucide-svelte/icons/file-text';
    import CloudAlert from 'lucide-svelte/icons/cloud-alert';
    import { appState, handlePDFSubmit } from '$lib/state/app.svelte.js';
    import { browser } from '$app/environment';

    let inputField = $state<HTMLInputElement>();
    let drag = $state(false);
    let online = $state(browser ? navigator.onLine : true);

    function handleChange() {
        appState.file = inputField?.files?.[0] ?? null
        handlePDFSubmit();
    }
</script>


<svelte:window onoffline={() => (online = false)} ononline={() => (online = true)}/>

<div class="border-4 border-blue-500 dark:border-blue-400 rounded-2xl relative shadow-lg transition-colors"
     class:bg-blue-100={drag}
     class:bg-slate-100={!drag}
     class:border-red-500={!online}
     class:dark:bg-blue-900={drag}
     class:dark:bg-transparent={!drag}
     class:dark:border-red-400={!online}>

    <input accept="application/pdf"
           bind:this={inputField}
           class="w-80 h-48 opacity-0"
           class:cursor-not-allowed={!online}
           class:cursor-pointer={!appState.loading && online}
           class:cursor-progress={appState.loading}
           disabled={appState.loading || !online}
           onchange={handleChange}
           ondragenter={() => (drag = true)}
           ondragleave={() => (drag = false)}
           ondrop={() => (drag = false)}
           type="file">

    <div class="absolute inset-0 flex flex-col items-center justify-center gap-4 p-4 text-center pointer-events-none">
        {#if online}
            <FileText class="w-8 h-8"/>
            <p class="font-bold">Clique ou dépose ton<br>dossier PDF ici</p>
        {:else}
            <CloudAlert class="w-8 h-8"/>
            <p class="font-bold">Hors ligne</p>
        {/if}
    </div>
</div>
