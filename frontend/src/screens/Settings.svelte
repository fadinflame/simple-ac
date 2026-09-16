<script>
    import { fade, fly, scale } from "svelte/transition";
    import { cubicOut } from "svelte/easing";
    import {
        BrowserOpenURL,
    } from "../../wailsjs/runtime/runtime.js";
    import {
        CheckForUpdates,
    } from "../../wailsjs/go/main/App.js";
    import {
        currentScreen,
        theme,
        appVersion,
        hasCheckedUpdate,
        updateRelease,
    } from "../store.js";
    import { AppState } from "../constants.js";

    let showResetConfirm = false;
    let isCheckingUpdate = false;
    let updateStatus = "";

    function goBack() {
        $currentScreen = AppState.MAIN;
    }

    function toggleTheme() {
        theme.update((t) => (t === "light" ? "dark" : "light"));
    }


    async function checkUpdate(manual = false) {
        isCheckingUpdate = true;
        try {
            const release = await CheckForUpdates();
            $hasCheckedUpdate = true;
            if (release) {
                $updateRelease = release;
                BrowserOpenURL(release.html_url);
            } else if (manual) {
                updateStatus = "Up to date";
                setTimeout(() => (updateStatus = ""), 2000);
            }
        } catch (err) {
            if (manual) {
                updateStatus = "Error";
                setTimeout(() => (updateStatus = ""), 3000);
            }
        } finally {
            isCheckingUpdate = false;
        }
    }

    function handleReset() {
        showResetConfirm = true;
    }

    function confirmReset() {
        showResetConfirm = false;
        $currentScreen = AppState.SETUP;
    }
</script>

<div
    class="settings-container"
    in:fly={{ x: 20, duration: 400 }}
    out:fade={{ duration: 200 }}
>
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
        <h2>Settings</h2>
        <div class="header-spacer"></div>
    </header>

    <main>
        <div class="settings-list">
            <section>
                <h2 class="section-title">Preference</h2>
                <div class="settings-item">
                    <div class="settings-info">
                        <span class="settings-label">Dark Mode</span>
                        <span class="settings-desc"
                            >Switch between light and dark themes</span
                        >
                    </div>
                    <button
                        class="toggle-switch"
                        class:active={$theme === "dark"}
                        on:click={toggleTheme}
                        aria-label="Toggle Dark Mode"
                    >
                        <div class="toggle-knob"></div>
                    </button>
                </div>

            </section>

            <section>
                <h2 class="section-title">Account</h2>
                <div class="settings-item">
                    <div class="settings-info">
                        <span class="settings-label">VPN Credentials</span>
                        <span class="settings-desc"
                            >Update server, login and OTP secret</span
                        >
                    </div>
                    <button
                        class="btn btn-secondary btn-sm"
                        on:click={() =>
                            ($currentScreen = AppState.EDIT_CREDENTIALS)}
                    >
                        Edit
                    </button>
                </div>
            </section>

            <section>
                <h2 class="section-title">About</h2>
                <div class="settings-item">
                    <div class="settings-info">
                        <span class="settings-label">Application Version</span>
                        <span class="settings-desc"
                            >Currently installed version</span
                        >
                    </div>
                    <div class="version-container">
                        <span class="version-tag"
                            >{$appVersion || "v1.0.0"}</span
                        >
                        <button
                            class="btn btn-primary btn-sm"
                            class:tooltip-visible={isCheckingUpdate ||
                                updateStatus}
                            data-tooltip={isCheckingUpdate
                                ? "Checking..."
                                : updateStatus || undefined}
                            on:click={() => checkUpdate(true)}
                            disabled={isCheckingUpdate}
                        >
                            Check
                        </button>
                    </div>
                </div>
            </section>

            <section>
                <h2 class="section-title danger">Danger Zone</h2>
                <div class="settings-item danger-item">
                    <div class="settings-info">
                        <span class="settings-label">Reset Application</span>
                        <span class="settings-desc"
                            >Delete all configurations and lock</span
                        >
                    </div>
                    <button
                        class="btn btn-danger btn-sm"
                        on:click={handleReset}
                    >
                        Reset
                    </button>
                </div>
            </section>
        </div>
    </main>

    {#if showResetConfirm}
        <div class="modal-overlay" transition:fade={{ duration: 200 }}>
            <div
                class="modal"
                in:scale={{ start: 0.9, duration: 200, easing: cubicOut }}
            >
                <h3>Reset Settings?</h3>
                <p>
                    All saved VPN data will be deleted. You will need to
                    reconfigure the application.
                </p>
                <div class="modal-actions">
                    <button
                        class="btn btn-secondary"
                        on:click={() => (showResetConfirm = false)}
                        >Cancel</button
                    >
                    <button
                        class="btn btn-primary danger-action"
                        on:click={confirmReset}>Reset All</button
                    >
                </div>
            </div>
        </div>
    {/if}
</div>

<style>
    .settings-container {
        display: flex;
        flex-direction: column;
        height: 100vh;
        background: var(--bg-app);
        padding: 1rem 1.5rem;
        box-sizing: border-box;
    }

    header {
        display: grid;
        grid-template-columns: 32px 1fr 32px;
        align-items: center;
        margin-bottom: 1.25rem;
    }

    .header-spacer {
        width: 32px;
    }

    main {
        flex: 1;
        overflow-y: auto;
        /* Hide scrollbar for IE, Edge and Firefox */
        -ms-overflow-style: none;
        scrollbar-width: none;
    }

    /* Hide scrollbar for Chrome, Safari and Opera */
    main::-webkit-scrollbar {
        display: none;
    }

    .settings-list {
        display: flex;
        flex-direction: column;
        gap: 1.25rem;
        max-width: 600px;
        margin: 0 auto;
        width: 100%;
    }

    .section-title {
        font-size: 0.7rem;
        font-weight: 700;
        text-transform: uppercase;
        letter-spacing: 0.05em;
        color: var(--text-muted);
        margin-bottom: 0.5rem;
        padding-left: 0.25rem;
    }

    .section-title.danger {
        color: var(--danger);
    }

    .settings-item {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 0.75rem 1rem;
        background: var(--bg-card);
        border-radius: var(--radius-md);
        border: 1px solid var(--border);
        margin-bottom: 0.5rem;
        transition: all 0.2s;
        text-align: left;
    }

    .settings-info {
        display: flex;
        flex-direction: column;
        gap: 0.25rem;
    }

    .settings-label {
        font-weight: 600;
        font-size: 0.85rem;
        color: var(--text-main);
    }

    .settings-desc {
        font-size: 0.7rem;
        color: var(--text-muted);
    }

    .toggle-switch {
        width: 40px;
        height: 22px;
        background: var(--border);
        border-radius: 11px;
        position: relative;
        cursor: pointer;
        transition: all 0.3s;
        border: none;
        padding: 0;
    }

    .toggle-switch.active {
        background: var(--accent);
    }

    .toggle-knob {
        position: absolute;
        width: 16px;
        height: 16px;
        background: white;
        border-radius: 50%;
        top: 3px;
        left: 3px;
        transition: all 0.3s cubic-bezier(0.68, -0.55, 0.265, 1.55);
        box-shadow: 0 1px 3px rgba(0, 0, 0, 0.2);
    }

    .active .toggle-knob {
        left: calc(100% - 19px);
    }

    .version-container {
        display: flex;
        align-items: center;
        gap: 0.75rem;
    }

    .version-tag {
        background: var(--bg-app);
        color: var(--text-muted);
        padding: 4px 10px;
        border-radius: var(--radius-sm);
        font-size: 0.8rem;
        font-weight: 700;
        font-family: "JetBrains Mono", monospace;
    }

    .danger-item {
        border-left: 4px solid var(--danger);
    }

    .btn-sm {
        padding: 4px 12px;
        font-size: 0.75rem;
        border-radius: 8px;
    }

    .danger-action {
        background: var(--danger) !important;
    }

    .danger-action:hover {
        background: var(--danger) !important;
        opacity: 0.9;
    }

    /* Modal styles */
    .modal-overlay {
        position: fixed;
        inset: 0;
        background: rgba(0, 0, 0, 0.6);
        backdrop-filter: blur(8px);
        display: flex;
        align-items: center;
        justify-content: center;
        z-index: 1000;
        padding: 1.5rem;
    }

    .modal {
        background: var(--bg-card);
        width: 100%;
        max-width: 380px;
        border-radius: 24px;
        padding: 2rem;
        box-shadow: 0 20px 40px rgba(0, 0, 0, 0.2);
        text-align: center;
    }

    .modal h3 {
        margin: 0 0 1rem;
        font-size: 1.25rem;
        color: var(--text-main);
    }

    .modal p {
        margin: 0 0 2rem;
        color: var(--text-muted);
        line-height: 1.5;
    }

    .modal-actions {
        display: flex;
        gap: 1rem;
    }

    .modal-actions .btn {
        flex: 1;
    }
</style>
