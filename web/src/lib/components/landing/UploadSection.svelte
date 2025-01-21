<script lang="ts">
    import FileDropDown from '$lib/components/landing/FileDropDown.svelte';
    import AdvantageCard from '$lib/components/landing/AdvantageCard.svelte';
    import Gauge from 'lucide-svelte/icons/gauge';
    import MousePointerClick from 'lucide-svelte/icons/mouse-pointer-click';
    import ShieldCheck from 'lucide-svelte/icons/shield-check';
</script>


<section class="flex flex-col items-center justify-center min-h-svh p-2 lg:p-4 lg:pb-8">
    <div class="w-full max-w-[64rem] flex-1 flex flex-col lg:flex-row justify-center lg:justify-between items-center gap-6">
        <h1 class="text-4xl lg:text-5xl font-bold text-center lg:text-left">
            Un outil pour visualiser ton <br class="hidden lg:inline"/>
            <span class="text-blue-500 dark:text-blue-400">dossier de synthèse</span>
            CESI
        </h1>

        <FileDropDown/>
    </div>

    <div class="w-full max-w-[64rem] flex flex-col lg:flex-row gap-2 sm:gap-4 lg:gap-8">
        <AdvantageCard title="Pratique"
                       description="Upload to dossier en deux clics, et voilà !"
                       icon={Gauge}/>
        <AdvantageCard title="Intéractif"
                       description="Modifie tes notes et voit les résultats en direct"
                       icon={MousePointerClick}/>
        <AdvantageCard title="Sécurisé"
                       description="Aucun PDF n'est conservé, projet open-source"
                       icon={ShieldCheck}
                       href="https://github.com/julien-wff/cesi-dossier-synthese"/>
    </div>
</section>
