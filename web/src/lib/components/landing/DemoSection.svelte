<script lang="ts">
    import type { Section } from '$lib/types/grades';
    import GradeCategoryTable from '$lib/components/GradeCategoryTable.svelte';
    import GradesStatistics from '$lib/components/GradesStatistics.svelte';

    let grades = $state<Section[]>([ {
        name: 'Sciences et Méthodes de l\'Ingénieur',
        categories: [ {
            name: '[S5] 5.4 Génie logiciel',
            grades: [
                {
                    letter: 'A',
                    coefficient: 2,
                    name: 'Programmation système : Modélisation (Livrable)',
                },
                {
                    letter: 'C',
                    coefficient: 2,
                    name: 'Programmation système : Architectures logicielles et Modélisation UML (CCTL)',
                },
                {
                    letter: 'B',
                    coefficient: 1,
                    name: 'Programmation système : Environnement .net et cycle de développement (CCTL)',
                },
                {
                    letter: 'A',
                    coefficient: 1,
                    name: 'Programmation système : Application non synchronisée (Livrable)',
                },
                {
                    letter: 'A',
                    coefficient: 2,
                    name: 'Programmation système : Application finale (Livrable)',
                },
            ],
        } ],
    } ]);
</script>


<section class="flex flex-col items-center justify-center min-h-svh p-2">
    <div class="w-full max-w-[64rem] flex-1 flex flex-col justify-center items-center gap-8 lg:gap-12">
        <div>
            <h2 class="text-3xl lg:text-4xl font-bold mb-2">
                Fini les PDFs moches et les fichiers Excel
            </h2>

            <p class="italic text-lg">Clique sur une lettre pour la changer !</p>
        </div>

        <div class="flex flex-col lg:flex-row items-center lg:items-start justify-center gap-4 md:gap-8">
            <div class="sm:max-w-xl w-full p-2 sm:p-4 bg-slate-100 dark:bg-slate-700 shadow-md rounded-md">
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
    </div>
</section>
