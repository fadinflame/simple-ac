<script>
    import {
        SelectConfigPath,
        SaveConfig,
        CheckConfig,
    } from "../../wailsjs/go/main/App.js";
    import { navigate, refreshConfig } from "../store.js";
    import { onMount } from "svelte";
    import { fly } from "svelte/transition";
    import { cubicOut } from "svelte/easing";
    import { AppState } from "../constants.js";

    let currentStep = 1;
    let isSaving = false;
    let configExists = false;

    let creds = {
        server: "",
        username: "",
        password: "",
        group: 0,
        otp_secret: "",
    };

    let configPath = "";
    let encryptionPassword = "";
    let confirmEncryptionPassword = "";

    let errors = {
        server: false,
        username: false,
        password: false,
        configPath: false,
        encryptionPassword: false,
        confirmEncryptionPassword: false,
    };

    onMount(async () => {
        configExists = await CheckConfig();
    });

    async function handleBrowse() {
        const path = await SelectConfigPath();
        if (path) {
            configPath = path;
            errors.configPath = false;
        }
    }

    function nextStep() {
        if (currentStep === 1) {
            errors.server = !creds.server;
            errors.username = !creds.username;
            errors.password = !creds.password;

            if (errors.server || errors.username || errors.password) return;
            currentStep = 2;
        }
    }

    function prevStep() {
        currentStep = 1;
    }

    async function handleSave() {
        errors.configPath = !configPath;
        errors.encryptionPassword = !encryptionPassword;
        errors.confirmEncryptionPassword = !confirmEncryptionPassword;

        if (
            encryptionPassword &&
            confirmEncryptionPassword &&
            encryptionPassword !== confirmEncryptionPassword
        ) {
            errors.confirmEncryptionPassword = true;
        }

        if (
            errors.configPath ||
            errors.encryptionPassword ||
            errors.confirmEncryptionPassword
        )
            return;

        isSaving = true;

        try {
            await SaveConfig(creds, encryptionPassword, configPath);
            await refreshConfig(); // Update global store
            setTimeout(() => {
                navigate(AppState.MAIN);
            }, 500);
        } catch (err) {
            console.error("Failed to save config:", err);
        } finally {
            isSaving = false;
        }
    }
</script>

<div class="setup-container">
    <div class="card">
        <header>
            <div class="step-indicator">
                <span class="step" class:active={currentStep === 1}>1</span>
                <div class="line"></div>
                <span class="step" class:active={currentStep === 2}>2</span>
            </div>
            <h1>{currentStep === 1 ? "VPN Setup" : "Security"}</h1>
            <p>
                {currentStep === 1
                    ? "Enter your VPN credentials"
                    : "Choose a storage location and protect your data"}
            </p>
        </header>

        <div class="step-content">
            {#if currentStep === 1}
                <div
                    class="step-panel"
                    in:fly={{ x: -20, duration: 250, easing: cubicOut }}
                    out:fly={{ x: 20, duration: 250, easing: cubicOut }}
                >
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
                                        <circle cx="12" cy="12" r="10"
                                        ></circle><line
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
                        <label for="pass">VPN Password</label>
                        <input
                            id="pass"
                            type="password"
                            class:invalid={errors.password}
                            bind:value={creds.password}
                            on:input={() => (errors.password = false)}
                            placeholder="••••••••"
                        />
                    </div>

                    <div class="form-group full-width">
                        <label for="otp"
                            >OTP Secret <span class="not-needed"
                                >(optional)</span
                            ></label
                        >
                        <input
                            id="otp"
                            bind:value={creds.otp_secret}
                            placeholder="JBSWY3DPEHPK3PXP"
                        />
                    </div>
                </div>
            {:else if currentStep === 2}
                <div
                    class="step-panel"
                    in:fly={{ x: 20, duration: 250, easing: cubicOut }}
                    out:fly={{ x: -20, duration: 250, easing: cubicOut }}
                >
                    <div class="form-group full-width">
                        <label for="path">Config File Path</label>
                        <div class="input-with-button">
                            <input
                                id="path"
                                class:invalid={errors.configPath}
                                readonly
                                placeholder="Choose a path..."
                                bind:value={configPath}
                            />
                            <button
                                type="button"
                                class="btn btn-secondary browse-btn"
                                on:click={handleBrowse}>Browse</button
                            >
                        </div>
                    </div>

                    <div class="encryption-section">
                        <div class="form-group full-width">
                            <label for="enc-pass">Master Password</label>
                            <input
                                id="enc-pass"
                                type="password"
                                class:invalid={errors.encryptionPassword}
                                bind:value={encryptionPassword}
                                on:input={() =>
                                    (errors.encryptionPassword = false)}
                                placeholder="Master password"
                            />
                        </div>

                        <div class="form-group full-width">
                            <label for="confirm-enc-pass"
                                >Confirm Password</label
                            >
                            <input
                                id="confirm-enc-pass"
                                type="password"
                                class:invalid={errors.confirmEncryptionPassword}
                                bind:value={confirmEncryptionPassword}
                                on:input={() =>
                                    (errors.confirmEncryptionPassword = false)}
                                placeholder="Repeat password"
                            />
                            {#if errors.confirmEncryptionPassword && encryptionPassword !== confirmEncryptionPassword && confirmEncryptionPassword !== ""}
                                <span class="error-text"
                                    >Passwords do not match</span
                                >
                            {/if}
                        </div>
                    </div>
                </div>
            {/if}
        </div>

        <div class="actions">
            {#if currentStep === 1}
                {#if configExists}
                    <button
                        type="button"
                        class="btn btn-secondary"
                        on:click={() => navigate(AppState.PASSWORD)}
                        >Cancel</button
                    >
                {/if}
                <button type="button" class="btn btn-primary" on:click={nextStep}
                    >Next</button
                >
            {:else}
                <button type="button" class="btn btn-secondary" on:click={prevStep}
                    >Back</button
                >
                <button
                    type="button"
                    class="btn btn-primary"
                    disabled={isSaving}
                    on:click={handleSave}
                >
                    {isSaving ? "Saving..." : "Encrypt & Save"}
                </button>
            {/if}
        </div>
    </div>
</div>

<style>
    .setup-container {
        display: flex;
        justify-content: center;
        align-items: center;
        min-height: 100vh;
        background: var(--bg-app);
        padding: 0 1.5rem;
        box-sizing: border-box;
    }

    .card {
        background: var(--bg-card);
        padding: 1.5rem;
        border-radius: var(--radius-lg);
        box-shadow: var(--shadow);
        width: 100%;
        max-width: 380px;
        color: var(--text-main);
    }

    header {
        margin-bottom: 1rem;
        text-align: left;
    }

    .step-indicator {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        margin-bottom: 1.5rem;
    }

    .step {
        width: 24px;
        height: 24px;
        border-radius: 50%;
        background: var(--border);
        color: var(--text-muted);
        display: flex;
        align-items: center;
        justify-content: center;
        font-size: 0.75rem;
        font-weight: 700;
        transition: all 0.3s;
    }

    .step.active {
        background: var(--accent);
        color: white;
    }

    .line {
        flex: 1;
        max-width: 40px;
        height: 2px;
        background: var(--border);
    }

    h1 {
        margin: 0;
        font-size: 1.4rem;
        color: var(--text-main);
    }

    p {
        margin: 0.4rem 0 0;
        color: var(--text-muted);
        font-size: 0.85rem;
        line-height: 1.4;
    }

    .step-content {
        position: relative;
        min-height: 350px;
    }

    .step-panel {
        position: absolute;
        top: 0;
        left: 0;
        width: 100%;
    }

    .form-group {
        display: flex;
        flex-direction: column;
        gap: 0.4rem;
        margin-bottom: 1rem;
    }

    .full-width {
        width: 100%;
    }

    .row {
        display: flex;
        gap: 1rem;
    }

    label {
        font-size: 0.8rem;
        font-weight: 600;
        color: var(--text-main);
        display: flex;
        align-items: center;
        gap: 0.4rem;
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

    .not-needed {
        color: var(--text-muted);
        font-weight: 400;
    }

    input {
        padding: 0.75rem;
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

    input[readonly] {
        background: var(--bg-app);
        cursor: default;
    }

    .input-with-button {
        display: flex;
        gap: 0.5rem;
    }

    .actions {
        display: flex;
        gap: 1rem;
        margin-top: 1rem;
    }

    .actions .btn {
        flex: 1;
    }

    .error-text {
        color: var(--danger);
        font-size: 0.75rem;
        margin-top: 0.2rem;
        font-weight: 500;
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
