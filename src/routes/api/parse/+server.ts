import type { ExtractionTable } from '$lib/types/tabula';
import { parseTabulaResult } from '$lib/utils/parse';
import { json } from '@sveltejs/kit';
import { execSync } from 'child_process';
import fs from 'fs/promises';
import { v4 as uuid } from 'uuid';
import type { RequestHandler } from './$types';

export const POST = (async ({ request }) => {
    const form = await request.formData();

    // Get file from request
    if (!form.has('file'))
        return json({ error: 'Aucun fichier n\'a été trouvé dans la requête' }, { status: 400 });

    const file = form.get('file');

    // Check file validity
    if (!(file instanceof File))
        return json({ error: 'Le fichier n\'est pas valide' }, { status: 400 });
    if (file.type !== 'application/pdf')
        return json({ error: 'Le fichier n\'est pas un PDF' }, { status: 400 });
    if (file.size > 5e4) // 50KB
        return json({ error: 'Le fichier est trop volumineux' }, { status: 400 });

    // Check the temp folder exists, and create it if not
    try {
        await fs.access('./temp');
    } catch {
        await fs.mkdir('./temp');
    }

    // Save file to disk
    const buffer = await file.arrayBuffer();
    const path = `./temp/${uuid()}.pdf`;
    await fs.writeFile(path, Buffer.from(buffer));

    // Parse the PDF
    const commandArgs = [
        '-Dorg.slf4j.simpleLogger.defaultLogLevel=off', // Disable slf4j logs
        '-Dorg.apache.commons.logging.Log=org.apache.commons.logging.impl.NoOpLog', // Disable commons-logging logs
        '-jar tabula.jar', // Tabula jar file
        '-p all', // Parse all pages
        '-f JSON', // Output as JSON
        '-l', // Use lattice mode (alternate parsing mode)
        '-u', // Use line returns in cells
        '-i', // Silent mode
        path,
    ];
    const commandResult = execSync(`java ${commandArgs.join(' ')}`).toString();
    const parsingResult = JSON.parse(commandResult) as ExtractionTable[];
    const result = parseTabulaResult(parsingResult);

    // Delete the file
    await fs.unlink(path);

    return json({ data: result });
}) satisfies RequestHandler;
