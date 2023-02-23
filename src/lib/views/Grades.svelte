<script lang="ts">
    import GradesStatistics from '$lib/components/GradesStatistics.svelte';
    import GradeCategoryTable from '../components/GradeCategoryTable.svelte';
    import type { Section } from '../types/grades';

    export let content: Section[];
</script>


<div class="flex flex-col-reverse lg:flex-row items-center lg:items-start justify-center gap-4 md:gap-8 py-4 sm:p-4 md:p-8">
    <div class="sm:max-w-xl w-full p-2 sm:p-4 bg-indigo-50 sm:shadow-md sm:rounded-md">
        {#each content as { name: group, categories }, i}
            <h1 class="mb-4 text-xl" class:mt-6={i !== 0}>{group}</h1>
            {#each categories as { name: category, grades }}
                <GradeCategoryTable category={category} bind:grades={grades}/>
            {/each}
        {/each}
    </div>

    <div class="flex flex-col p-4 rounded-md bg-indigo-50 shadow-md">
        <GradesStatistics bind:content={content}/>
    </div>
</div>
