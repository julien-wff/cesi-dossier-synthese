<script lang="ts">
    import { onMount } from 'svelte';
    import type { Grade, Letter } from '../types/grades';
    import { calculateAverage, gradesWithCoefficientToList, validationLevel } from '../utils/grades';
    import GradeSelector from './GradeSelector.svelte';

    export let category: string;
    export let grades: Grade[] = [];

    let categoryAverage: number, categoryLevel: number | null;
    $: categoryAverage = calculateAverage(gradesWithCoefficientToList(grades));
    $: categoryLevel = validationLevel(categoryAverage);

    let initialGrades: (Letter | null)[] = [];
    onMount(() => {
        initialGrades = grades.map(grade => grade.letter);
    });

    let gradesChanged: boolean;
    $: gradesChanged = grades.some((grade, index) => grade.letter !== initialGrades[index]);

    function handleGradesRevert() {
        grades = grades.map((grade, index) => ({
            ...grade,
            letter: initialGrades[index],
        }) satisfies Grade);
    }
</script>


<div class="flex flex-col divide-y divide-gray-400 mt-4 rounded bg-indigo-50 shadow border border-gray-400">

    <div class="flex justify-between">
        <div class="p-1.5">
            {category}
        </div>
        <img src="icons/refresh.svg"
             alt="*"
             class="w-9 h-9 cursor-pointer p-1.5"
             aria-hidden={!gradesChanged}
             class:hidden={!gradesChanged}
             on:click={handleGradesRevert}/>
    </div>

    <div class="flex divide-x divide-gray-400">
        <div class="flex-1 divide-y divide-gray-400">
            {#each grades as { name, coefficient, letter }}
                <div class="flex divide-x divide-gray-400">
                    <div class="p-1.5 flex-1">{name}</div>
                    <div class="grid place-items-center p-1.5 w-9">{coefficient}</div>
                    <GradeSelector bind:grade={letter}/>
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
