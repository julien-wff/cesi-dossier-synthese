<script lang="ts">
    import Loading from '$lib/views/Loading.svelte';

    let email: string;
    let password: string;

    let eventSource: EventSource;
    let loading = false;
    let loadingMessage = 'Connexion en cours...';
    let userID: string;

    function handleSubmit() {
        loading = true;
        eventSource = new EventSource(`/api/download`);
        eventSource.onmessage = async (event) => {
            const message = JSON.parse(event.data);
            switch (message.event) {
                case 'connected':
                    loadingMessage = 'Transmission des données...';
                    userID = message.data.id;
                    sendCredentials();
                    break;
                case 'loading_browser':
                    loadingMessage = 'Chargement du navigateur...';
                    break;
                case 'loading_email_page':
                    loadingMessage = 'Chargement de la page d\'email...';
                    break;
                case 'loading_login_page':
                    loadingMessage = 'Chargement de la page de connexion...';
                    break;
                case 'error':
                    loading = false;
                    eventSource.close();
                    console.error(message.data);
                    break;
            }
        };
        eventSource.onerror = (event) => {
            console.error(event);
            loading = false;
            eventSource.close();
        };
    }

    const sendCredentials = () => fetch(`/api/download`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            id: userID,
            email,
            password,
        }),
    });
</script>


{#if !loading}
    <div class="w-full min-h-screen grid place-items-center">
        <form class="bg-indigo-50 p-8 rounded-md shadow-md w-96" on:submit|preventDefault={handleSubmit}>
            <h2 class="text-xl mb-6">Connexion au compte CESI</h2>
            <label for="email-input" class="block mb-2">Adresse email CESI</label>
            <input type="text"
                   required
                   bind:value={email}
                   pattern=".+@viacesi.fr$"
                   placeholder="prénom.nom@viacesi.fr"
                   class="w-full border border-gray-300 focus-visible:outline-none focus:border-indigo-500 bg-white p-2 rounded-md mb-4"
                   id="email-input">
            <label for="password-input" class="block mb-2">Mot de passe</label>
            <input type="password"
                   required
                   bind:value={password}
                   class="w-full border border-gray-300 focus-visible:outline-none focus:border-indigo-500 bg-white p-2 rounded-md mb-6"
                   id="password-input">
            <button type="submit"
                    class="w-full bg-indigo-500 hover:bg-indigo-600 text-white p-2 rounded-md">
                Se connecter
            </button>
        </form>
    </div>
{:else}
    <Loading message={loadingMessage}/>
{/if}
