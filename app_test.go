package main

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/fadinflame/simple-ac/config"
)

func newTestApp(t *testing.T) *App {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("shell-script fixtures are unix-only")
	}
	t.Setenv("HOME", t.TempDir())
	return NewApp()
}

func writeCLI(t *testing.T, body string, mode os.FileMode) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "vpn")
	if err := os.WriteFile(path, []byte("#!/bin/sh\n"+body+"\n"), mode); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestSetCLIPathRejectsInvalidFiles(t *testing.T) {
	notCisco := writeCLI(t, "echo hello", 0o755)
	notExec := writeCLI(t, `echo ">> state: Disconnected"`, 0o644)

	tests := []struct {
		name string
		path string
		want string
	}{
		{"missing", filepath.Join(t.TempDir(), "nope"), "File not found."},
		{"directory", t.TempDir(), "Not a file."},
		{"not executable", notExec, "Not an executable."},
		{"not the VPN client", notCisco, "This doesn't look like the VPN command-line client."},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := newTestApp(t)
			err := a.SetCLIPath(tt.path)
			if err == nil || err.Error() != tt.want {
				t.Fatalf("SetCLIPath() error = %v, want %q", err, tt.want)
			}
			if got := a.ConfigManager.GetConfig().CLIPath; got != "" {
				t.Errorf("rejected path must not be saved, got %q", got)
			}
		})
	}
}

func TestSetCLIPathSavesAndResets(t *testing.T) {
	a := newTestApp(t)
	good := writeCLI(t, `echo ">> state: Disconnected"`, 0o755)

	if err := a.ConfigManager.SaveConfig(&config.Config{CredentialsFilePath: "/tmp/creds.enc"}); err != nil {
		t.Fatal(err)
	}

	if err := a.SetCLIPath("  " + good + "  "); err != nil {
		t.Fatalf("SetCLIPath(valid) = %v", err)
	}

	// Reload from disk to prove it was persisted and the other settings survived.
	reloaded := config.NewManager()
	if err := reloaded.LoadConfig(); err != nil {
		t.Fatal(err)
	}
	cfg := reloaded.GetConfig()
	if cfg.CLIPath != good {
		t.Errorf("CLIPath = %q, want %q", cfg.CLIPath, good)
	}
	if cfg.CredentialsFilePath != "/tmp/creds.enc" {
		t.Errorf("CredentialsFilePath was lost: %q", cfg.CredentialsFilePath)
	}

	if err := a.SetCLIPath(""); err != nil {
		t.Fatalf("SetCLIPath(\"\") = %v", err)
	}
	if got := a.ConfigManager.GetConfig().CLIPath; got != "" {
		t.Errorf("CLIPath after reset = %q, want empty", got)
	}
}

func TestSetCLIPathRejectionKeepsPreviousPath(t *testing.T) {
	a := newTestApp(t)
	good := writeCLI(t, `echo ">> state: Disconnected"`, 0o755)

	if err := a.SetCLIPath(good); err != nil {
		t.Fatal(err)
	}
	if err := a.SetCLIPath(writeCLI(t, "echo hello", 0o755)); err == nil {
		t.Fatal("expected an error for a non-VPN binary")
	}
	if got := a.ConfigManager.GetConfig().CLIPath; got != good {
		t.Errorf("previous path should be kept, got %q", got)
	}
}
