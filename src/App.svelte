<script>
    import Home from './lib/views/Home.svelte';

    function handlePDFSubmit(ev) {
        console.log(ev.detail);
    }
</script>

<main>
    <Home on:submit={handlePDFSubmit}/>
</main>
