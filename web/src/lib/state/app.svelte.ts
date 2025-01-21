import type { Section } from '$lib/types/grades';

export enum State {
    Selection = 'selection',
    Loading = 'loading',
    Display = 'display',
    Error = 'error',
}

export const appState = $state({
    state: State.Selection,
    file: null as File | null,
    grades: [] as Section[],
    error: null as string | null,
});
