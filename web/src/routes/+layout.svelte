<script lang="ts">
    import inter400 from '@fontsource/inter/files/inter-latin-400-normal.woff2';
    import inter700 from '@fontsource/inter/files/inter-latin-700-normal.woff2';
    import type { Snippet } from 'svelte';

    interface Props {
        children: Snippet;
    }

    let { children }: Props = $props();
</script>

<svelte:head>
    <link as="font" crossorigin="anonymous" href={inter400} rel="preload" type="font/woff2">
    <link as="font" crossorigin="anonymous" href={inter700} rel="preload" type="font/woff2">
</svelte:head>

{@render children()}
