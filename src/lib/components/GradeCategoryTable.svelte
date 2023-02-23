<script>
    import { onMount } from 'svelte';
    import { calculateAverage, gradesWithCoefficientToList, validationLevel } from '../utils/grades.js';
    import GradeSelector from './GradeSelector.svelte';

    export let category;
    export let grades = [];

    let categoryAverage, categoryLevel;
    $: categoryAverage = calculateAverage(gradesWithCoefficientToList(grades));
    $: categoryLevel = validationLevel(categoryAverage);

    let initialGrades = [];
    onMount(() => {
        initialGrades = grades.map(grade => grade.grade);
    });

    let gradesChanged;
    $: gradesChanged = grades.some((grade, index) => grade.grade !== initialGrades[index]);

    function handleGradesRevert() {
        grades = grades.map((grade, index) => ({
            ...grade,
            grade: initialGrades[index],
        }));
    }
</script>


<div class="flex flex-col divide-y divide-gray-400 mt-4 rounded-md bg-indigo-50 shadow-md">

    <div class="flex justify-between">
        <div class="p-1.5">
            {category}
        </div>
        <img src="icons/refresh.svg" alt="*" class="w-9 h-9 cursor-pointer p-1.5"
             class:hidden={!gradesChanged}
             on:click={handleGradesRevert}/>
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
