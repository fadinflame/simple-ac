<script>
    import { vpnLogs, appLogs, currentScreen } from "../store.js";
    import { AppState } from "../constants.js";
    import { LogView, Segmented, Toolbar, ToolbarButton } from "../lib/ui";

    let activeTab = "vpn"; // 'vpn' or 'app'

    const tabs = [
        { value: "vpn", label: "VPN CLI", icon: "terminal" },
        { value: "app", label: "App", icon: "app" },
    ];

    $: text = activeTab === "vpn" ? $vpnLogs : $appLogs;

    function goBack() {
        $currentScreen = AppState.MAIN;
    }

    function clearLogs() {
        if (activeTab === "vpn") vpnLogs.set("");
        else appLogs.set("");
    }
</script>

<svelte:window on:keydown={(e) => e.key === "Escape" && goBack()} />

<div class="ac-window ac-window--native">
    <Toolbar>
        <ToolbarButton slot="left" icon="chevron-left" label="Back" on:click={goBack} />
        <Segmented slot="center" small label="Log source" items={tabs} bind:value={activeTab} />
        <ToolbarButton slot="right" icon="trash" label="Clear Log" disabled={!text} on:click={clearLogs} />
    </Toolbar>

    <div class="ac-content" style="padding: var(--space-1) var(--space-5) var(--space-5)">
        {#key activeTab}
            <LogView
                {text}
                emptyText={activeTab === "vpn" ? "No VPN output yet" : "No app events yet"}
            />
        {/key}
    </div>
</div>
