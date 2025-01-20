<script lang="ts">
    import FileText from 'lucide-svelte/icons/file-text';

    interface Props {
        file: File | null;
        hidden?: boolean;
        shiftKey?: boolean;
    }

    let { file = $bindable(), hidden = false, shiftKey = $bindable(false) }: Props = $props();

    let inputField = $state<HTMLInputElement>();
    let drag = $state(false);

    const getFieldFile = () => inputField?.files?.[0] ?? null;

    $effect(() => {
        if (!inputField)
            return;

        if (file && !getFieldFile()) {
            const dataTransfer = new DataTransfer();
            dataTransfer.items.add(file);
            inputField.files = dataTransfer.files;
        } else if (!file && getFieldFile()) {
            inputField.value = '';
        }
    });

    function handleDrop(ev: DragEvent) {
        drag = false;
        shiftKey = ev.shiftKey;
    }

    function handleChange() {
        file = getFieldFile();
    }
</script>


<div class="border-4 border-blue-500 dark:border-blue-400 rounded-2xl relative bg-slate-100 dark:bg-transparent shadow-lg"
     class:hidden
     class:bg-blue-100={drag}
     class:dark:bg-blue-900={drag}>

    <input type="file"
           ondragenter={() => (drag = true)}
           ondragleave={() => (drag = false)}
           ondrop={handleDrop}
           onchange={handleChange}
           bind:this={inputField}
           accept="application/pdf"
           class="w-80 h-48 opacity-0 cursor-pointer">

    <div class="absolute inset-0 flex flex-col items-center justify-center gap-4 p-4 text-center pointer-events-none">
        <FileText class="w-8 h-8"/>
        <p class="font-bold">Clique ou dépose ton<br>dossier PDF ici</p>
    </div>
</div>
