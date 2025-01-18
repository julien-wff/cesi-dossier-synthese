<script lang="ts">
    import FileDropDown from '../components/FileDropDown.svelte';
    import { fade } from 'svelte/transition';

    interface Props {
        loading?: boolean;
        onsubmit?: (formData: FormData) => void;
    }

    let { loading = false, onsubmit }: Props = $props();

    let picked = $state(false);
    let filename: string | null = $state(null);
    let fileInput = $state<HTMLInputElement>();

    function handleFileChange(ev: { shiftKey: boolean }) {
        if (!fileInput?.files?.length)
            return;

        filename = fileInput.files[0].name ?? null;

        if (!filename?.toLowerCase().endsWith('.pdf')) {
            fileInput.value = '';
            return;
        }

        picked = true;

        // Submit the form if the user pressed shift while dropping the file
        if (ev.shiftKey)
            handleSubmit();
    }

    function handleCancel() {
        picked = false;
        fileInput && (fileInput.value = '');
    }

    function handleSubmit() {
        if (!fileInput?.files?.[0])
            return;

        const formData = new FormData();
        formData.append('file', fileInput.files[0]);
        onsubmit?.(formData);
    }
</script>


<div class="w-full min-h-screen grid place-items-center">
    <div class="bg-indigo-50 p-8 rounded-md shadow-md w-96 relative">
        {#if loading}
            <div class="absolute inset-0 bg-black bg-opacity-20 rounded-md cursor-progress z-10"
                 transition:fade></div>
        {/if}

        <h2 class="text-xl mb-6">Dossier de synthèse CESI</h2>

        <FileDropDown bind:input={fileInput} hidden={picked} onchange={handleFileChange}/>

        {#if picked}
            <div class="flex items-center justify-start gap-2 mb-6">
                <img src="icons/document.svg" alt="document" class="w-6 h-6">
                <p class="truncate max-w-full" title={filename}>{filename || 'Fichier PDF inconnu'}</p>
            </div>

            <div class="flex gap-4">
                <button class="bg-red-500 hover:bg-red-700 text-white rounded-md px-4 py-2 flex-1"
                        onclick={handleCancel}>
                    Annuler
                </button>

                <button class="bg-indigo-500 hover:bg-indigo-700 text-white rounded-md px-4 py-2 flex-1"
                        onclick={handleSubmit}>
                    Envoyer
                </button>
            </div>
        {/if}
    </div>
</div>
