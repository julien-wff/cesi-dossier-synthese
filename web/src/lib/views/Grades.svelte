<script lang="ts">
    import GradesStatistics from '$lib/components/GradesStatistics.svelte';
    import GradeCategoryTable from '$lib/components/GradeCategoryTable.svelte';
    import { page } from '$app/state';

    let grades = $state(page.state.grades ?? []);

    $effect(() => {
        grades = page.state.grades ?? [];
    });
</script>


<div class="min-h-svh flex flex-col-reverse lg:flex-row items-center lg:items-start justify-center gap-4 md:gap-8 py-4 sm:p-4 md:p-8">
    <div class="sm:max-w-xl w-full p-2 sm:p-4 bg-slate-100 dark:bg-slate-700 sm:shadow-md sm:rounded-md">
        {#each grades as { name: group, categories }, i}
            <h1 class="mb-4 text-xl" class:mt-6={i !== 0}>{group}</h1>
            {#each categories as { name: category }, i}
                <GradeCategoryTable category={category} bind:grades={categories[i].grades}/>
            {/each}
        {/each}
    </div>

    <div class="flex flex-col p-4 rounded-md bg-slate-100 dark:bg-slate-700 shadow-md">
        <GradesStatistics bind:content={grades}/>
    </div>
</div>
