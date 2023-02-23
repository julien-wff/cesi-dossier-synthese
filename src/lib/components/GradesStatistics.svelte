<script lang="ts">
    import type { Section } from '../types/grades';
    import { calculateCategoriesAverage, calculateGradesAverage, validationLevel } from '../utils/grades';

    export let content: Section[] = [];

    let gradesAverage: number, categoriesAverage: number, gradesLevel: number | null, categoryLevel: number | null;
    $: gradesAverage = calculateGradesAverage(content);
    $: categoriesAverage = calculateCategoriesAverage(content);
    $: gradesLevel = validationLevel(gradesAverage);
    $: categoryLevel = validationLevel(categoriesAverage);
</script>


<h1 class="text-xl mb-4">Moyennes</h1>

<div class="grid grid-cols-[1fr_auto] border border-gray-400 shadow rounded">
    <p class="p-1.5 border-gray-400 border-b border-r">
        Moyenne des notes
    </p>
    <p class="p-1.5 border-gray-400 border-b w-11 text-center text-white"
       class:bg-red-500={gradesLevel === 0}
       class:bg-green-500={gradesLevel === 1}
       class:bg-blue-500={gradesLevel === 2}>
        {isNaN(gradesAverage) ? '' : Math.round(gradesAverage * 100) / 100}
    </p>
    <p class="p-1.5 border-gray-400 border-r">
        Moyenne des U.E.
    </p>
    <p class="p-1.5 w-11 text-center text-white"
       class:bg-red-500={categoryLevel === 0}
       class:bg-green-500={categoryLevel === 1}
       class:bg-blue-500={categoryLevel === 2}>
        {isNaN(categoriesAverage) ? '' : Math.round(categoriesAverage * 100) / 100}
    </p>
</div>
