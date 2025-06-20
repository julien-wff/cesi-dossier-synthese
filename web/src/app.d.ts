// See https://kit.svelte.dev/docs/types#app
// for information about these interfaces
import type { AppView } from '$lib/state/app.svelte';
import type { Section } from '$lib/types/grades';

declare global {
    namespace App {
        // interface Error {}
        // interface Locals {}
        // interface PageData {}
        // interface Platform {}

        interface PageState {
            view?: AppView;
            grades?: Section[];
            error?: string;
        }
    }


    interface Window {
        launchQueue?: LaunchQueue;
    }
}

interface LaunchQueue {
    setConsumer: (consumer: (params: LaunchParams) => void) => void;
}

interface LaunchParams {
    targetURL: string;
    files: FileSystemFileHandle[];
}

export {};
