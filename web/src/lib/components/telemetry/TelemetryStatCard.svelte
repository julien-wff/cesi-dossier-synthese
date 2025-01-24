<script lang="ts">
    interface Props {
        value: number | string;
        label: string;
        description: string;
    }

    let { value, label, description }: Props = $props();
</script>


<div class="flex items-center justify-center flex-col bg-slate-100 dark:bg-slate-700 p-2 sm:p-4 rounded-sm shadow-sm">
    <p class="text-center font-bold">{label}</p>
    <p class="text-center text-4xl md:text-6xl font-bold">{value}</p>
    <p class="text-center text-xs sm:text-sm">{description}</p>
</div>
