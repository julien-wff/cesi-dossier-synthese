<script lang="ts">
    import type { Letter } from '$lib/types/grades';

    export let grade: Letter | null;
</script>


<div class="grid place-items-center w-9"
     class:bg-blue-500={grade === 'A'}
     class:bg-green-500={grade === 'B'}
     class:bg-amber-500={grade === 'C'}
     class:bg-red-500={grade === 'D'}>

    <select class="block w-9 h-full appearance-none bg-transparent text-center text-white cursor-pointer"
            bind:value={grade}>
        <option value={null}></option>
        <option value="A" class="text-blue-600">A</option>
        <option value="B" class="text-green-600">B</option>
        <option value="C" class="text-amber-600">C</option>
        <option value="D" class="text-red-600">D</option>
    </select>
</div>
