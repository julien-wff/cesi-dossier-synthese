<script lang="ts">
    interface Props {
        title?: string;
    }

    let { title = 'Dossier de synthèse CESI' }: Props = $props();
</script>


<svelte:head>
    <!-- Primary Meta Tags -->
    <title>{title}</title>
    <meta name="title" content="Dossier de synthèse CESI">
    <meta name="description" content="Outil pour visualiser et tester les notes de son dossier de synthèse">

    <!-- Favicons and PWA -->
    <link rel="icon" type="image/png" href="/favicons/favicon-96x96.png" sizes="96x96" />
    <link rel="icon" type="image/svg+xml" href="/favicons/favicon.svg" />
    <link rel="shortcut icon" href="/favicons/favicon.ico" />
    <link rel="apple-touch-icon" sizes="180x180" href="/favicons/apple-touch-icon.png" />
    <meta name="apple-mobile-web-app-title" content="Dossier CESI" />
    <link rel="manifest" href="/favicons/site.webmanifest" />

    <!-- Open Graph / Facebook -->
    <meta property="og:type" content="website">
    <meta property="og:url" content="https://dossier.bdecesinancy.fr/">
    <meta property="og:title" content="Dossier de synthèse CESI">
    <meta property="og:description" content="Outil pour visualiser et tester les notes de son dossier de synthèse">
    <meta property="og:image" content="/opengraph.png">
    <meta property="og:site_name" content="Dossier de synthèse CESI">

    <!-- Twitter -->
    <meta property="twitter:card" content="summary_large_image">
</svelte:head>
