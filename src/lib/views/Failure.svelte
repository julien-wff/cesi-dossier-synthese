<script lang="ts">
    import { createEventDispatcher } from 'svelte';

    export let error: string;

    const dispatch = createEventDispatcher();

    const handleBack = () => {
        dispatch('back');
    };
</script>


<div class="h-screen grid place-items-center">
    <div class="bg-indigo-50 p-8 rounded-md shadow-md w-96">
        <h1 class="text-xl mb-6">Oh non ! Une erreur est survenue :(</h1>
        <p class="mb-6">{error}</p>
        <button class="bg-red-500 hover:bg-red-700 w-full text-white px-4 py-2 rounded-md"
                on:click={handleBack}>
            Recommencer
        </button>
    </div>
</div>
