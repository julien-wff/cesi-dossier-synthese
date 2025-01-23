<script lang="ts">
    import type { Section } from '../types/grades';
    import {
        calculateCategoriesAverage,
        calculateGradesAverage,
        countLetterWithCoeff,
        countLetterWithoutCoeff,
        validationLevel,
    } from '$lib/utils/grades';

    interface Props {
        content?: Section[];
    }

    let { content = $bindable([]) }: Props = $props();

    let gradesAverage = $derived(calculateGradesAverage(content));
    let categoriesAverage = $derived(calculateCategoriesAverage(content));
    let gradesLevel = $derived(validationLevel(gradesAverage));
    let categoryLevel = $derived(validationLevel(categoriesAverage));
</script>


<h1 class="text-xl mb-4">Moyennes</h1>

<div class="grid grid-cols-[1fr_auto] border border-gray-400 shadow-sm rounded-sm">
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

<h1 class="text-xl mt-6 mb-4">Lettres</h1>

<div class="flex border border-gray-400 w-fit shadow-sm rounded-sm divide-x divide-gray-400">
    <div class="flex flex-col divide-y divide-gray-400">
        <div class="p-1.5 text-center">Lettre</div>
        <div class="p-1.5 text-center text-white bg-blue-500">A</div>
        <div class="p-1.5 text-center text-white bg-green-500">B</div>
        <div class="p-1.5 text-center text-white bg-amber-500">C</div>
        <div class="p-1.5 text-center text-white bg-red-500">D</div>
    </div>
    <div class="flex flex-col divide-y divide-gray-400">
        <div class="p-1.5 truncate">Sans Coeff.</div>
        <div class="p-1.5 text-center">{countLetterWithoutCoeff(content, 'A')}</div>
        <div class="p-1.5 text-center">{countLetterWithoutCoeff(content, 'B')}</div>
        <div class="p-1.5 text-center">{countLetterWithoutCoeff(content, 'C')}</div>
        <div class="p-1.5 text-center">{countLetterWithoutCoeff(content, 'D')}</div>
    </div>
    <div class="flex flex-col divide-y divide-gray-400">
        <div class="p-1.5 truncate">Avec Coeff.</div>
        <div class="p-1.5 text-center">{countLetterWithCoeff(content, 'A')}</div>
        <div class="p-1.5 text-center">{countLetterWithCoeff(content, 'B')}</div>
        <div class="p-1.5 text-center">{countLetterWithCoeff(content, 'C')}</div>
        <div class="p-1.5 text-center">{countLetterWithCoeff(content, 'D')}</div>
    </div>
</div>
