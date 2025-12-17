import { pushState } from '$app/navigation';

export const AppView = {
    SELECTION: undefined,
    DISPLAY: 'display',
    ERROR: 'error',
} as const;

export const ErrorType = {
    NETWORK: 'network',
    INVALID_FILE: 'invalid_file',
} as const;

type TErrorType = typeof ErrorType[keyof typeof ErrorType];

export const appState = $state({
    file: null as File | null,
    loading: false,
    errorType: ErrorType.NETWORK as TErrorType,
});


function handleError(type: TErrorType, message: string, error?: any) {
    appState.loading = false;
    appState.file = null;
    appState.errorType = type;

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

    let res: Response;
    try {
        const form = new FormData();
        form.append('file', appState.file);

        res = await fetch('/api/parse', {
            method: 'POST',
            body: form,
        });
    } catch (e) {
        return handleError(ErrorType.NETWORK, 'Impossible de contacter le serveur. Essayez de changer de connexion (le réseau CESI peut bloquer le service).', e);
    }

    if (res.status === 429) {
        return handleError(ErrorType.NETWORK, 'Trop de requêtes, veuillez réessayer d\'ici quelques minutes.', res);
    }

    let content = null;
    try {
        content = await res.json();
    } catch {
        // Ignore body
    }

    if (!res.ok || !content) {
        return handleError(ErrorType.INVALID_FILE, content?.message?.fr ?? 'Une erreur inconnue est survenue. Veuillez réessayer plus tard.', res);
    }

    if (!Array.isArray(content.data) || content.data.length === 0) {
        return handleError(ErrorType.INVALID_FILE, 'Aucune note n\'a été trouvée dans le PDF.', content);
    }

    appState.loading = false;
    appState.file = null;
    pushState('', { view: AppView.DISPLAY, grades: content.data });
}
