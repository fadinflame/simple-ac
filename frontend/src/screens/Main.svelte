<script>
    import { onMount, onDestroy } from "svelte";
    import { WindowSetTitle, BrowserOpenURL } from "../../wailsjs/runtime/runtime.js";
    import { Connect, Disconnect, Lock, GetCredentials } from "../../wailsjs/go/main/App.js";
    import {
        currentScreen,
        refreshConfig,
        isConnected,
        updateRelease,
        connectionSeconds,
        isConnecting,
        checkConnectionStatus,
    } from "../store.js";
    import { AppState } from "../constants.js";
    import {
        Button,
        Icon,
        Toolbar,
        ToolbarButton,
        StatusHeader,
        Group,
        Row,
    } from "../lib/ui";

    let connectionError = false;
    let errorTimeout;
    // Only the non-secret fields are kept for display.
    let profile = { server: "", username: "", group: 0 };

    onMount(async () => {
        WindowSetTitle("SimpleAC");
        await refreshConfig();
        try {
            const c = await GetCredentials();
            if (c) profile = { server: c.server, username: c.username, group: c.group };
        } catch (err) {
            console.error("Failed to read profile:", err);
        }
    });

    onDestroy(() => {
        if (errorTimeout) clearTimeout(errorTimeout);
    });

    $: formattedTime = (() => {
        const h = Math.floor($connectionSeconds / 3600).toString().padStart(2, "0");
        const m = Math.floor(($connectionSeconds % 3600) / 60).toString().padStart(2, "0");
        const s = ($connectionSeconds % 60).toString().padStart(2, "0");
        return `${h}:${m}:${s}`;
    })();

    $: status = $isConnecting
        ? $isConnected
            ? "disconnecting"
            : "connecting"
        : $isConnected
          ? "connected"
          : "disconnected";

    $: subtitle = connectionError
        ? "Couldn't connect. See Logs for details."
        : status === "connecting"
          ? "Waiting for the server"
          : profile.server;

    async function toggleConnection() {
        if ($isConnecting) return;
        try {
            if (!$isConnected) {
                $isConnecting = true;
                connectionError = false;
                const result = await Connect();
                if (!result) showError();
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
        errorTimeout = setTimeout(() => (connectionError = false), 8000);
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

    function onKey(e) {
        if (!e.metaKey && !e.ctrlKey) return;
        if (e.key === ",") { e.preventDefault(); goToSettings(); }
        if (e.key.toLowerCase() === "l") { e.preventDefault(); lockApp(); }
    }
</script>

<svelte:window on:keydown={onKey} />

<div class="ac-window ac-window--native">
    <Toolbar title="SimpleAC">
        <svelte:fragment slot="right">
            <ToolbarButton icon="sliders" label="Settings" on:click={goToSettings} />
            <ToolbarButton icon="doc" label="Logs" on:click={goToLogs} />
            <ToolbarButton icon="lock" label="Lock" on:click={lockApp} />
        </svelte:fragment>
    </Toolbar>

    <div class="ac-content" style="padding-top: var(--space-2)">
        {#if $updateRelease}
            <div class="ac-notice">
                <span class="ac-notice-icon"><Icon name="arrow-down-circle" /></span>
                <span class="ac-notice-text">
                    <b>SimpleAC {$updateRelease.tag_name}</b> is available.
                </span>
                <Button size="small" on:click={() => BrowserOpenURL($updateRelease.html_url)}>View…</Button>
                <button
                    class="ac-tbtn"
                    title="Dismiss"
                    aria-label="Dismiss"
                    on:click={() => ($updateRelease = null)}
                >
                    <Icon name="xmark" size={14} />
                </button>
            </div>
        {/if}

        <div style="margin-top: 36px">
            <StatusHeader state={status} {subtitle} error={connectionError} timer={formattedTime} />
        </div>

        <div style="display: flex; justify-content: center; margin-top: var(--space-5)">
            {#if status === "connected"}
                <Button wide on:click={toggleConnection}>Disconnect</Button>
            {:else if status === "disconnected"}
                <Button variant="default" wide on:click={toggleConnection}>Connect</Button>
            {:else}
                <Button wide disabled>{status === "connecting" ? "Connecting…" : "Disconnecting…"}</Button>
            {/if}
        </div>

        <div class="ac-spacer"></div>

        <Group>
            <Row label="Server" value={profile.server} />
            <Row label="Username" value={profile.username} />
            <Row label="Group" value={String(profile.group ?? 0)} />
        </Group>
    </div>
</div>
