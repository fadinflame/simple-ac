<script>
    import {
        CheckConfig,
        CenterWindow,
        GetVPNLogs,
        GetAppVersion,
    } from "../wailsjs/go/main/App.js";
    import { onMount } from "svelte";
    import { fade } from "svelte/transition";
    import {
        currentScreen,
        navigate,
        theme,
        vpnLogs,
        appLogs,
        appVersion,
        checkConnectionStatus,
        checkConnectionOnFocus,
        markWindowHidden,
        isConnecting,
    } from "./store.js";
    import { EventsOn, Environment } from "../wailsjs/runtime/runtime.js";

    import Setup from "./screens/Setup.svelte";
    import Password from "./screens/Password.svelte";
    import Main from "./screens/Main.svelte";
    import Logs from "./screens/Logs.svelte";
    import Settings from "./screens/Settings.svelte";
    import EditCredentials from "./screens/EditCredentials.svelte";
    import { AppState } from "./constants.js";

    const MAX_LOG_SIZE = 200000; // Limit for UI performance (chars)

    // checking if config exists on app startup
    onMount(async () => {
        // Platform drives small layout differences (traffic-light inset on macOS).
        Environment()
            .then((env) => (document.documentElement.dataset.platform = env.platform))
            .catch(() => {});

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
            const timestamp = new Date().toLocaleTimeString("en-GB", { hour12: false });
            const entry = `${timestamp}  ${msg}\n`;
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

        // Re-verify the real state when the window comes back into focus, so a
        // stale "connected" state doesn't linger after e.g. the system slept.
        const handleVisibilityChange = () => {
            if (document.hidden) markWindowHidden();
            else checkConnectionOnFocus();
        };
        window.addEventListener("blur", markWindowHidden);
        window.addEventListener("focus", checkConnectionOnFocus);
        document.addEventListener("visibilitychange", handleVisibilityChange);

        return () => {
            clearInterval(interval);
            window.removeEventListener("blur", markWindowHidden);
            window.removeEventListener("focus", checkConnectionOnFocus);
            document.removeEventListener("visibilitychange", handleVisibilityChange);
        };
    });

    const screens = {
        [AppState.SETUP]: Setup,
        [AppState.PASSWORD]: Password,
        [AppState.MAIN]: Main,
        [AppState.LOGS]: Logs,
        [AppState.EDIT_CREDENTIALS]: EditCredentials,
        [AppState.SETTINGS]: Settings,
    };
</script>

<svelte:window on:contextmenu|preventDefault />

<div class="app-container ac-root">
    {#key $currentScreen}
        <div class="screen-wrapper" in:fade={{ duration: 120 }}>
            <svelte:component this={screens[$currentScreen]} />
        </div>
    {/key}
</div>

<style>
    .app-container {
        position: relative;
        width: 100%;
        height: 100vh;
        overflow: hidden;
        background: var(--bg-window);
    }

    .screen-wrapper {
        position: absolute;
        inset: 0;
        display: flex;
        flex-direction: column;
    }
</style>
