package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/fadinflame/simple-ac/config"
	"github.com/fadinflame/simple-ac/vpn"
)

var testCreds = config.Credentials{Server: "vpn.example.com", Username: "alice", Password: "pw", Group: 1}

// newUnlockableApp returns an App whose VPN client runs a fake CLI, never a real one.
func newUnlockableApp(t *testing.T, cliOutput string) *App {
	t.Helper()
	a := newTestApp(t)
	cli := writeCLI(t, "cat >/dev/null\necho '"+cliOutput+"'", 0o755)
	if err := a.ConfigManager.SaveConfig(&config.Config{CLIPath: cli}); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = a.Lock() })
	return a
}

func savedApp(t *testing.T) (*App, string) {
	t.Helper()
	a := newUnlockableApp(t, ">> state: Disconnected")
	path := filepath.Join(t.TempDir(), "creds.enc")
	if err := a.SaveConfig(testCreds, "master", path); err != nil {
		t.Fatal(err)
	}
	return a, path
}

func TestGetDefaultCLIPath(t *testing.T) {
	if got := NewApp().GetDefaultCLIPath(); got != vpn.DefaultCLIPath() {
		t.Errorf("GetDefaultCLIPath() = %q", got)
	}
}

func TestSaveConfigKeepsCLIPathAndLoadsClient(t *testing.T) {
	a, path := savedApp(t)

	cfg := a.GetConfig()
	if cfg.CredentialsFilePath != path {
		t.Errorf("CredentialsFilePath = %q, want %q", cfg.CredentialsFilePath, path)
	}
	if cfg.CLIPath == "" {
		t.Error("setup must not drop the custom CLI path")
	}
	if !a.CheckConfig() {
		t.Error("CheckConfig() should be true after setup")
	}

	creds, err := a.GetCredentials()
	if err != nil || *creds != testCreds {
		t.Errorf("GetCredentials() = %+v, %v", creds, err)
	}
}

func TestCheckConfig(t *testing.T) {
	t.Run("no config", func(t *testing.T) {
		if newTestApp(t).CheckConfig() {
			t.Error("want false")
		}
	})
	t.Run("credentials file missing", func(t *testing.T) {
		a, path := savedApp(t)
		if err := os.Remove(path); err != nil {
			t.Fatal(err)
		}
		if a.CheckConfig() {
			t.Error("a moved credentials file should send the user back to setup")
		}
	})
}

func TestUnlock(t *testing.T) {
	a, _ := savedApp(t)
	if err := a.Lock(); err != nil {
		t.Fatal(err)
	}

	if ok, err := a.Unlock("wrong"); ok || err == nil {
		t.Errorf("Unlock(wrong) = %v, %v", ok, err)
	}
	if _, err := a.GetCredentials(); err == nil {
		t.Error("a failed unlock must not load credentials")
	}

	if ok, err := a.Unlock("master"); !ok || err != nil {
		t.Fatalf("Unlock(master) = %v, %v", ok, err)
	}
	creds, err := a.GetCredentials()
	if err != nil || *creds != testCreds {
		t.Errorf("GetCredentials() = %+v, %v", creds, err)
	}
}

func TestUnlockWithoutConfig(t *testing.T) {
	if ok, err := newTestApp(t).Unlock("x"); ok || err == nil {
		t.Errorf("Unlock() = %v, %v; want an error", ok, err)
	}
}

func TestLockClearsCredentials(t *testing.T) {
	a, _ := savedApp(t)
	if err := a.Lock(); err != nil {
		t.Fatal(err)
	}
	if _, err := a.GetCredentials(); err == nil {
		t.Error("credentials must be gone after Lock")
	}
	if connected, _ := a.IsConnected(); connected {
		t.Error("locked app must report disconnected")
	}
	if a.GetVPNLogs() != "" {
		t.Error("locked app must not expose logs")
	}
}

func TestVPNActionsRequireUnlock(t *testing.T) {
	a := newTestApp(t)
	if _, err := a.Connect(); err == nil {
		t.Error("Connect() should fail when locked")
	}
	if _, err := a.Disconnect(); err == nil {
		t.Error("Disconnect() should fail when locked")
	}
	if a.IsConnecting() {
		t.Error("a rejected request must not leave the app in the connecting state")
	}
}

func TestConnectUsesConfiguredCLI(t *testing.T) {
	a := newUnlockableApp(t, ">> state: Connected")
	if err := a.SaveConfig(testCreds, "master", filepath.Join(t.TempDir(), "creds.enc")); err != nil {
		t.Fatal(err)
	}

	ok, err := a.Connect()
	if err != nil || !ok {
		t.Fatalf("Connect() = %v, %v", ok, err)
	}
	if a.IsConnecting() {
		t.Error("connecting flag must be cleared when Connect returns")
	}
	if connected, _ := a.IsConnected(); !connected {
		t.Error("expected connected state")
	}
}

func TestUpdateCredentials(t *testing.T) {
	updated := config.Credentials{Server: "other.example.com", Username: "bob", Password: "new", Group: 4}

	t.Run("invalid credentials", func(t *testing.T) {
		a, _ := savedApp(t)
		bad := updated
		bad.Server = ""
		if err := a.UpdateCredentials(bad, "master"); err == nil {
			t.Error("want a validation error")
		}
	})

	t.Run("wrong master password", func(t *testing.T) {
		a, path := savedApp(t)
		if err := a.UpdateCredentials(updated, "wrong"); err == nil || err.Error() != "invalid master password" {
			t.Errorf("error = %v", err)
		}
		onDisk, err := a.ConfigManager.LoadCredentials("master", path)
		if err != nil || *onDisk != testCreds {
			t.Errorf("file must be unchanged, got %+v, %v", onDisk, err)
		}
		if live, _ := a.GetCredentials(); *live != testCreds {
			t.Errorf("live credentials must be unchanged, got %+v", live)
		}
	})

	t.Run("success", func(t *testing.T) {
		a, path := savedApp(t)
		if err := a.UpdateCredentials(updated, "master"); err != nil {
			t.Fatal(err)
		}
		onDisk, err := a.ConfigManager.LoadCredentials("master", path)
		if err != nil || *onDisk != updated {
			t.Errorf("file = %+v, %v", onDisk, err)
		}
		if live, _ := a.GetCredentials(); *live != updated {
			t.Errorf("live credentials = %+v", live)
		}
	})

	t.Run("no credentials file configured", func(t *testing.T) {
		a := newUnlockableApp(t, ">> state: Disconnected")
		if err := a.UpdateCredentials(updated, "master"); err == nil {
			t.Error("want an error")
		}
	})
}
