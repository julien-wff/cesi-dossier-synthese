<script lang="ts">
    import { page } from '$app/state';

    const email = 'julien.wolff1@viacesi.fr'
    const emailSubject = encodeURIComponent('Erreur de lecture de dossier de synthèse');

    function handleBackClick() {
        window.history.back();
    }
</script>


<div class="min-h-svh p-2 grid place-items-center">
    <div class="bg-slate-100 dark:bg-slate-700 p-4 rounded-md shadow-md max-w-md w-full">
        <h1 class="text-xl font-bold mb-4">Oh non ! Une erreur est survenue&nbsp;😢</h1>

        <p class="mb-2">{page.state.error}</p>

        <p class="mb-6">
            S'il s'agit d'un dossier valide, vous pouvez me l'envoyer par email à
            <a class="text-blue-500 dark:text-blue-400 underline"
               href="mailto:{email}?subject={emailSubject}">
                {email}
            </a>
            afin que je puisse corriger le problème. Merci d'avance !
        </p>

        <button class="bg-red-500 dark:bg-red-400 hover:bg-red-700 dark:hover:bg-red-600 w-full text-white px-4 py-2 rounded-md transition-colors cursor-pointer"
                onclick={handleBackClick}>
            Retour
        </button>
    </div>
</div>
