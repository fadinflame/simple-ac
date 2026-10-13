package vpn

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/fadinflame/simple-ac/config"
)

func writeScript(t *testing.T, name, body string, mode os.FileMode) string {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("shell-script fixtures are unix-only")
	}
	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, []byte("#!/bin/sh\n"+body+"\n"), mode); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestDefaultCLIPath(t *testing.T) {
	want := unixCLI
	if runtime.GOOS == "windows" {
		want = windowsCLI
	}
	if got := DefaultCLIPath(); got != want {
		t.Errorf("DefaultCLIPath() = %q, want %q", got, want)
	}
}

func TestVerifyCLI(t *testing.T) {
	ciscoLike := writeScript(t, "vpn", `echo ">> state: Disconnected"`, 0o755)
	ciscoExitsNonZero := writeScript(t, "vpn", "echo '>> state: Disconnected'\nexit 1", 0o755)
	foreign := writeScript(t, "tool", `echo hello`, 0o755)
	silent := writeScript(t, "tool", `exit 0`, 0o755)
	notExecutable := writeScript(t, "vpn", `echo ">> state: Disconnected"`, 0o644)

	tests := []struct {
		name    string
		path    string
		wantErr string // substring; empty means success
	}{
		{"cisco-like output", ciscoLike, ""},
		{"cisco-like output with non-zero exit", ciscoExitsNonZero, ""},
		{"foreign executable", foreign, "doesn't look like"},
		{"executable with no output", silent, "doesn't look like"},
		{"missing file", filepath.Join(t.TempDir(), "nope"), "not found"},
		{"not executable", notExecutable, "No permission"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := VerifyCLI(tt.path)
			if tt.wantErr == "" {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}
			var cliErr *CLIError
			if !errors.As(err, &cliErr) {
				t.Fatalf("want *CLIError, got %T (%v)", err, err)
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("error %q does not contain %q", err, tt.wantErr)
			}
		})
	}
}

func TestUserError(t *testing.T) {
	cliErr := &CLIError{"boom"}
	if got := userError(cliErr, "fallback"); got != error(cliErr) {
		t.Errorf("CLIError should pass through, got %v", got)
	}

	wrapped := errors.Join(errors.New("ctx"), cliErr)
	var target *CLIError
	if got := userError(wrapped, "fallback"); !errors.As(got, &target) {
		t.Errorf("wrapped CLIError should be kept, got %v", got)
	}

	if got := userError(errors.New("exit status 1"), "fallback"); got.Error() != "fallback" {
		t.Errorf("other errors should use the fallback, got %q", got)
	}
}

func newTestClient(cliPath func() string) *Client {
	c := NewClient(&config.Credentials{Server: "vpn.example.com", Username: "u", Password: "p"}, cliPath)
	return c
}

func TestClientUsesCustomCLIPath(t *testing.T) {
	script := writeScript(t, "vpn", `echo ">> state: Connected"`, 0o755)
	c := newTestClient(func() string { return script })
	defer c.Stop()

	connected, err := c.SyncStatus()
	if err != nil {
		t.Fatal(err)
	}
	if !connected {
		t.Error("expected the custom CLI to be used and report Connected")
	}
}

func TestClientReadsCLIPathOnEveryCommand(t *testing.T) {
	connectedCLI := writeScript(t, "vpn", `echo ">> state: Connected"`, 0o755)
	disconnectedCLI := writeScript(t, "vpn", `echo ">> state: Disconnected"`, 0o755)

	current := connectedCLI
	c := newTestClient(func() string { return current })
	defer c.Stop()

	if ok, _ := c.SyncStatus(); !ok {
		t.Fatal("expected Connected from the first CLI")
	}
	current = disconnectedCLI
	if ok, _ := c.SyncStatus(); ok {
		t.Error("expected Disconnected after the path changed")
	}
}

func TestConnectWithMissingCLI(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "nope")
	c := newTestClient(func() string { return missing })
	defer c.Stop()

	var logged strings.Builder
	c.OnLog = func(s string) { logged.WriteString(s) }

	ok, err := c.Connect()
	if ok {
		t.Error("Connect should fail")
	}
	var cliErr *CLIError
	if !errors.As(err, &cliErr) || !strings.Contains(err.Error(), "not found") {
		t.Errorf("want a 'not found' CLIError, got %v", err)
	}
	if !strings.Contains(logged.String(), "[ERROR]") || !strings.Contains(logged.String(), missing) {
		t.Errorf("start failure should be logged with the path, got %q", logged.String())
	}
}

func TestConnectAndDisconnectWithFailingCLI(t *testing.T) {
	script := writeScript(t, "tool", "echo nope\nexit 1", 0o755)
	c := newTestClient(func() string { return script })
	defer c.Stop()

	if _, err := c.Connect(); err == nil || err.Error() != "Connection failed. See Logs for details." {
		t.Errorf("Connect error = %v", err)
	}
	if _, err := c.Disconnect(); err == nil || err.Error() != "Disconnect failed. See Logs for details." {
		t.Errorf("Disconnect error = %v", err)
	}
}

func TestConnectSucceedsWhenCLIReportsConnected(t *testing.T) {
	script := writeScript(t, "vpn", `cat >/dev/null; echo ">> state: Connected"`, 0o755)
	c := newTestClient(func() string { return script })
	defer c.Stop()

	ok, err := c.Connect()
	if err != nil || !ok {
		t.Fatalf("Connect() = %v, %v", ok, err)
	}
	if connected, _ := c.IsConnected(); !connected {
		t.Error("cached state should be connected")
	}
}
