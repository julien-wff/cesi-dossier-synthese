import { json } from '@sveltejs/kit';
import puppeteer, { type Browser, type Page } from 'puppeteer';
import { v4 as uuid } from 'uuid';
import type { RequestHandler } from './$types';


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
    let browser: Browser, page: Page;

    try {
        browser = await puppeteer.launch({
            executablePath: import.meta.env.VITE_CHROMIUM_PATH || undefined,
        });
        page = await browser.newPage();
        await page.setViewport({ width: 1280, height: 720 });
        await page.setUserAgent('Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109303030 Safari/537.36');

        user.send('loading_email_page');

        // Login page
        await page.goto('https://ent.cesi.fr/');
        await page.waitForSelector('input#login', { timeout: 5000 });

        // Input login
        await page.type('input#login', data.email);
        await page.click('#submit');

        user.send('loading_login_page');

        // Input password
        await page.waitForSelector('input#passwordInput', { timeout: 5000 });
        await page.type('input#passwordInput', data.password);
        await page.click('#submitButton');

        // Wait for the page to load
        await page.waitForNavigation({ timeout: 5000, waitUntil: 'networkidle2' });

        // Goto the page
        await page.goto('https://ent.cesi.fr/cpi-a2-info-22-23-nancy-ny2ap201/dossiers-de-synthese', { waitUntil: 'networkidle2' });

    } catch (error) {
        console.log(error);
        user.send('error', { error: (error as Error).message });
        if (typeof browser !== 'undefined' && typeof page !== 'undefined') {
            await page.screenshot({ path: 'error.png' });
        }
    } finally {
        if (typeof browser !== 'undefined')
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
