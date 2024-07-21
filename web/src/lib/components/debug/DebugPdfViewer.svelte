<script lang="ts">
    import { DebugLineDirection, type DebugResponse, DrawMode } from '$lib/types/debug';
    import { onDestroy, onMount } from 'svelte';

    export let margin = 16;
    export let resolution = 1;
    export let data: DebugResponse | null = null;
    export let page: number = 0;
    export let mode: DrawMode = DrawMode.Page;
    export let debugColors = false;
    export let displaySingleText = false;
    export let textIndex = 0;
    export let displaySingleLine = false;
    export let lineIndex = 0;
    export let displaySingleSquare = false;
    export let squareIndex = 0;
    export let showNeighbours = false;

    let canvas: HTMLCanvasElement;
    let hidden = true;

    $: if (canvas && data)
        drawCanvas(
            data,
            page,
            mode,
            debugColors,
            !displaySingleText ? -1 : textIndex,
            !displaySingleLine ? -1 : lineIndex,
            !displaySingleSquare ? -1 : squareIndex,
            showNeighbours,
        );

    function drawCanvas(data: DebugResponse,
                        page: number,
                        mode: DrawMode,
                        debugColors: boolean,
                        displaySingleText: number,
                        displaySingleLine: number,
                        displaySingleSquare: number,
                        showNeighbours: boolean) {
        hidden = false;
        const ctx = canvas?.getContext('2d');

        if (!ctx)
            return;

        let scaleFactor = getScaleFactor(data, page);
        const FADED_OPACITY = 0.40;

        canvas.style.width = `${data.pages[page].size.width * scaleFactor}px`;
        canvas.style.height = `${data.pages[page].size.height * scaleFactor}px`;

        scaleFactor *= resolution;

        canvas.width = data.pages[page].size.width * scaleFactor;
        canvas.height = data.pages[page].size.height * scaleFactor;

        if (mode === DrawMode.Page) {
            for (let i = 0; i < data.pages[page].text.length; i++) {
                const text = data.pages[page].text[i];
                if (displaySingleText !== -1 && displaySingleText !== i)
                    continue;

                ctx.font = `${text.font_size * scaleFactor}px sans-serif`;
                ctx.fillStyle = debugColors ? 'red' : 'black';
                ctx.globalAlpha = debugColors ? FADED_OPACITY : 1;
                ctx.fillText(text.content, text.position.x * scaleFactor, text.position.y * scaleFactor);
            }

            for (let i = 0; i < data.pages[page].lines.length; i++) {
                const line = data.pages[page].lines[i];
                if (displaySingleLine !== -1 && displaySingleLine !== i)
                    continue;

                ctx.strokeStyle = debugColors ? 'blue' : 'black';
                ctx.globalAlpha = debugColors ? FADED_OPACITY : 1;
                ctx.lineWidth = scaleFactor;
                ctx.beginPath();
                ctx.moveTo(line.x1 * scaleFactor, line.y1 * scaleFactor);
                ctx.lineTo(line.x2 * scaleFactor, line.y2 * scaleFactor);
                ctx.stroke();
            }

            for (let i = 0; i < data.pages[page].rectangles.length; i++) {
                const rect = data.pages[page].rectangles[i];
                if (displaySingleSquare !== -1 && displaySingleSquare !== i)
                    continue;

                ctx.strokeStyle = debugColors ? 'green' : 'black';
                ctx.globalAlpha = debugColors ? FADED_OPACITY : 1;
                ctx.lineWidth = scaleFactor;
                const x1 = rect.position.x * scaleFactor;
                const y1 = rect.position.y * scaleFactor;
                const x2 = (rect.position.x + rect.size.width) * scaleFactor;
                const y2 = (rect.position.y + rect.size.height) * scaleFactor;
                ctx.strokeRect(x1, y1, x2 - x1, y2 - y1);
            }
        } else if (mode === DrawMode.Line) {
            const pickedLine = displaySingleLine !== -1 ? data.lines[page].lines[displaySingleLine] : null;
            const pickedNeighbours = pickedLine
                ? [ ...pickedLine.start_neighbours_ids, ...pickedLine.end_neighbours_ids ]
                : null;

            for (let i = 0; i < data.lines[page].lines.length; i++) {
                const line = data.lines[page].lines[i];
                if (!(displaySingleLine === -1 || displaySingleLine === i || (showNeighbours && pickedNeighbours?.includes(line.id))))
                    continue;

                ctx.strokeStyle = debugColors
                    ? line.direction === DebugLineDirection.Horizontal ? 'blue' : 'red'
                    : 'black';
                ctx.globalAlpha = debugColors && displaySingleLine !== i
                    ? FADED_OPACITY
                    : 1;
                ctx.lineWidth = scaleFactor * (Number(debugColors && displaySingleLine === i) + 1);
                ctx.beginPath();
                ctx.moveTo(line.x1 * scaleFactor, line.y1 * scaleFactor);
                ctx.lineTo(line.x2 * scaleFactor, line.y2 * scaleFactor);
                ctx.stroke();

                // Draw points at beginning and end of line
                if (debugColors) {
                    ctx.fillStyle = 'purple';
                    ctx.beginPath();
                    ctx.arc(line.x1 * scaleFactor, line.y1 * scaleFactor, 3 * scaleFactor, 0, 2 * Math.PI);
                    ctx.fill();

                    ctx.beginPath();
                    ctx.arc(line.x2 * scaleFactor, line.y2 * scaleFactor, 3 * scaleFactor, 0, 2 * Math.PI);
                    ctx.fill();
                }
            }
        } else if (mode === DrawMode.Square) {
            const squares = data.squares[page].squares.flat();
            for (let i = 0; i < squares.length; i++) {
                const square = squares[i];
                if (displaySingleSquare !== -1 && displaySingleSquare !== i)
                    continue;

                ctx.strokeStyle = debugColors ? 'red' : 'black';
                ctx.globalAlpha = debugColors ? FADED_OPACITY : 1;
                ctx.lineWidth = scaleFactor;
                const x1 = square.x1 * scaleFactor;
                const y1 = square.y1 * scaleFactor;
                const x2 = square.x2 * scaleFactor;
                const y2 = square.y2 * scaleFactor;
                ctx.strokeRect(x1, y1, x2 - x1, y2 - y1);

                // Add text content
                ctx.font = `${8 * scaleFactor}px sans-serif`;
                ctx.fillStyle = debugColors ? 'blue' : 'black';
                ctx.globalAlpha = debugColors ? FADED_OPACITY : 1;
                ctx.textAlign = 'center';
                ctx.textBaseline = 'middle';
                // Prevent text from being drawn outside of square by truncating it
                let text = square.content;
                if (ctx.measureText(square.content).width > (x2 - x1) * 0.9) {
                    const words = square.content.split(' ');
                    text = '';
                    for (const word of words) {
                        if (ctx.measureText(text + word).width > (x2 - x1) * 0.9) {
                            text += '...';
                            break;
                        }
                        text += word + ' ';
                    }
                }
                ctx.fillText(text, (x1 + x2) / 2, (y1 + y2) / 2);

                // Draw points at corners of square
                if (debugColors) {
                    ctx.fillStyle = 'purple';
                    ctx.beginPath();
                    ctx.arc(x1, y1, 3 * scaleFactor, 0, 2 * Math.PI);
                    ctx.fill();

                    ctx.beginPath();
                    ctx.arc(x2, y1, 3 * scaleFactor, 0, 2 * Math.PI);
                    ctx.fill();

                    ctx.beginPath();
                    ctx.arc(x1, y2, 3 * scaleFactor, 0, 2 * Math.PI);
                    ctx.fill();

                    ctx.beginPath();
                    ctx.arc(x2, y2, 3 * scaleFactor, 0, 2 * Math.PI);
                    ctx.fill();
                }
            }
        }
    }

    function getScaleFactor(data: DebugResponse, page: number) {
        const canvasParent = canvas.parentElement!;
        const parentWidth = canvasParent.clientWidth;
        const parentHeight = canvasParent.clientHeight;

        return Math.min(
            (parentWidth - margin * 2) / data.pages[page].size.width,
            (parentHeight - margin * 2) / data.pages[page].size.height,
        );
    }

    let parentObserver: ResizeObserver;

    onMount(() => {
        // if parent element size changes, redraw canvas
        parentObserver = new ResizeObserver(() => {
            if (data) drawCanvas(
                data,
                page,
                mode,
                debugColors,
                !displaySingleText ? -1 : textIndex,
                !displaySingleLine ? -1 : lineIndex,
                !displaySingleSquare ? -1 : squareIndex,
                showNeighbours,
            );
        });

        parentObserver.observe(canvas.parentElement!);
    });

    onDestroy(() => {
        parentObserver?.disconnect();
    });
</script>

<canvas bind:this={canvas} class="bg-white" class:hidden/>
