import { json } from '@sveltejs/kit';
import fs from 'fs/promises';
import path from 'path';
import puppeteer, { type Browser, type Page } from 'puppeteer';
import { v4 as uuid } from 'uuid';
import type { RequestHandler } from './$types';


const PPTR_TIMEOUT = 10000;


const userSessions = new Map<string, UserSession>();

export const GET = (async ({ getClientAddress }) => {
    const user = new UserSession(getClientAddress());
    userSessions.set(user.id, user);

    const stream = new ReadableStream({
        start(controller: ReadableStreamDefaultController<string>) {
            user.streamController = controller;
            user.send('connected', { id: user.id });
        },
        cancel() {
            userSessions.delete(user.id);
        },
    });

    return new Response(stream, {
        headers: { 'Content-Type': 'text/event-stream' },
    });
}) satisfies RequestHandler;


export const POST = (async ({ request, getClientAddress }) => {
    let data;
    try {
        data = await request.json();
    } catch (error) {
        return json({ error: 'Invalid JSON' }, { status: 400 });
    }

    const user = userSessions.get(data.id);

    if (!user) {
        return json({ error: 'User not found' }, { status: 404 });
    }

    if (user.address !== getClientAddress()) {
        return json({ error: 'Unauthorized' }, { status: 401 });
    }

    if (!data.email || !data.password) {
        return json({ error: 'Invalid data' }, { status: 400 });
    }

    user.send('loading_browser');
    let browser: Browser | null = null, page: Page | null = null;

    try {
        browser = await puppeteer.launch({
            executablePath: import.meta.env.VITE_CHROMIUM_PATH || undefined,
            headless: false,
        });
        page = await browser.newPage();
        await page.setViewport({ width: 1000, height: 500 });

        await page.setUserAgent('Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109303030 Safari/537.36');

        user.send('loading_email_page');
        // Login page
        await page.goto('https://ent.cesi.fr/', { timeout: PPTR_TIMEOUT });

        await page.waitForSelector('input#login', { timeout: PPTR_TIMEOUT });
        // Input login
        await page.type('input#login', data.email);

        await page.click('#submit');

        user.send('loading_login_page');
        // Input password
        await page.waitForSelector('input#passwordInput', { timeout: PPTR_TIMEOUT });
        await page.type('input#passwordInput', data.password);

        await page.click('#submitButton');
        // Wait for the page to load

        await page.waitForNavigation({ timeout: PPTR_TIMEOUT, waitUntil: 'networkidle2' });
        // Goto the pdf page
        user.send('loading_pdf_page');
        const filesLink = await page.evaluate(() =>
            (document.querySelector('a[title="Dossiers de synthèse"]') as HTMLAnchorElement)?.href,
        );
        if (!filesLink) {
            throw new Error('Impossible de trouver le lien vers les dossiers de synthèse');
        }

        await page.goto(filesLink, { waitUntil: 'networkidle2', timeout: PPTR_TIMEOUT });
        // Get the pdf link

        const pdfLink = await page.evaluate(() => {
            return ([ ...document.querySelectorAll('.semestres__notes a[title="Lien vers le dossier"]') ] as HTMLAnchorElement[])
                .map(a => ({ href: a.href, text: a.innerText }));
        });
        if (!pdfLink || !pdfLink.length) {
            throw new Error('Impossible de trouver le lien vers le dossier de synthèse');

        }
        // Download the pdf
        console.log(pdfLink);

        user.send('loading_pdf');
        // Intercept the pdf request
        const client = await page.target().createCDPSession();
        // Enable downloading
        await client.send('Page.setDownloadBehavior', {
            behavior: 'allow',
            downloadPath: path.resolve('./temp'),
        });

        await page.setRequestInterception(true);

        page.on('request', async (request) => {
            await request.continue();
        });

        page.on('response', async (request) => {
            const headers = request.headers();
            if (headers['content-type'] !== 'application/pdf')
                return;
            // Wait 1s
            await new Promise(resolve => setTimeout(resolve, 1000));
            // Rename the file to an uuid if it exists
            const filePath = path.resolve('./temp', headers['content-disposition'].split('filename=')[1].replace(/"/g, ''));
            const newFilePath = path.resolve('./temp', `${uuid()}.pdf`);
            try {
                await fs.rename(filePath, newFilePath);
            } catch (_) {
            }
        });

        // Click on the link
        await page.click(`a[href="${pdfLink[0].href}"]`);

        await new Promise(resolve => setTimeout(resolve, 5000));

        throw new Error('Not implemented');

    } catch (error) {
        let screenshot = browser && await page?.screenshot({ encoding: 'base64', quality: 20, type: 'webp' });
        user.send('error', { error: (error as Error).message, screenshot });
    } finally {
        await browser?.close();
    }

    return json({ message: 'ok' });
}) satisfies RequestHandler;


class UserSession {
    constructor(public address: string) {
    }

    id = uuid();
    streamController?: ReadableStreamDefaultController<string>;

    send(event: string, data: any = null) {
        this.streamController?.enqueue(`data: ${JSON.stringify({ event, data })}\n\n`);
    }
}
