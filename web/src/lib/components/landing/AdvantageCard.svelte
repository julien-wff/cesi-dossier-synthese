<script lang="ts">
    import type { Icon as IconType } from 'lucide-svelte';

    interface Props {
        title: string;
        description: string;
        icon: typeof IconType;
        href?: string;
    }

    let { title, description, icon: Icon, href }: Props = $props();
</script>


{#snippet body()}
    <div class="flex flex-col gap-2">
        <div class="flex gap-2 items-center">
            <Icon class="text-blue-500 dark:text-blue-400"/>
            <p class="text-xl">
                {title}
            </p>
        </div>
        {description}
    </div>
{/snippet}


{#if href}
    <a {href} rel="noreferrer noopener" class="flex-1 p-4 rounded-lg shadow-lg bg-slate-100 dark:bg-slate-700">
        {@render body()}
    </a>
{:else}
    <div class="flex-1 p-4 rounded-lg shadow-lg bg-slate-100 dark:bg-slate-700">
        {@render body()}
    </div>
{/if}
