<script>
    import { onMount } from "svelte";
    import { GetCredentials, UpdateCredentials } from "../../wailsjs/go/main/App.js";
    import { currentScreen } from "../store.js";
    import { AppState } from "../constants.js";
    import {
        Button,
        TextField,
        NumberField,
        Group,
        FieldRow,
        Toolbar,
        ToolbarButton,
    } from "../lib/ui";

    let isLoading = true;
    let isSaving = false;
    let loadError = "";
    let saveError = "";

    let creds = {
        server: "",
        username: "",
        password: "",
        group: 0,
        otp_secret: "",
    };

    let masterPassword = "";

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
            loadError = "Couldn't load the current credentials.";
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

        if (errors.server || errors.username || errors.password || errors.masterPassword) return;

        isSaving = true;
        try {
            await UpdateCredentials(creds, masterPassword);
            goBack();
        } catch (err) {
            saveError = "Wrong master password, or the file couldn't be saved.";
            console.error("Failed to update credentials:", err);
        } finally {
            isSaving = false;
        }
    }
</script>

<svelte:window on:keydown={(e) => e.key === "Escape" && goBack()} />

<div class="ac-window ac-window--native">
    <Toolbar title="VPN Credentials">
        <ToolbarButton slot="left" round icon="chevron-left" label="Back" on:click={goBack} />
    </Toolbar>

    {#if isLoading}
        <div class="ac-content"><div class="ac-log-empty">Loading…</div></div>
    {:else if loadError}
        <div class="ac-content"><div class="ac-log-empty" style="color: var(--danger)">{loadError}</div></div>
    {:else}
        <form id="edit-creds" class="ac-content ac-content--scroll" on:submit|preventDefault={handleSave}>
            <Group footer="Group is the number of the connection profile in AnyConnect's list, usually 0.">
                <FieldRow label="Server" forId="server">
                    <TextField
                        id="server"
                        placeholder="vpn.example.com"
                        bind:value={creds.server}
                        error={errors.server}
                        on:input={() => (errors.server = false)}
                    />
                </FieldRow>
                <FieldRow label="Group" forId="group">
                    <NumberField id="group" bind:value={creds.group} min={0} max={99} />
                </FieldRow>
                <FieldRow label="Username" forId="user">
                    <TextField
                        id="user"
                        placeholder="name"
                        bind:value={creds.username}
                        error={errors.username}
                        on:input={() => (errors.username = false)}
                    />
                </FieldRow>
                <FieldRow label="Password" forId="pass">
                    <TextField
                        id="pass"
                        type="password"
                        placeholder="Required"
                        bind:value={creds.password}
                        error={errors.password}
                        on:input={() => (errors.password = false)}
                        revealable
                    />
                </FieldRow>
                <FieldRow label="OTP secret" forId="otp">
                    <TextField
                        id="otp"
                        type="password"
                        mono
                        placeholder="Optional"
                        bind:value={creds.otp_secret}
                        revealable
                    />
                </FieldRow>
            </Group>

            <Group title="Confirm with master password" footer="Needed to re-encrypt the file with your changes.">
                <FieldRow label="Password" forId="master-pass">
                    <TextField
                        id="master-pass"
                        type="password"
                        placeholder="Master password"
                        bind:value={masterPassword}
                        error={saveError || errors.masterPassword}
                        on:input={() => {
                            errors.masterPassword = false;
                            saveError = "";
                        }}
                    />
                </FieldRow>
            </Group>
        </form>

        <footer class="ac-assist-foot">
            <span class="ac-assist-step"></span>
            <Button on:click={goBack}>Cancel</Button>
            <Button type="submit" form="edit-creds" variant="default" disabled={isSaving}>
                {isSaving ? "Saving…" : "Save"}
            </Button>
        </footer>
    {/if}
</div>
