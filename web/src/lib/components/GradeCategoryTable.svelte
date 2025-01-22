<script lang="ts">
    import { onMount } from 'svelte';
    import type { Grade, Letter } from '../types/grades';
    import { calculateAverage, gradesWithCoefficientToList, validationLevel } from '../utils/grades';
    import GradeSelector from './GradeSelector.svelte';
    import Undo from 'lucide-svelte/icons/undo-2';

    interface Props {
        category: string;
        grades?: Grade[];
    }

    let { category, grades = $bindable([]) }: Props = $props();

    let categoryAverage = $derived(calculateAverage(gradesWithCoefficientToList(grades)));
    let categoryLevel = $derived(validationLevel(categoryAverage));
    let initialGrades: (Letter | null)[] = $state([]);

    onMount(() => {
        initialGrades = grades.map(grade => grade.letter);
    });

    let gradesChanged = $derived(grades.some((grade, index) => grade.letter !== initialGrades[index]));

    function handleGradesRevert() {
        grades = grades.map((grade, index) => ({
            ...grade,
            letter: initialGrades[index],
        }) satisfies Grade);
    }
</script>


<div class="flex flex-col divide-y divide-gray-400 mt-4 rounded shadow border border-gray-400">

    <div class="flex justify-between">
        <div class="p-1.5">
            {category}
        </div>
        <Undo class="w-9 h-9 cursor-pointer p-1.5 {!gradesChanged ? 'hidden' : ''}"
              aria-hidden={!gradesChanged}
              onclick={handleGradesRevert}/>
    </div>

    <div class="flex divide-x divide-gray-400">
        <div class="flex-1 divide-y divide-gray-400">
            {#each grades as { name, coefficient }, i}
                <div class="flex divide-x divide-gray-400">
                    <div class="p-1.5 flex-1">{name}</div>
                    <div class="grid place-items-center p-1.5 w-9">{coefficient}</div>
                    <GradeSelector bind:grade={grades[i].letter}/>
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
