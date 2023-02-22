<script>
    import { calculateAverage, gradesWithCoefficientToList, validationLevel } from '../utils/grades.js';
    import GradeSelector from './GradeSelector.svelte';

    export let category;
    export let grades = [];

    let categoryAverage, categoryLevel;
    $: categoryAverage = calculateAverage(gradesWithCoefficientToList(grades));
    $: categoryLevel = validationLevel(categoryAverage);
</script>


<div class="flex flex-col divide-y divide-gray-400 mt-4 rounded-md bg-indigo-50 shadow-md">

    <div class="p-1.5">
        {category}
    </div>

    <div class="flex divide-x divide-gray-400">
        <div class="flex-1 divide-y divide-gray-400">
            {#each grades as { name, coefficient, grade }}
                <div class="flex divide-x divide-gray-400">
                    <div class="p-1.5 flex-1">{name}</div>
                    <div class="grid place-items-center p-1.5 w-9">{coefficient}</div>
                    <GradeSelector bind:grade={grade}/>
                </div>
            {/each}
        </div>

        <div class="grid place-items-center p-1.5 w-10 text-white"
             class:bg-red-500={categoryLevel === 0}
             class:bg-green-500={categoryLevel === 1}
             class:bg-blue-500={categoryLevel === 2}>
            {#if !isNaN(categoryAverage)}
                {Math.round(categoryAverage * 100) / 100}
            {/if}
        </div>
    </div>

</div>
