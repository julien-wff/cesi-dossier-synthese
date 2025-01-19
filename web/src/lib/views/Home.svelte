<script lang="ts">
    import FileDropDown from '../components/FileDropDown.svelte';
    import { fade } from 'svelte/transition';

    interface Props {
        loading?: boolean;
        onsubmit?: () => void;
        selectedFile: File | null;
    }

    let { loading = false, onsubmit, selectedFile = $bindable() }: Props = $props();

    let fileInput = $state<HTMLInputElement>();
    let shiftKey = $state(false);

    $effect(() => {
        if (selectedFile && shiftKey)
            onsubmit?.();
    });

    function handleCancel() {
        selectedFile = null;
        fileInput && (fileInput.value = '');
    }
</script>


<div class="w-full min-h-svh grid place-items-center">
    <div class="bg-slate-100 dark:bg-slate-700 p-8 rounded-md shadow-md w-96 relative">
        {#if loading}
            <div class="absolute inset-0 bg-black bg-opacity-20 dark:bg-opacity-40 rounded-md cursor-progress z-10"
                 transition:fade></div>
        {/if}

        <h2 class="text-xl mb-6">Dossier de synthèse CESI</h2>

        <FileDropDown bind:file={selectedFile} hidden={!!selectedFile} bind:shiftKey/>

        {#if selectedFile}
            <div class="flex items-center justify-start gap-2 mb-6">
                <img src="icons/document.svg" alt="document" class="w-6 h-6">
                <p class="truncate max-w-full" title={selectedFile.name}>
                    {selectedFile.name || 'Fichier PDF inconnu'}
                </p>
            </div>

            <div class="flex gap-4">
                <button class="bg-red-500 hover:bg-red-700 text-white rounded-md px-4 py-2 flex-1"
                        onclick={handleCancel}>
                    Annuler
                </button>

                <button class="bg-indigo-500 hover:bg-indigo-700 text-white rounded-md px-4 py-2 flex-1"
                        onclick={onsubmit}>
                    Envoyer
                </button>
            </div>
        {/if}
    </div>
</div>
