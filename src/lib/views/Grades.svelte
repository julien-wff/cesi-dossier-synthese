<script>
    import GradeCategoryTable from '../components/GradeCategoryTable.svelte';

    export let content;
</script>


<div class="flex flex-col items-center">
    <div class="max-w-xl">
        {#each content as { name: group, categories }}
            <h1 class="mt-6 mb-4 text-xl">{group}</h1>
            {#each categories as { name: category, grades }}
                <GradeCategoryTable category={category} bind:grades={grades}/>
            {/each}
        {/each}
    </div>
</div>
