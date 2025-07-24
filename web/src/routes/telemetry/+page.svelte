<script lang="ts">
    import '$lib/app.css';
    import Meta from '$lib/components/Meta.svelte';
    import Login from '$lib/views/telemetry/Login.svelte';
    import Telemetry from '$lib/views/telemetry/Telemetry.svelte';
    import { browser } from '$app/environment';
    import { onMount } from 'svelte';
    import { telemetryState } from '$lib/state/telemetry.svelte';

    const telemetryAuthKey = 'telemetry-auth';

    let loaded = $state(false);
    let authToken = $state(browser ? localStorage.getItem(telemetryAuthKey) : null);
    let error = $state<string | null>('');

    function handleLogin(user: string, password: string) {
        authToken = btoa(`${user}:${password}`);
        localStorage.setItem(telemetryAuthKey, authToken);
        fetchTelemetryData();
    }

    async function fetchTelemetryData() {
        error = null;

        const res = await fetch('/api/telemetry', {
            headers: {
                Authorization: `Basic ${authToken}`,
            },
        });

        if (res.status === 401) {
            authToken = null;
            localStorage.removeItem(telemetryAuthKey);
            error = 'Invalid credentials';
            return;
        }

        try {
            const resJson = await res.json();
            if (resJson.error) {
                error = resJson.error;
                return;
            }
            telemetryState.stats = resJson.data;
            telemetryState.loaded = true;
        } catch (e) {
            error = 'Failed to fetch telemetry data';
            return;
        }
    }

    onMount(() => {
        loaded = true;
        if (authToken)
            fetchTelemetryData();
    });
</script>

<Meta title="Télémétrie - Dossier de synthèse"/>

<main class="min-h-svh">
    {#if authToken && telemetryState.loaded && !error}
        <Telemetry/>
    {:else if !error && (!loaded || (loaded && authToken))}
        <div class="min-h-svh grid place-content-center">
            <div class="w-8 h-8 border-t-2 border-slate-900 dark:border-slate-50 rounded-full animate-spin"></div>
        </div>
    {:else}
        <Login onlogin={handleLogin} {error}/>
    {/if}
</main>
