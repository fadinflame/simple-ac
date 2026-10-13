<script>
    import { fade } from "svelte/transition";
    import { BrowserOpenURL } from "../../wailsjs/runtime/runtime.js";
    import { CheckForUpdates, GetCredentials } from "../../wailsjs/go/main/App.js";
    import { onMount } from "svelte";
    import {
        currentScreen,
        theme,
        glass,
        appVersion,
        hasCheckedUpdate,
        updateRelease,
    } from "../store.js";
    import { AppState } from "../constants.js";
    import {
        Button,
        LinkButton,
        Segmented,
        Group,
        Row,
        Toolbar,
        ToolbarButton,
        Icon,
    } from "../lib/ui";

    let showResetConfirm = false;
    let isCheckingUpdate = false;
    let updateStatus = "";
    let account = "";

    const glassModes = [
        { value: "clear", label: "Light" },
        { value: "tinted", label: "Window" },
    ];

    const themes = [
        { value: "system", label: "System" },
        { value: "light", label: "Light" },
        { value: "dark", label: "Dark" },
    ];

    onMount(async () => {
        try {
            const c = await GetCredentials();
            if (c) account = [c.server, c.username].filter(Boolean).join(" · ");
        } catch {
            account = "";
        }
    });

    function goBack() {
        $currentScreen = AppState.MAIN;
    }

    async function checkUpdate() {
        isCheckingUpdate = true;
        updateStatus = "";
        try {
            const release = await CheckForUpdates();
            $hasCheckedUpdate = true;
            if (release) {
                $updateRelease = release;
                updateStatus = `${release.tag_name} is available`;
                BrowserOpenURL(release.html_url);
            } else {
                updateStatus = "SimpleAC is up to date";
            }
        } catch (err) {
            updateStatus = "Couldn't check for updates";
        } finally {
            isCheckingUpdate = false;
        }
    }

    function confirmReset() {
        showResetConfirm = false;
        $currentScreen = AppState.SETUP;
    }

    function onKey(e) {
        if (e.key !== "Escape") return;
        if (showResetConfirm) showResetConfirm = false;
        else goBack();
    }
</script>

<svelte:window on:keydown={onKey} />

<div class="ac-window ac-window--native">
    <Toolbar title="Settings">
        <ToolbarButton slot="left" round icon="chevron-left" label="Back" on:click={goBack} />
    </Toolbar>

    <div class="ac-content ac-content--scroll">
        <Group title="Appearance">
            <Row label="Theme">
                <Segmented small label="Theme" items={themes} bind:value={$theme} />
            </Row>
            <Row label="Toolbar">
                <Segmented small label="Toolbar" items={glassModes} bind:value={$glass} />
            </Row>
        </Group>

        <Group title="Account">
            <Row label="VPN credentials" sub={account || "Server, login and OTP secret"}>
                <Button size="small" on:click={() => ($currentScreen = AppState.EDIT_CREDENTIALS)}>Edit…</Button>
            </Row>
        </Group>

        <Group title="About">
            <Row label="Version" sub={updateStatus ? `${$appVersion || "dev"} — ${updateStatus}` : $appVersion || "dev"}>
                <Button size="small" disabled={isCheckingUpdate} on:click={checkUpdate}>
                    {isCheckingUpdate ? "Checking…" : "Check for Updates…"}
                </Button>
            </Row>
        </Group>

        <Group footer="Runs the first-time setup again. Saving to the same file replaces the current configuration.">
            <Row>
                <LinkButton slot="label" tone="danger" on:click={() => (showResetConfirm = true)}>
                    Set Up Again…
                </LinkButton>
            </Row>
        </Group>
    </div>

    {#if showResetConfirm}
        <div class="ac-overlay" transition:fade={{ duration: 120 }}>
            <div class="ac-alert" role="alertdialog" aria-labelledby="reset-title" aria-describedby="reset-text">
                <div class="ac-unlock-icon"><Icon name="alert" size={24} stroke={1.4} /></div>
                <h3 class="ac-alert-title" id="reset-title">Set up SimpleAC again?</h3>
                <p class="ac-alert-text" id="reset-text">
                    You'll enter your VPN details and choose a master password again.
                </p>
                <div class="ac-alert-actions">
                    <Button variant="danger" block on:click={confirmReset}>Set Up Again</Button>
                    <Button block on:click={() => (showResetConfirm = false)}>Cancel</Button>
                </div>
            </div>
        </div>
    {/if}
</div>
