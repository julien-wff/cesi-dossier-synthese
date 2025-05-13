<script lang="ts">
    interface Props {
        onlogin: (user: string, password: string) => void;
        error: string | null;
    }

    let { onlogin, error }: Props = $props();

    let user = $state('');
    let password = $state('');

    function handleSubmit(ev: SubmitEvent) {
        ev.preventDefault();
        onlogin(user, password);
    }
</script>


<div class="min-h-svh grid place-content-center">
    <form class="w-full max-w-64 bg-slate-100 dark:bg-slate-700 p-4 rounded shadow" onsubmit={handleSubmit}>
        <h2 class="bold text-xl">Login Required</h2>

        {#if error}
            <p class="text-red-500 mt-2">{error}</p>
        {/if}

        <label for="user" class="block mt-2 mb-1">Username</label>
        <input type="text"
               class="w-full bg-slate-300 dark:bg-slate-500 rounded-sm px-2 py-1"
               id="user"
               required
               bind:value={user}/>

        <label for="password" class="block mt-2 mb-1">Password</label>
        <input type="password"
               class="w-full bg-slate-300 dark:bg-slate-500 rounded-sm px-2 py-1"
               id="password"
               required
               bind:value={password}/>

        <button type="submit"
                class="w-full mt-4 bg-blue-500 dark:bg-blue-400 text-slate-50 rounded-sm px-2 py-1 cursor-pointer">
            Login
        </button>
    </form>
</div>
