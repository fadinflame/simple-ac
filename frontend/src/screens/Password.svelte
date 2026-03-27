<script>
    import { onMount } from "svelte";
    import { fade, fly, scale } from "svelte/transition";
    import { cubicOut } from "svelte/easing";
    import { Unlock } from "../../wailsjs/go/main/App.js";
    import { navigate } from "../store.js";
    import { AppState } from "../constants.js";

    let password = "";
    let showPassword = false;
    let passwordInput;
    let isUnlocking = false;
    let error = "";

    onMount(() => {
        if (passwordInput) {
            passwordInput.focus();
        }
    });

    async function handleUnlock() {
        if (!password) {
            error = "Please enter your password";
            return;
        }

        isUnlocking = true;
        error = "";

        try {
            const success = await Unlock(password);
            if (success) {
                setTimeout(() => {
                    navigate(AppState.MAIN);
                }, 500);
            } else {
                error = "Incorrect master password";
                password = "";
            }
        } catch (err) {
            console.error("Unlock error:", err);
            error = "Error verifying password";
            password = "";
        } finally {
            isUnlocking = false;
        }
    }

    function goToSetup() {
        navigate(AppState.SETUP);
    }
</script>

<div class="lock-container">
    <div
        class="lock-card"
        in:scale={{ start: 0.9, duration: 400, easing: cubicOut }}
    >
        <header>
            <div class="lock-icon">
                <svg
                    xmlns="http://www.w3.org/2000/svg"
                    width="32"
                    height="32"
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
            </div>
            <h1>SimpleAC</h1>
            <p>Enter your master password to access the configuration</p>
        </header>

        <form on:submit|preventDefault={handleUnlock}>
            <div class="form-group">
                <input
                    type="password"
                    bind:value={password}
                    class:invalid={error}
                    bind:this={passwordInput}
                    placeholder="Master password"
                    on:input={() => (error = "")}
                />
                {#if error}
                    <span class="error-text" transition:fade>{error}</span>
                {/if}
            </div>

            <div class="actions">
                <button
                    type="submit"
                    class="btn btn-primary"
                    disabled={isUnlocking}
                >
                    {isUnlocking ? "Loading..." : "Unlock"}
                </button>
            </div>
        </form>

        <footer>
            <button type="button" class="btn-text" on:click={goToSetup}>
                Forgot password or want to reconfigure?
            </button>
        </footer>
    </div>
</div>

<style>
    .lock-container {
        display: flex;
        justify-content: center;
        align-items: center;
        height: 100vh;
        background: var(--bg-app);
        padding: 2rem;
        box-sizing: border-box;
    }

    .lock-card {
        background: var(--bg-card);
        padding: 2.5rem 2rem;
        border-radius: var(--radius-lg);
        box-shadow: var(--shadow-lg);
        width: 100%;
        max-width: 380px;
        text-align: center;
        color: var(--text-main);
    }

    header {
        margin-bottom: 2rem;
    }

    .lock-icon {
        width: 64px;
        height: 64px;
        background: var(--bg-app);
        color: var(--accent);
        border-radius: var(--radius-md);
        display: flex;
        align-items: center;
        justify-content: center;
        margin: 0 auto 1.5rem;
    }

    h1 {
        margin: 0;
        font-size: 1.6rem;
        font-weight: 800;
        color: var(--text-main);
    }

    p {
        margin: 0.5rem 0 0;
        color: var(--text-muted);
        font-size: 0.9rem;
        line-height: 1.5;
    }

    form {
        margin-top: 2rem;
    }

    .form-group {
        margin-bottom: 1.5rem;
        text-align: left;
    }

    input {
        width: 100%;
        padding: 0.85rem 1rem;
        border: 1px solid var(--border);
        border-radius: var(--radius-md);
        font-size: 1rem;
        transition: all 0.2s;
        box-sizing: border-box;
        background: var(--input-bg);
        color: var(--text-main);
    }

    input:focus {
        outline: none;
        border-color: var(--accent);
        box-shadow: 0 0 0 4px rgba(100, 108, 255, 0.1);
    }

    input.invalid {
        border-color: var(--danger);
        background: var(--input-bg);
        animation: shake 0.4s ease-in-out;
    }

    .error-text {
        color: var(--danger);
        font-size: 0.8rem;
        margin-top: 0.5rem;
        display: block;
        font-weight: 600;
    }

    .actions .btn {
        width: 100%;
        padding: 0.9rem;
        font-size: 1rem;
    }

    footer {
        margin-top: 2rem;
        border-top: 1px solid var(--border);
        padding-top: 1.5rem;
    }

    @keyframes shake {
        0%,
        100% {
            transform: translateX(0);
        }
        25% {
            transform: translateX(8px);
        }
        75% {
            transform: translateX(-8px);
        }
    }
</style>
