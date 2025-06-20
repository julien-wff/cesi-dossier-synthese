<script lang="ts">
    import { page } from '$app/state';

    function handleBackClick() {
        window.history.back();
    }
</script>


<div class="min-h-svh grid place-items-center">
    <div class="bg-slate-100 dark:bg-slate-700 p-4 rounded-md shadow-md w-96">
        <h1 class="text-xl mb-6">Oh non ! Une erreur est survenue&nbsp;😢</h1>
        <p class="mb-6">{page.state.error}</p>
        <button class="bg-red-500 hover:bg-red-700 w-full text-white px-4 py-2 rounded-md transition-colors cursor-pointer"
                onclick={handleBackClick}>
            Retour
        </button>
    </div>
</div>
