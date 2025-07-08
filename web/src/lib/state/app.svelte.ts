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


function handleError(message: string, error?: any) {
    if (error) {
        console.error(new Error(message), 'Source error:', error);
    } else {
        console.error(new Error(message));
    }
    pushState('', {
        view: AppView.ERROR,
        error: message,
    });
}

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

        if (res.status == 429) {
            return handleError('Trop de requêtes, veuillez réessayez d\'ici quelques minutes.');
        }

        let content = null;
        try {
            content = await res.json();
        } catch {
            // Ignore body
        }

        if (!res.ok || !content)
            return handleError(content?.message?.fr ?? 'Une erreur inconnue est survenue. Veuillez réessayer plus tard.');

        if (!Array.isArray(content.data) || content.data.length === 0)
            return handleError('Aucune note n\'a été trouvée dans le PDF.');

        pushState('', { view: AppView.DISPLAY, grades: content.data });
    } catch (e) {
        handleError('Une erreur inconnue est survenue lors de l\'envoi du PDF. Veuillez réessayer plus tard.', e);
    } finally {
        appState.loading = false;
        appState.file = null;
    }
}
