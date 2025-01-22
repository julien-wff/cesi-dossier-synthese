<script lang="ts">
    import FileText from 'lucide-svelte/icons/file-text';
    import CloudAlert from 'lucide-svelte/icons/cloud-alert';
    import { appState, State } from '$lib/state/app.svelte.js';
    import { browser } from '$app/environment';

    let inputField = $state<HTMLInputElement>();
    let drag = $state(false);
    let online = $state(browser ? navigator.onLine : true);

    const getFieldFile = () => inputField?.files?.[0] ?? null;

    $effect(() => {
        if (!inputField)
            return;

        if (appState.file && !getFieldFile()) {
            const dataTransfer = new DataTransfer();
            dataTransfer.items.add(appState.file);
            inputField.files = dataTransfer.files;
        } else if (!appState.file && getFieldFile()) {
            inputField.value = '';
        }
    });

    function handleChange() {
        appState.file = getFieldFile();
    }
</script>


<svelte:window ononline={() => (online = true)} onoffline={() => (online = false)}/>

<div class="border-4 border-blue-500 dark:border-blue-400 rounded-2xl relative shadow-lg transition-colors"
     class:bg-slate-100={!drag}
     class:dark:bg-transparent={!drag}
     class:bg-blue-100={drag}
     class:dark:bg-blue-900={drag}
     class:border-red-500={!online}
     class:dark:border-red-400={!online}>

    <input type="file"
           ondragenter={() => (drag = true)}
           ondragleave={() => (drag = false)}
           ondrop={() => (drag = false)}
           onchange={handleChange}
           bind:this={inputField}
           disabled={appState.state === State.Loading || !online}
           class:cursor-pointer={appState.state !== State.Loading && online}
           class:cursor-progress={appState.state === State.Loading}
           class:cursor-not-allowed={!online}
           accept="application/pdf"
           class="w-80 h-48 opacity-0">

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
