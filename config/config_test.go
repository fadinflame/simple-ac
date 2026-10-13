package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// newManager points the manager's config dir at a temp dir.
func newManager(t *testing.T) *Manager {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Setenv("APPDATA", t.TempDir())
	} else {
		t.Setenv("HOME", t.TempDir())
	}
	return NewManager()
}

func TestCredentialsValidate(t *testing.T) {
	valid := Credentials{Server: "s", Username: "u", Password: "p"}
	if err := valid.Validate(); err != nil {
		t.Errorf("valid credentials rejected: %v", err)
	}

	tests := []struct {
		name   string
		mutate func(*Credentials)
		want   string
	}{
		{"no server", func(c *Credentials) { c.Server = "" }, "server"},
		{"no username", func(c *Credentials) { c.Username = "" }, "username"},
		{"no password", func(c *Credentials) { c.Password = "" }, "password"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := valid
			tt.mutate(&c)
			err := c.Validate()
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Errorf("Validate() = %v, want error mentioning %q", err, tt.want)
			}
		})
	}
}

func TestCredentialsOptionalFields(t *testing.T) {
	c := Credentials{Server: "s", Username: "u", Password: "p"}
	if err := c.Validate(); err != nil {
		t.Errorf("group and OTP secret must be optional: %v", err)
	}
}

func TestConfigValidate(t *testing.T) {
	if err := (&Config{}).Validate(); err == nil {
		t.Error("empty credentials path must be invalid")
	}
	if err := (&Config{CredentialsFilePath: "/x"}).Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestConfigJSONOmitsEmptyCLIPath(t *testing.T) {
	data, err := json.Marshal(Config{CredentialsFilePath: "/x"})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data), "cli_path") {
		t.Errorf("empty cli_path should be omitted, got %s", data)
	}

	data, _ = json.Marshal(Config{CredentialsFilePath: "/x", CLIPath: "/y"})
	if !strings.Contains(string(data), `"cli_path":"/y"`) {
		t.Errorf("cli_path missing from %s", data)
	}
}

func TestManagerConfigRoundTrip(t *testing.T) {
	m := newManager(t)
	if m.ConfigExists() {
		t.Fatal("fresh manager should have no config")
	}

	want := &Config{CredentialsFilePath: "/creds.enc", CLIPath: "/opt/vpn"}
	if err := m.SaveConfig(want); err != nil {
		t.Fatal(err)
	}
	if !m.ConfigExists() {
		t.Error("config should exist after saving")
	}

	fresh := NewManager()
	if err := fresh.LoadConfig(); err != nil {
		t.Fatal(err)
	}
	if got := *fresh.GetConfig(); got != *want {
		t.Errorf("loaded %+v, want %+v", got, *want)
	}
}

func TestManagerConfigPermissions(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("unix permission bits")
	}
	m := newManager(t)
	if err := m.SaveConfig(&Config{CredentialsFilePath: "/x"}); err != nil {
		t.Fatal(err)
	}

	file, err := os.Stat(m.ConfigPath())
	if err != nil {
		t.Fatal(err)
	}
	if perm := file.Mode().Perm(); perm != 0o600 {
		t.Errorf("config file mode = %o, want 600", perm)
	}
	dir, _ := os.Stat(filepath.Dir(m.ConfigPath()))
	if perm := dir.Mode().Perm(); perm != 0o700 {
		t.Errorf("config dir mode = %o, want 700", perm)
	}
}

func TestGetConfigDefaultsToEmpty(t *testing.T) {
	m := newManager(t)
	cfg := m.GetConfig()
	if cfg == nil || cfg.CredentialsFilePath != "" || cfg.CLIPath != "" {
		t.Errorf("GetConfig() on a missing file = %+v, want empty config", cfg)
	}
}

func TestLoadConfigRejectsInvalidJSON(t *testing.T) {
	m := newManager(t)
	if err := os.MkdirAll(filepath.Dir(m.ConfigPath()), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(m.ConfigPath(), []byte("{not json"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := m.LoadConfig(); err == nil {
		t.Error("invalid JSON should be an error")
	}
}

func TestLoadConfigIgnoresUnknownAndOldFiles(t *testing.T) {
	m := newManager(t)
	if err := os.MkdirAll(filepath.Dir(m.ConfigPath()), 0o700); err != nil {
		t.Fatal(err)
	}
	// A config written before cli_path existed, with a stale extra key.
	old := `{"credentials_file_path":"/creds.enc","show_console_log":true}`
	if err := os.WriteFile(m.ConfigPath(), []byte(old), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := m.LoadConfig(); err != nil {
		t.Fatal(err)
	}
	cfg := m.GetConfig()
	if cfg.CredentialsFilePath != "/creds.enc" || cfg.CLIPath != "" {
		t.Errorf("got %+v", cfg)
	}
}

func TestManagerCredentialsRoundTrip(t *testing.T) {
	m := newManager(t)
	path := filepath.Join(t.TempDir(), "creds.enc")
	want := &Credentials{Server: "vpn.example.com", Username: "u", Password: "p", Group: 3, OTPSecret: "JBSWY3DPEHPK3PXP"}

	if err := m.SaveCredentials(want, "master", path); err != nil {
		t.Fatal(err)
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(raw), "vpn.example.com") || strings.Contains(string(raw), `"p"`) {
		t.Error("credentials file must not contain plaintext")
	}

	got, err := m.LoadCredentials("master", path)
	if err != nil {
		t.Fatal(err)
	}
	if *got != *want {
		t.Errorf("loaded %+v, want %+v", *got, *want)
	}
}

func TestLoadCredentialsErrors(t *testing.T) {
	m := newManager(t)
	path := filepath.Join(t.TempDir(), "creds.enc")
	if err := m.SaveCredentials(&Credentials{Server: "s", Username: "u", Password: "p"}, "master", path); err != nil {
		t.Fatal(err)
	}

	if _, err := m.LoadCredentials("wrong", path); err == nil {
		t.Error("wrong master password must fail")
	}
	if _, err := m.LoadCredentials("master", filepath.Join(t.TempDir(), "missing")); err == nil {
		t.Error("missing file must fail")
	}
}

func TestSaveCredentialsFileMode(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("unix permission bits")
	}
	m := newManager(t)
	path := filepath.Join(t.TempDir(), "creds.enc")
	if err := m.SaveCredentials(&Credentials{Server: "s", Username: "u", Password: "p"}, "pw", path); err != nil {
		t.Fatal(err)
	}
	info, _ := os.Stat(path)
	if perm := info.Mode().Perm(); perm != 0o600 {
		t.Errorf("credentials file mode = %o, want 600", perm)
	}
}
