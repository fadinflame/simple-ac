<script>
    import { SelectConfigPath, SaveConfig, CheckConfig } from "../../wailsjs/go/main/App.js";
    import { navigate, refreshConfig } from "../store.js";
    import { onMount } from "svelte";
    import { AppState } from "../constants.js";
    import {
        Button,
        TextField,
        NumberField,
        PathField,
        Group,
        FieldRow,
        Toolbar,
    } from "../lib/ui";

    let currentStep = 1;
    let isSaving = false;
    let configExists = false;
    let saveError = "";

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
        confirmEncryptionPassword: "",
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
        errors.server = !creds.server.trim();
        errors.username = !creds.username.trim();
        errors.password = !creds.password;
        if (errors.server || errors.username || errors.password) return;
        currentStep = 2;
    }

    function prevStep() {
        currentStep = 1;
    }

    async function handleSave() {
        errors.configPath = !configPath;
        errors.encryptionPassword = !encryptionPassword;
        errors.confirmEncryptionPassword = !confirmEncryptionPassword
            ? true
            : encryptionPassword !== confirmEncryptionPassword
              ? "Passwords don't match."
              : "";

        if (errors.configPath || errors.encryptionPassword || errors.confirmEncryptionPassword) return;

        isSaving = true;
        saveError = "";
        try {
            await SaveConfig(creds, encryptionPassword, configPath);
            await refreshConfig();
            navigate(AppState.MAIN);
        } catch (err) {
            console.error("Failed to save config:", err);
            saveError = "Couldn't save the file. Choose another location and try again.";
        } finally {
            isSaving = false;
        }
    }

    function onKey(e) {
        if (e.key !== "Escape") return;
        if (currentStep === 2) prevStep();
        else if (configExists) navigate(AppState.PASSWORD);
    }
</script>

<svelte:window on:keydown={onKey} />

<div class="ac-window ac-window--native">
    <Toolbar />

    {#if currentStep === 1}
        <form id="setup-step-1" class="ac-content" style="overflow-y: auto" on:submit|preventDefault={nextStep}>
            <div class="ac-assist-head">
                <h1 class="ac-assist-title">Connect to your VPN</h1>
                <p class="ac-assist-sub">SimpleAC passes these to the AnyConnect command-line client.</p>
            </div>

            <Group footer="Group is the number of the connection profile in AnyConnect's list, usually 0.">
                <FieldRow label="Server" forId="server">
                    <TextField
                        id="server"
                        placeholder="vpn.example.com"
                        bind:value={creds.server}
                        error={errors.server}
                        on:input={() => (errors.server = false)}
                        autofocus
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
            </Group>

            <Group footer="Optional. The Base32 key from your two-factor setup; SimpleAC generates the codes.">
                <FieldRow label="OTP secret" forId="otp">
                    <TextField id="otp" mono placeholder="JBSWY3DPEHPK3PXP" bind:value={creds.otp_secret} />
                </FieldRow>
            </Group>
        </form>

        <footer class="ac-assist-foot">
            <span class="ac-assist-step">Step 1 of 2</span>
            {#if configExists}
                <Button on:click={() => navigate(AppState.PASSWORD)}>Cancel</Button>
            {/if}
            <Button type="submit" form="setup-step-1" variant="default">Continue</Button>
        </footer>
    {:else}
        <form id="setup-step-2" class="ac-content" style="overflow-y: auto" on:submit|preventDefault={handleSave}>
            <div class="ac-assist-head">
                <h1 class="ac-assist-title">Protect your settings</h1>
                <p class="ac-assist-sub">Your credentials are encrypted with a master password before they're saved.</p>
            </div>

            <Group>
                <FieldRow label="Save to">
                    <PathField path={configPath} error={errors.configPath}>
                        <Button size="small" on:click={handleBrowse}>Choose…</Button>
                    </PathField>
                </FieldRow>
            </Group>

            <Group footer="There's no way to recover this password. If you forget it, you'll need to set up SimpleAC again.">
                <FieldRow label="Password" forId="enc-pass">
                    <TextField
                        id="enc-pass"
                        type="password"
                        placeholder="Required"
                        bind:value={encryptionPassword}
                        error={errors.encryptionPassword}
                        on:input={() => (errors.encryptionPassword = false)}
                        autofocus
                    />
                </FieldRow>
                <FieldRow label="Confirm" forId="confirm-enc-pass">
                    <TextField
                        id="confirm-enc-pass"
                        type="password"
                        placeholder="Repeat password"
                        bind:value={confirmEncryptionPassword}
                        error={errors.confirmEncryptionPassword}
                        on:input={() => (errors.confirmEncryptionPassword = "")}
                    />
                </FieldRow>
            </Group>

            {#if saveError}
                <p class="ac-group-footer" style="color: var(--danger)" role="alert">{saveError}</p>
            {/if}
        </form>

        <footer class="ac-assist-foot">
            <span class="ac-assist-step">Step 2 of 2</span>
            <Button on:click={prevStep}>Back</Button>
            <Button type="submit" form="setup-step-2" variant="default" disabled={isSaving}>
                {isSaving ? "Saving…" : "Encrypt & Save"}
            </Button>
        </footer>
    {/if}
</div>
