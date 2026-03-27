import { writable, get } from 'svelte/store';
import { AppState } from "./constants.js";
import { GetConfig, IsConnected, SyncStatus } from "../wailsjs/go/main/App.js";

// Navigation
export const currentScreen = writable(null);

export function navigate(screen) {
    if (Object.values(AppState).includes(screen)) {
        currentScreen.set(screen);
    }
}

// Theme
const initialTheme = localStorage.getItem('theme') || 'dark';
export const theme = writable(initialTheme);

theme.subscribe(val => {
    localStorage.setItem('theme', val);
    if (val === 'dark') {
        document.documentElement.classList.add('dark-theme');
    } else {
        document.documentElement.classList.remove('dark-theme');
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