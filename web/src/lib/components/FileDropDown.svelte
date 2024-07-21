<script lang="ts">
    import { createEventDispatcher } from 'svelte';

    export let input;
    export let hidden = false;

    const dispatch = createEventDispatcher();

    let drag = false;
    let shiftKey = false;

    function handleDrop(ev: DragEvent) {
        drag = false;
        shiftKey = ev.shiftKey;
    }

    function handleChange() {
        dispatch('change', { shiftKey });
        shiftKey = false;
    }
</script>


<div class="border border-gray-400 rounded relative"
     class:hidden
     class:bg-blue-100={drag}
     class:border-indigo-500={drag}>

    <input type="file"
           on:dragenter={() => (drag = true)}
           on:dragleave={() => (drag = false)}
           on:drop={handleDrop}
           on:change={handleChange}
           bind:this={input}
           accept="application/pdf"
           class="w-full h-48 opacity-0 cursor-pointer">

    <div class="absolute inset-0 flex flex-col items-center justify-center gap-2 p-4 text-center pointer-events-none">
        <img src="icons/document.svg" alt="Document" class="w-8 h-8">
        <p class="text-sm text-gray-500">Cliquez ou glissez-déposez le PDF</p>
    </div>
</div>
