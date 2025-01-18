<script lang="ts">
    interface Props {
        checked?: boolean;
        label: string;
    }

    let { checked = $bindable(false), label }: Props = $props();
</script>

<label class="block select-none">
    <input type="checkbox" bind:checked class="mr-2"/>
    {label}
</label>
