<script lang="ts">
    export let checked = false;
    export let label: string;
</script>

<label class="block select-none">
    <input type="checkbox" bind:checked class="mr-2"/>
    {label}
</label>
