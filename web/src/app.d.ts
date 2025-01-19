// See https://kit.svelte.dev/docs/types#app
// for information about these interfaces
declare global {
    namespace App {
        // interface Error {}
        // interface Locals {}
        // interface PageData {}
        // interface Platform {}
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
