<script>
    import { onMount, onDestroy } from "svelte";
    import { fade, fly, scale } from "svelte/transition";
    import { cubicOut } from "svelte/easing";
    import {
        WindowSetTitle,
        BrowserOpenURL,
    } from "../../wailsjs/runtime/runtime.js";
    import {
        Connect,
        Disconnect,
        IsConnected,
        SyncStatus,
        Lock,
        UpdateConfig,
    } from "../../wailsjs/go/main/App.js";
    import { EventsOn, EventsOff } from "../../wailsjs/runtime/runtime.js";
    import {
        currentScreen,
        config,
        refreshConfig,
        isConnected,
        theme,
        appVersion,
        hasCheckedUpdate,
        updateRelease,
        connectionSeconds,
        startGlobalTimer,
        stopGlobalTimer,
        isConnecting,
        checkConnectionStatus,
    } from "../store.js";
    import { AppState } from "../constants.js";

    let connectionError = false;
    let errorTimeout;

    onMount(async () => {
        WindowSetTitle("");
        await refreshConfig();
    });

    onDestroy(() => {
        if (errorTimeout) clearTimeout(errorTimeout);
    });

    $: formattedTime = (() => {
        const h = Math.floor($connectionSeconds / 3600)
            .toString()
            .padStart(2, "0");
        const m = Math.floor(($connectionSeconds % 3600) / 60)
            .toString()
            .padStart(2, "0");
        const s = ($connectionSeconds % 60).toString().padStart(2, "0");
        return `${h}:${m}:${s}`;
    })();

    async function toggleConnection() {
        if ($isConnecting) return;
        try {
            if (!$isConnected) {
                $isConnecting = true;
                connectionError = false;
                const result = await Connect();
                if (!result) {
                    showError();
                }
            } else {
                await Disconnect();
            }
        } catch (err) {
            console.error("VPN Error:", err);
            showError();
        } finally {
            $isConnecting = false;
            await checkConnectionStatus(true);
        }
    }

    function showError() {
        connectionError = true;
        if (errorTimeout) clearTimeout(errorTimeout);
        errorTimeout = setTimeout(() => {
            connectionError = false;
        }, 5000);
    }

    function goToLogs() {
        $currentScreen = AppState.LOGS;
    }

    function goToSettings() {
        $currentScreen = AppState.SETTINGS;
    }

    async function lockApp() {
        try {
            await Lock();
            $currentScreen = AppState.PASSWORD;
        } catch (err) {
            console.error("Failed to lock app:", err);
        }
    }

</script>

<div class="main-container">
    <header>
        <div class="logo">
            <svg
                xmlns="http://www.w3.org/2000/svg"
                width="24"
                height="24"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
            >
                <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"></path>
            </svg>
            <span>SimpleAC</span>
        </div>
        <div class="header-actions">
            <button
                class="btn btn-icon"
                on:click={goToSettings}
                data-tooltip="Settings"
            >
                <svg
                    xmlns="http://www.w3.org/2000/svg"
                    width="20"
                    height="20"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                >
                    <circle cx="12" cy="12" r="3"></circle>
                    <path
                        d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"
                    ></path>
                </svg>
            </button>
            <button
                class="btn btn-icon"
                on:click={goToLogs}
                data-tooltip="Connection Logs"
            >
                <svg
                    xmlns="http://www.w3.org/2000/svg"
                    width="20"
                    height="20"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                >
                    <path d="M4 19h16"></path><path d="M4 14h16"></path><path
                        d="M4 9h16"
                    ></path><path d="M4 4h16"></path>
                </svg>
            </button>
            <button
                class="btn btn-icon lock"
                on:click={lockApp}
                data-tooltip="Lock"
            >
                <svg
                    xmlns="http://www.w3.org/2000/svg"
                    width="20"
                    height="20"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                >
                    <rect x="3" y="11" width="18" height="11" rx="2" ry="2"
                    ></rect>
                    <path d="M7 11V7a5 5 0 0 1 10 0v4"></path>
                </svg>
            </button>
        </div>
    </header>

    <div class="banner-overlay">
        {#if connectionError}
            <div
                class="error-banner"
                in:fly={{ y: -20, duration: 400 }}
                out:fade={{ duration: 200 }}
                on:click={() => (connectionError = false)}
                role="button"
                tabindex="0"
                on:keydown={(e) =>
                    e.key === "Enter" && (connectionError = false)}
            >
                <div class="error-info">
                    <svg
                        xmlns="http://www.w3.org/2000/svg"
                        width="18"
                        height="18"
                        viewBox="0 0 24 24"
                        fill="none"
                        stroke="currentColor"
                        stroke-width="2.5"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                    >
                        <circle cx="12" cy="12" r="10"></circle>
                        <line x1="12" y1="8" x2="12" y2="12"></line>
                        <line x1="12" y1="16" x2="12.01" y2="16"></line>
                    </svg>
                    <span class="error-title"
                        >Connection failed. Check logs for details.</span
                    >
                </div>
            </div>
        {/if}

        {#if $updateRelease}
            <div
                class="update-banner"
                in:fly={{ y: -20, duration: 400 }}
                on:click={() => {
                    BrowserOpenURL($updateRelease.html_url);
                    $updateRelease = null;
                }}
                role="button"
                tabindex="0"
                on:keydown={(e) =>
                    e.key === "Enter" &&
                    BrowserOpenURL($updateRelease.html_url)}
            >
                <button
                    class="btn btn-icon close-update"
                    on:click|stopPropagation={() => ($updateRelease = null)}
                    aria-label="Close update banner"
                >
                    <svg
                        xmlns="http://www.w3.org/2000/svg"
                        width="16"
                        height="16"
                        viewBox="0 0 24 24"
                        fill="none"
                        stroke="currentColor"
                        stroke-width="2.5"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                    >
                        <line x1="18" y1="6" x2="6" y2="18"></line><line
                            x1="6"
                            y1="6"
                            x2="18"
                            y2="18"
                        ></line>
                    </svg>
                </button>
                <div class="update-info">
                    <span class="update-title"
                        >Update Available: {$updateRelease.tag_name}</span
                    >
                    <span class="update-link">Click to view on GitHub</span>
                </div>
            </div>
        {/if}
    </div>

    <main>
        <div class="status-card" class:connected={$isConnected}>
            <div class="status-indicator">
                <div class="dot" class:pulse={$isConnecting}></div>
                <span
                    >{$isConnecting
                        ? $isConnected
                            ? "Disconnecting..."
                            : "Connecting..."
                        : $isConnected
                          ? "Connected"
                          : "Disconnected"}</span
                >
            </div>
        </div>

        <div class="connection-control">
            <button
                class="power-button"
                class:active={$isConnected}
                class:loading={$isConnecting}
                on:click={toggleConnection}
            >
                <div class="button-inner">
                    <svg
                        xmlns="http://www.w3.org/2000/svg"
                        width="48"
                        height="48"
                        viewBox="0 0 24 24"
                        fill="none"
                        stroke="currentColor"
                        stroke-width="2"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                    >
                        <path d="M18.36 6.64a9 9 0 1 1-12.73 0"></path>
                        <line x1="12" y1="2" x2="12" y2="12"></line>
                    </svg>
                </div>
                <div class="ring"></div>
            </button>
            <div class="timer" class:visible={$isConnected || $isConnecting}>
                {formattedTime}
            </div>
        </div>
    </main>

    <footer>
        <div class="footer-left"><span>{$appVersion || "v1.0.0"}</span></div>
        <button
            class="btn-text"
            on:click={() =>
                BrowserOpenURL("https://github.com/fadinflame/simple-ac")}
            aria-label="Open Repository"
        >
            <span>Made with ❤️ & Golang</span>
        </button>
    </footer>
</div>

<style>
    .main-container {
        display: flex;
        flex-direction: column;
        height: 100vh;
        background: var(--bg-app);
        padding: 1.5rem 2rem 1rem;
        box-sizing: border-box;
        position: relative;
    }

    header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 1.25rem;
    }

    .header-actions {
        display: flex;
        gap: 0.5rem;
    }

    .logo {
        display: flex;
        align-items: center;
        gap: 0.75rem;
        color: var(--text-main);
        font-weight: 700;
        font-size: 1.1rem;
    }

    .logo svg {
        color: var(--accent);
    }

    .btn-icon.lock:hover {
        color: var(--danger);
        background: #fff6f6;
    }

    :global(.dark-theme) .btn-icon.lock:hover {
        background: #2d1d1d;
    }

    main {
        flex: 1;
        display: flex;
        flex-direction: column;
        align-items: center;
        gap: 1.5rem;
        justify-content: center;
        margin-bottom: 1.5rem;
    }

    .status-card {
        background: var(--bg-card);
        padding: 0.6rem 1.25rem;
        border-radius: var(--radius-lg);
        box-shadow: var(--shadow);
        display: flex;
        align-items: center;
        justify-content: center;
        width: auto;
        min-width: 140px;
        transition: all 0.3s;
    }

    .status-card.connected {
        background: #f0fdf4;
        box-shadow: 0 4px 20px rgba(34, 197, 94, 0.1);
    }

    :global(.dark-theme) .status-card.connected {
        background: rgba(34, 197, 94, 0.1);
    }

    .status-indicator {
        display: flex;
        align-items: center;
        gap: 0.6rem;
        font-weight: 600;
        font-size: 0.85rem;
        color: var(--text-muted);
    }

    .connected .status-indicator {
        color: var(--success);
    }

    .dot {
        width: 8px;
        height: 8px;
        border-radius: 50%;
        background: var(--text-muted);
    }

    .connected .dot {
        background: var(--success);
    }
    .pulse {
        animation: pulse 1.5s infinite;
        background: var(--accent);
    }

    .timer {
        font-family: "JetBrains Mono", monospace;
        font-size: 1.25rem;
        font-weight: 700;
        color: var(--text-muted);
        transition: all 0.3s ease;
        opacity: 0.5;
        margin-top: 0.5rem;
    }

    .timer.visible {
        color: var(--text-main);
        opacity: 1;
    }

    .connection-control {
        display: flex;
        flex-direction: column;
        align-items: center;
        gap: 1rem;
    }

    .power-button {
        width: 120px;
        height: 120px;
        border-radius: 50%;
        border: none;
        background: var(--bg-card);
        color: var(--text-muted);
        cursor: pointer;
        position: relative;
        display: flex;
        align-items: center;
        justify-content: center;
        transition: all 0.3s ease;
        box-shadow: var(--shadow);
    }

    .power-button:focus {
        outline: none;
    }

    .power-button:hover {
        background: var(--bg-card);
        color: var(--accent);
        box-shadow: 0 0 20px rgba(100, 108, 255, 0.15);
    }

    :global(.dark-theme) .power-button:hover {
        background: #242f44;
    }

    .power-button:active {
        transform: translateY(1px);
        box-shadow: var(--shadow);
    }

    .power-button.active {
        color: var(--success);
        background: #f0fdf4;
        box-shadow: 0 0 15px rgba(34, 197, 94, 0.1);
    }

    .power-button.active:hover {
        background: #ecfdf5;
        box-shadow: 0 0 25px rgba(34, 197, 94, 0.2);
    }

    :global(.dark-theme) .power-button.active {
        background: rgba(34, 197, 94, 0.05);
    }

    .power-button.loading .button-inner {
        animation: spin 2s linear infinite;
        color: var(--accent);
    }

    .ring {
        position: absolute;
        inset: -8px;
        border-radius: 50%;
        border: 2px solid transparent;
        transition: all 0.3s;
    }

    .active .ring {
        border-color: var(--success);
        box-shadow: 0 0 20px rgba(34, 197, 94, 0.3);
        animation: breathe 3s ease-in-out infinite;
    }

    footer {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding-top: 1rem;
        border-top: 1px solid var(--border);
        color: var(--text-muted);
        font-size: 0.75rem;
    }

    .banner-overlay {
        position: absolute;
        top: 4.5rem;
        left: 0;
        right: 0;
        display: flex;
        flex-direction: column;
        align-items: center;
        gap: 0.75rem;
        z-index: 100;
        pointer-events: none;
    }

    .banner-overlay > * {
        pointer-events: auto;
    }

    .update-banner {
        width: 100%;
        max-width: 380px;
        background: var(--accent);
        color: white;
        padding: 0.75rem 1rem;
        border-radius: 12px;
        display: flex;
        align-items: center;
        justify-content: flex-start;
        cursor: pointer;
        position: relative;
        transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
        box-shadow: var(--shadow);
    }

    .update-banner:hover {
        box-shadow: var(--shadow-lg);
        background: var(--accent-hover);
    }

    .update-banner:active {
        box-shadow: var(--shadow);
    }

    .update-title {
        font-weight: 600;
        font-size: 0.85rem;
        display: block;
    }

    .update-link {
        font-size: 0.75rem;
        opacity: 0.9;
    }

    .update-info {
        text-align: left;
        display: flex;
        flex-direction: column;
        gap: 0.1rem;
    }

    .close-update {
        background: rgba(255, 255, 255, 0.2);
        border: none;
        width: 28px;
        height: 28px;
        padding: 0;
        border-radius: 8px;
        color: white;
        cursor: pointer;
        position: absolute;
        top: 10px;
        right: 10px;
        box-shadow: none;
    }

    .close-update:hover {
        background: rgba(255, 255, 255, 0.3);
    }

    .error-banner {
        width: 100%;
        max-width: 400px;
        background: var(--danger);
        color: white;
        padding: 0.85rem 1.25rem;
        border-radius: 14px;
        display: flex;
        align-items: center;
        justify-content: flex-start;
        cursor: pointer;
        transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
        box-shadow: var(--shadow);
    }

    .error-banner:hover {
        transform: translateY(-2px);
        box-shadow: 0 6px 20px rgba(239, 68, 68, 0.35);
    }

    .error-info {
        display: flex;
        align-items: center;
        gap: 0.75rem;
    }

    .error-title {
        font-weight: 600;
        font-size: 0.85rem;
    }

    @keyframes spin {
        from {
            transform: rotate(0deg);
        }
        to {
            transform: rotate(360deg);
        }
    }

    @keyframes breathe {
        0%,
        100% {
            transform: scale(1);
            opacity: 0.5;
        }
        50% {
            transform: scale(1.02);
            opacity: 1;
        }
    }

    @keyframes pulse {
        0% {
            transform: scale(0.95);
            box-shadow: 0 0 0 0 rgba(100, 108, 255, 0.7);
        }
        70% {
            transform: scale(1);
            box-shadow: 0 0 0 10px rgba(100, 108, 255, 0);
        }
        100% {
            transform: scale(0.95);
            box-shadow: 0 0 0 0 rgba(100, 108, 255, 0);
        }
    }
</style>
