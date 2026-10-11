import { writable, get } from 'svelte/store';
import { AppState } from "./constants.js";
import { GetConfig, IsConnected, SyncStatus } from "../wailsjs/go/main/App.js";
import { WindowSetLightTheme, WindowSetDarkTheme, WindowSetSystemDefaultTheme } from "../wailsjs/runtime/runtime.js";

// Navigation
export const currentScreen = writable(null);

export function navigate(screen) {
    if (Object.values(AppState).includes(screen)) {
        currentScreen.set(screen);
    }
}

// Theme: 'system' | 'light' | 'dark'. 'system' follows the OS appearance (prefers-color-scheme).
// Old installs stored 'light'/'dark' under the same key — those values are still valid.
function readTheme() {
    try {
        return localStorage.getItem('theme') || 'system';
    } catch {
        return 'system';
    }
}
export const theme = writable(readTheme());

theme.subscribe(val => {
    try {
        localStorage.setItem('theme', val);
    } catch {
        /* storage unavailable: keep the choice for this session only */
    }
    const root = document.documentElement;
    if (val === 'light' || val === 'dark') root.setAttribute('data-theme', val);
    else root.removeAttribute('data-theme');
    // Native title bar / window chrome (takes effect on Windows; macOS follows the system).
    try {
        if (val === 'light') WindowSetLightTheme();
        else if (val === 'dark') WindowSetDarkTheme();
        else WindowSetSystemDefaultTheme();
    } catch {
        /* runtime not ready (e.g. vite preview in a browser) */
    }
});

// App State
export const isConnected = writable(false);
export const isConnecting = writable(false);
export const vpnLogs = writable("");
export const appLogs = writable("");
export const appVersion = writable("");
export const hasCheckedUpdate = writable(false);
export const updateRelease = writable(null);
export const config = writable({
    credentials_file_path: "",
    show_console_log: false,
});

// Timer State
export const connectionSeconds = writable(0);
let timerInterval = null;

export function startGlobalTimer() {
    if (timerInterval) return;
    timerInterval = setInterval(() => {
        connectionSeconds.update(n => n + 1);
    }, 1000);
}

export function stopGlobalTimer() {
    if (timerInterval) {
        clearInterval(timerInterval);
        timerInterval = null;
    }
    connectionSeconds.set(0);
}

export async function checkConnectionStatus(force = false) {
    if (get(isConnecting)) return;
    try {
        const status = force ? await SyncStatus() : await IsConnected();
        const currentStatus = get(isConnected);
        if (status !== currentStatus) {
            isConnected.set(status);
            if (status) startGlobalTimer();
            else stopGlobalTimer();
        }
    } catch (err) {
        console.error("Failed to get status:", err);
    }
}

// Re-verify the real VPN state when the window regains focus, since the cached
// status can go stale while the app was hidden/backgrounded (e.g. system sleep).
const MIN_FORCE_CHECK_INTERVAL = 5000; // don't re-check more than once per 5s
const MIN_HIDDEN_DURATION = 10000; // ignore quick alt-tabs, only react to longer absences
let hiddenAt = null;
let lastFocusCheck = 0;

export function markWindowHidden() {
    if (hiddenAt === null) hiddenAt = Date.now();
}

export function checkConnectionOnFocus() {
    const now = Date.now();
    const hiddenDuration = hiddenAt !== null ? now - hiddenAt : 0;
    hiddenAt = null;

    if (hiddenDuration < MIN_HIDDEN_DURATION) return;
    if (now - lastFocusCheck < MIN_FORCE_CHECK_INTERVAL) return;

    lastFocusCheck = now;
    checkConnectionStatus(true);
}

export async function refreshConfig() {
    try {
        const c = await GetConfig();
        if (c) config.set(c);
        return c;
    } catch (err) {
        console.error("Failed to refresh config in store:", err);
        return null;
    }
}