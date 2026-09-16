<script>
    import { onMount } from "svelte";
    import { fade, fly } from "svelte/transition";
    import {
        GetCredentials,
        UpdateCredentials,
    } from "../../wailsjs/go/main/App.js";
    import { currentScreen } from "../store.js";
    import { AppState } from "../constants.js";

    let isLoading = true;
    let isSaving = false;
    let loadError = "";
    let saveError = "";
    let saveSuccess = false;

    let creds = {
        server: "",
        username: "",
        password: "",
        group: 0,
        otp_secret: "",
    };

    let masterPassword = "";
    let showPassword = false;
    let showOtpSecret = false;

    let errors = {
        server: false,
        username: false,
        password: false,
        masterPassword: false,
    };

    onMount(async () => {
        try {
            const current = await GetCredentials();
            if (current) creds = { ...current };
        } catch (err) {
            loadError = "Failed to load current credentials";
            console.error(err);
        } finally {
            isLoading = false;
        }
    });

    function goBack() {
        $currentScreen = AppState.SETTINGS;
    }

    async function handleSave() {
        errors.server = !creds.server;
        errors.username = !creds.username;
        errors.password = !creds.password;
        errors.masterPassword = !masterPassword;
        saveError = "";

        if (errors.server || errors.username || errors.password || errors.masterPassword)
            return;

        isSaving = true;
        try {
            await UpdateCredentials(creds, masterPassword);
            saveSuccess = true;
            setTimeout(() => {
                $currentScreen = AppState.SETTINGS;
            }, 600);
        } catch (err) {
            saveError = "Invalid master password or save failed";
            console.error("Failed to update credentials:", err);
        } finally {
            isSaving = false;
        }
    }
</script>

<div
    class="edit-creds-screen"
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
        <h2>Edit Credentials</h2>
        <div class="header-spacer"></div>
    </header>

    <main>
        {#if isLoading}
            <p class="status-text">Loading...</p>
        {:else if loadError}
            <p class="status-text error-text">{loadError}</p>
        {:else}
        <div class="card">
            <div class="row">
                <div class="form-group">
                    <label for="server">Server (Domain)</label>
                    <input
                        id="server"
                        class:invalid={errors.server}
                        bind:value={creds.server}
                        on:input={() => (errors.server = false)}
                        placeholder="vpn.example.com"
                    />
                </div>

                <div class="form-group">
                    <label for="group">
                        Group
                        <div class="tooltip-container">
                            <svg
                                xmlns="http://www.w3.org/2000/svg"
                                width="14"
                                height="14"
                                viewBox="0 0 24 24"
                                fill="none"
                                stroke="currentColor"
                                stroke-width="2"
                                stroke-linecap="round"
                                stroke-linejoin="round"
                                class="info-icon"
                            >
                                <circle cx="12" cy="12" r="10"></circle><line
                                    x1="12"
                                    y1="16"
                                    x2="12"
                                    y2="12"
                                ></line><line
                                    x1="12"
                                    y1="8"
                                    x2="12.01"
                                    y2="8"
                                ></line>
                            </svg>
                            <span class="tooltip-text"
                                >Group index (0 = Default)</span
                            >
                        </div>
                    </label>
                    <input
                        id="group"
                        type="number"
                        bind:value={creds.group}
                        placeholder="0"
                    />
                </div>
            </div>

            <div class="form-group full-width">
                <label for="user">Username</label>
                <input
                    id="user"
                    class:invalid={errors.username}
                    bind:value={creds.username}
                    on:input={() => (errors.username = false)}
                    placeholder="username"
                />
            </div>

            <div class="form-group full-width">
                <label for="pass">Password</label>
                <div class="field-wrapper">
                    <input
                        id="pass"
                        type={showPassword ? "text" : "password"}
                        class:invalid={errors.password}
                        value={creds.password}
                        on:input={(e) => {
                            creds.password = e.target.value;
                            errors.password = false;
                        }}
                        placeholder="••••••••"
                    />
                    <button
                        type="button"
                        class="reveal-btn"
                        tabindex="-1"
                        aria-label={showPassword ? "Hide password" : "Show password"}
                        on:click={() => (showPassword = !showPassword)}
                    >
                        {#if showPassword}
                            <svg
                                xmlns="http://www.w3.org/2000/svg"
                                width="16"
                                height="16"
                                viewBox="0 0 24 24"
                                fill="none"
                                stroke="currentColor"
                                stroke-width="2"
                                stroke-linecap="round"
                                stroke-linejoin="round"
                            >
                                <path
                                    d="M17.94 17.94A10.94 10.94 0 0 1 12 20c-7 0-11-8-11-8a18.5 18.5 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"
                                ></path>
                                <line x1="1" y1="1" x2="23" y2="23"></line>
                            </svg>
                        {:else}
                            <svg
                                xmlns="http://www.w3.org/2000/svg"
                                width="16"
                                height="16"
                                viewBox="0 0 24 24"
                                fill="none"
                                stroke="currentColor"
                                stroke-width="2"
                                stroke-linecap="round"
                                stroke-linejoin="round"
                            >
                                <path
                                    d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"
                                ></path>
                                <circle cx="12" cy="12" r="3"></circle>
                            </svg>
                        {/if}
                    </button>
                </div>
            </div>

            <div class="form-group full-width">
                <label for="otp"
                    >OTP Secret <span class="not-needed">(optional)</span
                    ></label
                >
                <div class="field-wrapper">
                    <input
                        id="otp"
                        type={showOtpSecret ? "text" : "password"}
                        value={creds.otp_secret}
                        on:input={(e) => (creds.otp_secret = e.target.value)}
                        placeholder="JBSWY3DPEHPK3PXP"
                    />
                    <button
                        type="button"
                        class="reveal-btn"
                        tabindex="-1"
                        aria-label={showOtpSecret ? "Hide secret" : "Show secret"}
                        on:click={() => (showOtpSecret = !showOtpSecret)}
                    >
                        {#if showOtpSecret}
                            <svg
                                xmlns="http://www.w3.org/2000/svg"
                                width="16"
                                height="16"
                                viewBox="0 0 24 24"
                                fill="none"
                                stroke="currentColor"
                                stroke-width="2"
                                stroke-linecap="round"
                                stroke-linejoin="round"
                            >
                                <path
                                    d="M17.94 17.94A10.94 10.94 0 0 1 12 20c-7 0-11-8-11-8a18.5 18.5 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"
                                ></path>
                                <line x1="1" y1="1" x2="23" y2="23"></line>
                            </svg>
                        {:else}
                            <svg
                                xmlns="http://www.w3.org/2000/svg"
                                width="16"
                                height="16"
                                viewBox="0 0 24 24"
                                fill="none"
                                stroke="currentColor"
                                stroke-width="2"
                                stroke-linecap="round"
                                stroke-linejoin="round"
                            >
                                <path
                                    d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"
                                ></path>
                                <circle cx="12" cy="12" r="3"></circle>
                            </svg>
                        {/if}
                    </button>
                </div>
            </div>

            <div class="divider"></div>

            <div class="form-group full-width">
                <label for="master-pass">Master Password</label>
                <input
                    id="master-pass"
                    type="password"
                    class:invalid={errors.masterPassword}
                    bind:value={masterPassword}
                    on:input={() => (errors.masterPassword = false)}
                    placeholder="Required to save changes"
                />
                <span class="hint-text"
                    >Confirms your identity before re-encrypting the file</span
                >
                {#if saveError}
                    <span class="error-text">{saveError}</span>
                {/if}
            </div>
        </div>
        {/if}
    </main>

    {#if !isLoading && !loadError}
        <div class="actions">
            <button type="button" class="btn btn-secondary" on:click={goBack}
                >Cancel</button
            >
            <button
                type="button"
                class="btn btn-primary"
                disabled={isSaving}
                on:click={handleSave}
            >
                {isSaving ? "Saving..." : saveSuccess ? "Saved!" : "Save Changes"}
            </button>
        </div>
    {/if}
</div>

<style>
    .edit-creds-screen {
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
        margin-bottom: 1rem;
        flex-shrink: 0;
    }

    .header-spacer {
        width: 32px;
    }

    main {
        flex: 1;
        display: flex;
        flex-direction: column;
        overflow-y: auto;
        -ms-overflow-style: none;
        scrollbar-width: none;
    }

    main::-webkit-scrollbar {
        display: none;
    }

    .card {
        flex: 1;
        display: flex;
        flex-direction: column;
        background: var(--bg-card);
        padding: 1.1rem 1.25rem;
        border-radius: var(--radius-lg);
        box-shadow: var(--shadow);
        color: var(--text-main);
    }

    .status-text {
        text-align: center;
        color: var(--text-muted);
        margin-top: 2rem;
    }

    .row {
        display: flex;
        gap: 1rem;
    }

    .form-group {
        display: flex;
        flex-direction: column;
        gap: 0.35rem;
        margin-bottom: 0.85rem;
        flex: 1;
    }

    .full-width {
        width: 100%;
    }

    label {
        font-size: 0.8rem;
        font-weight: 600;
        color: var(--text-main);
        display: flex;
        align-items: center;
        gap: 0.4rem;
    }

    .not-needed {
        color: var(--text-muted);
        font-weight: 400;
    }

    input {
        padding: 0.65rem 0.75rem;
        border: 1px solid var(--border);
        border-radius: var(--radius-sm);
        font-size: 0.9rem;
        transition: all 0.2s;
        width: 100%;
        box-sizing: border-box;
        background: var(--input-bg);
        color: var(--text-main);
    }

    input.invalid {
        border-color: var(--danger);
        background-color: var(--input-bg);
        animation: shake 0.2s ease-in-out 0s 2;
    }

    input:focus {
        outline: none;
        border-color: var(--accent);
    }

    .field-wrapper {
        position: relative;
        display: flex;
        align-items: center;
    }

    .field-wrapper input {
        padding-right: 2.5rem;
    }

    .reveal-btn {
        position: absolute;
        right: 0.5rem;
        display: flex;
        align-items: center;
        justify-content: center;
        background: none;
        border: none;
        padding: 0.3rem;
        color: var(--text-muted);
        cursor: pointer;
        opacity: 0;
        transition: opacity 0.15s ease-out, color 0.15s;
    }

    .field-wrapper:hover .reveal-btn,
    .field-wrapper:focus-within .reveal-btn {
        opacity: 1;
    }

    .reveal-btn:hover {
        color: var(--accent);
    }

    .tooltip-container {
        position: relative;
        display: inline-flex;
        align-items: center;
        cursor: help;
    }

    .info-icon {
        color: var(--text-muted);
        transition: color 0.2s;
    }

    .tooltip-container:hover .info-icon {
        color: var(--accent);
    }

    .tooltip-text {
        visibility: hidden;
        width: 180px;
        background-color: var(--tooltip-bg);
        color: var(--tooltip-text);
        text-align: center;
        border-radius: 6px;
        padding: 6px 10px;
        position: absolute;
        z-index: 50;
        top: 125%;
        left: 50%;
        transform: translateX(-50%);
        opacity: 0;
        transition: opacity 0.1s ease-out;
        font-size: 0.75rem;
        font-weight: 400;
        line-height: 1.2;
        pointer-events: none;
        box-shadow: 0 4px 6px -1px rgb(0 0 0 / 0.1);
    }

    .tooltip-text::after {
        content: "";
        position: absolute;
        bottom: 100%;
        left: 50%;
        margin-left: -5px;
        border-width: 5px;
        border-style: solid;
        border-color: transparent transparent var(--tooltip-bg) transparent;
    }

    .tooltip-container:hover .tooltip-text {
        visibility: visible;
        opacity: 1;
    }

    .divider {
        height: 1px;
        background: var(--border);
        margin: 0.35rem 0 1rem;
    }

    .hint-text {
        font-size: 0.72rem;
        color: var(--text-muted);
    }

    .error-text {
        color: var(--danger);
        font-size: 0.75rem;
        margin-top: 0.2rem;
        font-weight: 500;
    }

    .actions {
        display: flex;
        gap: 1rem;
        margin-top: 1rem;
        flex-shrink: 0;
    }

    .actions .btn {
        flex: 1;
    }

    @keyframes shake {
        0%,
        100% {
            transform: translateX(0);
        }
        25% {
            transform: translateX(4px);
        }
        75% {
            transform: translateX(-4px);
        }
    }
</style>
