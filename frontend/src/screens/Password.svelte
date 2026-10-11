<script>
    import { Unlock } from "../../wailsjs/go/main/App.js";
    import { navigate } from "../store.js";
    import { AppState } from "../constants.js";
    import { Button, LinkButton, TextField, Icon, Toolbar } from "../lib/ui";

    let password = "";
    let isUnlocking = false;
    let error = "";

    async function handleUnlock() {
        if (!password || isUnlocking) return;

        isUnlocking = true;
        error = "";

        try {
            const success = await Unlock(password);
            if (success) {
                navigate(AppState.MAIN);
            } else {
                error = "Wrong password. Try again.";
                password = "";
            }
        } catch (err) {
            console.error("Unlock error:", err);
            error = "Couldn't read the configuration file.";
            password = "";
        } finally {
            isUnlocking = false;
        }
    }

    function goToSetup() {
        navigate(AppState.SETUP);
    }

    // Clear the error as soon as the user types again.
    $: if (password) error = "";
</script>

<div class="ac-window ac-window--native">
    <Toolbar />
    <div class="ac-content">
        <div class="ac-unlock">
            <div class="ac-unlock-icon">
                <Icon name="lock" size={28} stroke={1.4} />
            </div>
            <h1 class="ac-status-title">SimpleAC</h1>
            <p class="ac-status-sub">Enter your master password to unlock.</p>

            <form on:submit|preventDefault={handleUnlock}>
                <TextField
                    type="password"
                    placeholder="Master password"
                    bind:value={password}
                    {error}
                    autofocus
                />
                <Button type="submit" variant="default" block disabled={!password || isUnlocking}>
                    {isUnlocking ? "Unlocking…" : "Unlock"}
                </Button>
            </form>
        </div>

        <p class="ac-group-footer" style="text-align: center">
            Forgot your password?
            <LinkButton small on:click={goToSetup}>Set Up Again…</LinkButton>
        </p>
    </div>
</div>
