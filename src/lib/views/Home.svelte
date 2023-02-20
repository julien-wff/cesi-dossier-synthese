<script>
    import { createEventDispatcher } from 'svelte';
    import FileDropDown from '../components/FileDropDown.svelte';

    let picked = false;
    let filename = null;
    let fileInput;

    const dispatch = createEventDispatcher();

    function handleFileChange() {
        filename = fileInput.files[0].name;

        if (!filename.toLowerCase().endsWith('.pdf')) {
            fileInput.value = null;
            return;
        }

        picked = true;
    }

    function handleCancel() {
        picked = false;
        fileInput.value = null;
    }

    function handleSubmit() {
        const formData = new FormData();
        formData.append('file', fileInput.files[0]);
        dispatch('submit', formData);
    }
</script>


<div class="w-full min-h-screen grid place-items-center">
    <div class="bg-indigo-50 p-8 rounded-md shadow-md w-96">
        <h2 class="text-xl mb-6">Dossier de synthèse CESI</h2>

        <FileDropDown bind:input={fileInput} hidden={picked} on:change={handleFileChange}/>

        {#if picked}
            <div class="flex items-center justify-start gap-2 mb-6">
                <img src="icons/document.svg" alt="document" class="w-6 h-6">
                <p class="truncate max-w-full" title={filename}>{filename || 'Fichier PDF inconnu'}</p>
            </div>

            <div class="flex gap-4">
                <button class="bg-red-500 hover:bg-red-700 text-white rounded-md px-4 py-2 flex-1"
                        on:click={handleCancel}>
                    Annuler
                </button>

                <button class="bg-indigo-500 hover:bg-indigo-700 text-white rounded-md px-4 py-2 flex-1"
                        on:click={handleSubmit}>
                    Envoyer
                </button>
            </div>
        {/if}
    </div>
</div>
