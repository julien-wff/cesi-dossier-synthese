<script>
    export let input;
    export let hidden = false;

    let drag;
</script>


<div class="border border-gray-400 rounded relative"
     class:hidden
     class:bg-blue-100={drag}
     class:border-indigo-500={drag}>

    <input type="file"
           on:dragenter={() => (drag = true)}
           on:dragleave={() => (drag = false)}
           on:drop={() => (drag = false)}
           bind:this={input}
           on:change
           accept="application/pdf"
           class="w-full h-48 opacity-0 cursor-pointer">

    <div class="absolute inset-0 flex flex-col items-center justify-center gap-2 p-4 text-center pointer-events-none">
        <img src="icons/document.svg" alt="Document" class="w-8 h-8">
        <p class="text-sm text-gray-500">Cliquez ou glissez-déposez le PDF</p>
    </div>
</div>
