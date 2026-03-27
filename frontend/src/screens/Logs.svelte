<script>
    import { fade, fly } from "svelte/transition";
    import { vpnLogs, appLogs, currentScreen } from "../store.js";
    import { AppState } from "../constants.js";

    let activeTab = "vpn"; // 'vpn' or 'app'
    let vpnContainer;
    let appContainer;

    function scrollToBottom(container) {
        if (!container) return;
        setTimeout(() => {
            container.scrollTop = container.scrollHeight;
        }, 10);
    }

    $: if (activeTab === "vpn" && vpnContainer) scrollToBottom(vpnContainer);
    $: if (activeTab === "app" && appContainer) scrollToBottom(appContainer);

    function goBack() {
        $currentScreen = AppState.MAIN;
    }

    function clearLogs() {
        if (activeTab === "vpn") vpnLogs.set("");
        else appLogs.set("");
    }
</script>

<div class="logs-screen" in:fade={{ duration: 300 }}>
    <header>
        <button class="btn btn-icon" on:click={goBack} aria-label="Go back">
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
                <line x1="19" y1="12" x2="5" y2="12"></line>
                <polyline points="12 19 5 12 12 5"></polyline>
            </svg>
        </button>
        <h2>System Logs</h2>
        <button
            class="btn btn-icon clear-button"
            on:click={clearLogs}
            title="Clear Logs"
        >
            <svg
                xmlns="http://www.w3.org/2000/svg"
                width="18"
                height="18"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
                ><path d="M3 6h18"></path><path
                    d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6"
                ></path><path d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2"
                ></path><line x1="10" y1="11" x2="10" y2="17"></line><line
                    x1="14"
                    y1="11"
                    x2="14"
                    y2="17"
                ></line></svg
            >
        </button>
    </header>

    <div class="tabs-container">
        <div class="segmented-control">
            <div class="indicator" class:app={activeTab === "app"}></div>
            <button
                class:active={activeTab === "vpn"}
                on:click={() => (activeTab = "vpn")}
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
                    ><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"
                    ></path></svg
                >
                <span>VPN CLI</span>
            </button>
            <button
                class:active={activeTab === "app"}
                on:click={() => (activeTab = "app")}
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
                    ><path d="M22 12h-4l-3 9L9 3l-3 9H2"></path></svg
                >
                <span>Application</span>
            </button>
        </div>
    </div>

    <main>
        {#if activeTab === "vpn"}
            {#if $vpnLogs}
                <div
                    class="log-view"
                    bind:this={vpnContainer}
                    in:fly={{ y: 5, duration: 200 }}
                >
                    <pre>{$vpnLogs}</pre>
                </div>
            {:else}
                <div class="empty" in:fade>
                    <svg
                        xmlns="http://www.w3.org/2000/svg"
                        width="48"
                        height="48"
                        viewBox="0 0 24 24"
                        fill="none"
                        stroke="currentColor"
                        stroke-width="1.5"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        ><circle cx="12" cy="12" r="10"></circle><line
                            x1="12"
                            y1="8"
                            x2="12"
                            y2="12"
                        ></line><line x1="12" y1="16" x2="12.01" y2="16"
                        ></line></svg
                    >
                    <p>No VPN logs captured yet.</p>
                </div>
            {/if}
        {:else if $appLogs}
            <div
                class="log-view"
                bind:this={appContainer}
                in:fly={{ y: 5, duration: 200 }}
            >
                <pre>{$appLogs}</pre>
            </div>
        {:else}
            <div class="empty" in:fade>
                <svg
                    xmlns="http://www.w3.org/2000/svg"
                    width="48"
                    height="48"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="1.5"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    ><circle cx="12" cy="12" r="10"></circle><line
                        x1="12"
                        y1="8"
                        x2="12"
                        y2="12"
                    ></line><line x1="12" y1="16" x2="12.01" y2="16"
                    ></line></svg
                >
                <p>No application events recorded.</p>
            </div>
        {/if}
    </main>
</div>

<style>
    .logs-screen {
        display: flex;
        flex-direction: column;
        height: 100vh;
        background: var(--bg-app);
        padding: 1.25rem;
        box-sizing: border-box;
    }

    header {
        display: grid;
        grid-template-columns: 100px 1fr 100px;
        align-items: center;
        margin-bottom: 1.25rem;
    }

    .clear-button {
        justify-self: end;
        box-shadow: none;
    }

    .tabs-container {
        display: flex;
        justify-content: center;
        margin-bottom: 1rem;
    }

    .segmented-control {
        position: relative;
        display: flex;
        width: 100%;
        background: var(--bg-card);
        padding: 3px;
        border-radius: var(--radius-md);
        border: 1px solid var(--border);
        box-shadow: inset 0 1px 2px rgba(0, 0, 0, 0.02);
    }

    .indicator {
        position: absolute;
        top: 3px;
        left: 3px;
        width: calc(50% - 3px);
        height: calc(100% - 6px);
        background: var(--accent);
        border-radius: 9px;
        transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1);
        z-index: 1;
        box-shadow: 0 4px 12px rgba(100, 108, 255, 0.25);
    }

    .indicator.app {
        transform: translateX(100%);
    }

    .segmented-control button {
        flex: 1;
        position: relative;
        z-index: 2;
        padding: 0.65rem;
        border: none;
        background: transparent;
        color: var(--text-muted);
        font-weight: 600;
        font-size: 0.85rem;
        cursor: pointer;
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 0.5rem;
        transition: color 0.3s ease;
    }

    .segmented-control button.active {
        color: white;
    }

    .segmented-control button:hover:not(.active) {
        color: var(--text-main);
    }

    main {
        flex: 1;
        background: var(--bg-card);
        border-radius: var(--radius-lg);
        border: 1px solid var(--border);
        overflow: hidden;
        display: flex;
        flex-direction: column;
        box-shadow: var(--shadow);
    }

    .log-view {
        flex: 1;
        overflow-y: auto;
        padding: 1.25rem;
        font-family: "JetBrains Mono", "Fira Code", "Menlo", monospace;
        font-size: 0.75rem;
        line-height: 1.6;
        color: var(--text-main);
        background: rgba(0, 0, 0, 0.02);
    }

    :global(.dark-theme) .log-view {
        background: rgba(0, 0, 0, 0.1);
    }

    pre {
        margin: 0;
        white-space: pre-wrap;
        word-break: break-all;
        opacity: 0.9;
    }

    .empty {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        height: 100%;
        color: var(--text-muted);
        text-align: center;
        padding: 2rem;
    }

    .empty p {
        margin: 0;
        margin-top: 1rem;
        font-style: italic;
        font-size: 0.9rem;
    }

    /* Scrollbar Styling */
    .log-view::-webkit-scrollbar {
        width: 6px;
    }

    .log-view::-webkit-scrollbar-track {
        background: transparent;
    }

    .log-view::-webkit-scrollbar-thumb {
        background: var(--border);
        border-radius: 10px;
    }

    .log-view::-webkit-scrollbar-thumb:hover {
        background: var(--text-muted);
    }
</style>
