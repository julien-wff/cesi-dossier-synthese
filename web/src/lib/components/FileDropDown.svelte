<script lang="ts">
    interface Props {
        input: HTMLInputElement | undefined;
        hidden?: boolean;
        onchange?: ({ shiftKey }: { shiftKey: boolean }) => void;
    }

    let { input = $bindable(), hidden = false, onchange }: Props = $props();

    let drag = $state(false);
    let shiftKey = false;

    function handleDrop(ev: DragEvent) {
        drag = false;
        shiftKey = ev.shiftKey;
    }

    function handleChange() {
        onchange?.({ shiftKey });
        shiftKey = false;
    }
</script>


<div class="border border-slate-400 dark:border-slate-400 rounded relative"
     class:hidden
     class:bg-blue-100={drag}
     class:border-indigo-500={drag}
     class:dark:bg-blue-900={drag}
     class:dark:border-indigo-500={drag}>

    <input type="file"
           ondragenter={() => (drag = true)}
           ondragleave={() => (drag = false)}
           ondrop={handleDrop}
           onchange={handleChange}
           bind:this={input}
           accept="application/pdf"
           class="w-full h-48 opacity-0 cursor-pointer">

    <div class="absolute inset-0 flex flex-col items-center justify-center gap-2 p-4 text-center pointer-events-none">
        <img src="icons/document.svg" alt="Document" class="w-8 h-8">
        <p class="text-sm text-gray-500 dark:text-gray-400">Cliquez ou glissez-déposez le PDF</p>
    </div>
</div>
