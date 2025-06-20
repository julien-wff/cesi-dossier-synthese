import { pushState } from '$app/navigation';

export const AppView = {
    SELECTION: undefined,
    DISPLAY: 'display',
    ERROR: 'error',
} as const;

export const appState = $state({
    file: null as File | null,
    loading: false,
});


export async function handlePDFSubmit() {
    if (!appState.file || appState.loading)
        return;

    appState.loading = true;

    try {
        const form = new FormData();
        form.append('file', appState.file);

        const res = await fetch('/api/parse', {
            method: 'POST',
            body: form,
        });

        if (res.status == 429)
            throw new Error('Trop de requêtes, réessaye d\'ici quelques minutes');

        const content = await res.json();

        if (!res.ok)
            throw new Error(content?.message?.fr || 'Une erreur inconnue est survenue');
        if (!('data' in content))
            throw new Error('Aucune donnée n\'a été trouvée dans le PDF');

        pushState('', { view: AppView.DISPLAY, grades: content.data });
    } catch (e) {
        console.error(e);
        pushState('', {
            view: AppView.ERROR,
            error: (e as Error).message,
        });
    } finally {
        appState.loading = false;
        appState.file = null;
    }
}
