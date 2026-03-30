<script>
    import {
        CheckConfig,
        CenterWindow,
        GetVPNLogs,
        GetAppVersion,
    } from "../wailsjs/go/main/App.js";
    import { onMount } from "svelte";
    import { fade, fly, blur } from "svelte/transition";
    import { cubicOut } from "svelte/easing";
    import {
        currentScreen,
        navigate,
        theme,
        vpnLogs,
        appLogs,
        appVersion,
        checkConnectionStatus,
        isConnecting,
    } from "./store.js";
    import { EventsOn } from "../wailsjs/runtime/runtime.js";

    import Setup from "./screens/Setup.svelte";
    import Password from "./screens/Password.svelte";
    import Main from "./screens/Main.svelte";
    import Logs from "./screens/Logs.svelte";
    import Settings from "./screens/Settings.svelte";
    import { AppState } from "./constants.js";

    const MAX_LOG_SIZE = 200000; // Limit for UI performance (chars)

    // checking if config exists on app startup
    onMount(async () => {
        const initialLogs = await GetVPNLogs();
        if (initialLogs) vpnLogs.set(initialLogs);

        const version = await GetAppVersion();
        if (version) appVersion.set(version);

        try {
            const configExist = await CheckConfig();
            if (configExist) {
                // go to master password screen
                navigate(AppState.PASSWORD);
            } else {
                // go to setup screen
                navigate(AppState.SETUP);
            }
        } catch (err) {
            console.error("Initialization error:", err);
            navigate(AppState.SETUP);
        }

        // Global logging listeners
        EventsOn("vpn-log", (log) => {
            vpnLogs.update((n) => {
                const combined = n + log;
                return combined.length > MAX_LOG_SIZE
                    ? combined.slice(combined.length - MAX_LOG_SIZE)
                    : combined;
            });
        });

        EventsOn("app-log", (msg) => {
            const timestamp = new Date().toLocaleTimeString();
            const entry = `[${timestamp}] ${msg}\n`;
            appLogs.update((n) => {
                const combined = n + entry;
                return combined.length > MAX_LOG_SIZE
                    ? combined.slice(combined.length - MAX_LOG_SIZE)
                    : combined;
            });
        });

        EventsOn("connection-processing", (val) => {
            isConnecting.set(val);
        });

        // Periodic status syncing
        checkConnectionStatus(true);
        const interval = setInterval(() => checkConnectionStatus(), 1000);
        return () => clearInterval(interval);
    });

    const screens = {
        [AppState.SETUP]: Setup,
        [AppState.PASSWORD]: Password,
        [AppState.MAIN]: Main,
        [AppState.LOGS]: Logs,
        [AppState.SETTINGS]: Settings,
    };
</script>

<svelte:window on:contextmenu|preventDefault />

<div class="app-container">
    {#key $currentScreen}
        <div
            class="screen-wrapper"
            in:blur={{ amount: 5, duration: 400, delay: 200, easing: cubicOut }}
            out:fade={{ duration: 200 }}
        >
            <svelte:component this={screens[$currentScreen]} />
        </div>
    {/key}
</div>

<style>
    :global(html),
    :global(body) {
        margin: 0;
        background: var(--bg-app);
        overscroll-behavior: none;
        font-family:
            "Inter",
            system-ui,
            -apple-system,
            sans-serif;
    }

    :global(button) {
        outline: none;
    }

    :global(:root) {
        --bg-app: #f4f7f6;
        --bg-card: #ffffff;
        --text-main: #1a1a1a;
        --text-muted: #666666;
        --accent: #646cff;
        --accent-hover: #535bf2;
        --border: #e0e0e0;
        --shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1),
            0 2px 4px -1px rgba(0, 0, 0, 0.06);
        --shadow-lg: 0 10px 15px -3px rgba(0, 0, 0, 0.1),
            0 4px 6px -2px rgba(0, 0, 0, 0.05);
        --input-bg: #ffffff;
        --success: #22c55e;
        --danger: #ef4444;
        --tooltip-bg: #1e293b;
        --tooltip-text: #ffffff;
        --radius-sm: 8px;
        --radius-md: 12px;
        --radius-lg: 20px;
    }

    :global(.dark-theme) {
        --bg-app: #0f172a;
        --bg-card: #1e293b;
        --text-main: #f8fafc;
        --text-muted: #94a3b8;
        --accent: #818cf8;
        --accent-hover: #a5b4fc;
        --border: #334155;
        --shadow: 0 4px 20px rgba(0, 0, 0, 0.3);
        --shadow-lg: 0 10px 25px rgba(0, 0, 0, 0.4);
        --input-bg: #0f172a;
        --success: #4ade80;
        --danger: #f87171;
        --tooltip-bg: #334155;
        --tooltip-text: #f8fafc;
    }

    :global(.btn) {
        display: inline-flex;
        align-items: center;
        justify-content: center;
        gap: 0.5rem;
        padding: 0.65rem 1.25rem;
        border-radius: var(--radius-md);
        font-size: 0.9rem;
        font-weight: 600;
        cursor: pointer;
        transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
        border: 1px solid transparent;
        font-family: inherit;
        line-height: 1;
        background: var(--bg-card);
        color: var(--text-main);
        box-shadow: var(--shadow);
    }

    :global(.btn:hover) {
        box-shadow: var(--shadow-lg);
    }

    :global(.btn:active) {
        transform: translateY(0);
        box-shadow: var(--shadow);
    }

    :global(.btn:disabled) {
        opacity: 0.5;
        cursor: not-allowed;
        transform: none !important;
        box-shadow: none !important;
    }

    :global(.btn-primary) {
        background: var(--accent);
        color: white;
    }

    :global(.btn-primary:hover) {
        background: var(--accent-hover);
    }

    :global(.btn-secondary) {
        background: var(--bg-card);
        border: 1px solid var(--border);
        color: var(--text-main);
    }

    :global(.btn-secondary:hover) {
        background: var(--bg-app);
    }

    :global(.btn-danger) {
        background: var(--bg-card);
        border: 1px solid var(--border);
        color: var(--danger);
    }

    :global(.btn-danger:hover) {
        background: #fffafa;
    }

    :global(.dark-theme .btn-danger:hover) {
        background: rgba(248, 113, 113, 0.05);
    }

    :global(.btn-icon) {
        width: 36px;
        height: 36px;
        padding: 0;
        border-radius: 10px;
        display: flex;
        align-items: center;
        justify-content: center;
        color: var(--text-muted);
    }

    :global(.btn-icon:hover) {
        color: var(--accent);
    }

    :global(.btn-text) {
        background: none;
        border: none;
        padding: 4px 8px;
        color: var(--text-muted);
        box-shadow: none;
        font-weight: 500;
        font-size: 0.8rem;
        cursor: pointer;
        transition: color 0.2s ease;
    }

    :global(.btn-text:hover) {
        color: var(--accent);
    }

    :global([data-tooltip]) {
        position: relative;
    }

    :global([data-tooltip]::after) {
        content: attr(data-tooltip);
        position: absolute;
        bottom: -32px;
        left: 50%;
        transform: translateX(-50%) translateY(5px);
        background: var(--tooltip-bg);
        color: var(--tooltip-text);
        padding: 4px 10px;
        border-radius: 6px;
        font-size: 0.75rem;
        font-weight: 500;
        white-space: nowrap;
        opacity: 0;
        visibility: hidden;
        transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
        z-index: 200;
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
        pointer-events: none;
    }

    :global([data-tooltip]:hover::after),
    :global([data-tooltip].tooltip-visible::after) {
        opacity: 1;
        visibility: visible;
        transform: translateX(-50%) translateY(0);
    }

    :global(h2) {
        margin: 0;
        font-size: 1.15rem;
        font-weight: 700;
        color: var(--text-main);
        letter-spacing: -0.01em;
        text-align: center;
    }

    .app-container {
        position: relative;
        width: 100%;
        height: 100vh;
        overflow: hidden;
        background: var(--bg-app);
        color: var(--text-main);
        transition:
            background 0.3s ease,
            color 0.3s ease;
    }

    .screen-wrapper {
        position: absolute;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        display: flex;
        flex-direction: column;
    }
</style>
