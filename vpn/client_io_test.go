package vpn

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/fadinflame/simple-ac/config"
)

// recordingCLI saves its args and stdin, then prints out. Background `status`
// polls are answered without recording so they can't overwrite the files.
func recordingCLI(t *testing.T, out string) (path, argsFile, stdinFile string) {
	t.Helper()
	dir := t.TempDir()
	argsFile = filepath.Join(dir, "args")
	stdinFile = filepath.Join(dir, "stdin")
	body := fmt.Sprintf("[ \"$1\" = status ] && { echo %q; exit 0; }\necho \"$@\" > %q\ncat > %q\necho %q", out, argsFile, stdinFile, out)
	return writeScript(t, "vpn", body, 0o755), argsFile, stdinFile
}

func read(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func clientWith(creds *config.Credentials, cli string) *Client {
	return NewClient(creds, func() string { return cli })
}

func TestConnectPipesCredentialsWithoutOTP(t *testing.T) {
	cli, argsFile, stdinFile := recordingCLI(t, ">> state: Connected")
	c := clientWith(&config.Credentials{Server: "vpn.example.com", Username: "alice", Password: "pw", Group: 2}, cli)
	defer c.Stop()

	if ok, err := c.Connect(); err != nil || !ok {
		t.Fatalf("Connect() = %v, %v", ok, err)
	}

	if got := strings.TrimSpace(read(t, argsFile)); got != "-s connect vpn.example.com" {
		t.Errorf("args = %q", got)
	}
	if got, want := read(t, stdinFile), "2\nalice\npw\ny\n"; got != want {
		t.Errorf("stdin = %q, want %q", got, want)
	}
}

func TestConnectAddsTOTPCodeWhenSecretIsSet(t *testing.T) {
	cli, _, stdinFile := recordingCLI(t, ">> state: Connected")
	c := clientWith(&config.Credentials{
		Server: "s", Username: "alice", Password: "pw", Group: 0, OTPSecret: "JBSWY3DPEHPK3PXP",
	}, cli)
	defer c.Stop()

	if ok, err := c.Connect(); err != nil || !ok {
		t.Fatalf("Connect() = %v, %v", ok, err)
	}
	if !regexp.MustCompile(`^0\nalice\npw\n\d{6}\ny\n$`).MatchString(read(t, stdinFile)) {
		t.Errorf("stdin = %q, want group, user, password, 6-digit code, y", read(t, stdinFile))
	}
}

func TestConnectFailsOnInvalidOTPSecret(t *testing.T) {
	cli, _, _ := recordingCLI(t, ">> state: Connected")
	c := clientWith(&config.Credentials{Server: "s", Username: "u", Password: "p", OTPSecret: "not base32!"}, cli)
	defer c.Stop()

	if ok, err := c.Connect(); ok || err == nil {
		t.Errorf("Connect() = %v, %v; want failure before running the CLI", ok, err)
	}
}

func TestConnectNotConnectedWhenOutputLacksState(t *testing.T) {
	cli, _, _ := recordingCLI(t, "login failed")
	c := clientWith(&config.Credentials{Server: "s", Username: "u", Password: "p"}, cli)
	defer c.Stop()

	ok, err := c.Connect()
	if err != nil || ok {
		t.Errorf("Connect() = %v, %v; want false, nil", ok, err)
	}
	if connected, _ := c.IsConnected(); connected {
		t.Error("cache must stay disconnected")
	}
}

func TestDisconnectUpdatesCache(t *testing.T) {
	cli, argsFile, _ := recordingCLI(t, ">> state: Disconnected")
	c := clientWith(&config.Credentials{Server: "s", Username: "u", Password: "p"}, cli)
	defer c.Stop()

	c.mu.Lock()
	c.connected = true
	c.mu.Unlock()

	if ok, err := c.Disconnect(); err != nil || !ok {
		t.Fatalf("Disconnect() = %v, %v", ok, err)
	}
	if connected, _ := c.IsConnected(); connected {
		t.Error("cache should be disconnected")
	}
	if got := strings.TrimSpace(read(t, argsFile)); got != "disconnect" {
		t.Errorf("args = %q", got)
	}
}

func TestStatusChecksAreNotLogged(t *testing.T) {
	cli, _, _ := recordingCLI(t, ">> state: Connected")
	c := clientWith(&config.Credentials{Server: "s", Username: "u", Password: "p"}, cli)
	defer c.Stop()

	var streamed strings.Builder
	c.OnLog = func(s string) { streamed.WriteString(s) }

	if _, err := c.SyncStatus(); err != nil {
		t.Fatal(err)
	}
	if c.GetLogs() != "" || streamed.Len() != 0 {
		t.Errorf("status check leaked into logs: %q / %q", c.GetLogs(), streamed.String())
	}

	if _, err := c.Connect(); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(c.GetLogs(), ">> state: Connected") || !strings.Contains(streamed.String(), ">> state: Connected") {
		t.Errorf("connect output should be logged and streamed: %q / %q", c.GetLogs(), streamed.String())
	}
}

func TestLogBufferIsCapped(t *testing.T) {
	cli := writeScript(t, "vpn", "cat >/dev/null\nhead -c 400000 /dev/zero | tr '\\0' x\necho\necho '>> state: Connected'", 0o755)
	c := clientWith(&config.Credentials{Server: "s", Username: "u", Password: "p"}, cli)
	defer c.Stop()

	if ok, err := c.Connect(); err != nil || !ok {
		t.Fatalf("Connect() = %v, %v", ok, err)
	}
	logs := c.GetLogs()
	if len(logs) > maxLogSize {
		t.Errorf("log buffer = %d bytes, cap is %d", len(logs), maxLogSize)
	}
	if !strings.HasSuffix(strings.TrimSpace(logs), ">> state: Connected") {
		t.Error("the newest output must be kept when trimming")
	}
}

func TestGetCredentialsReturnsACopy(t *testing.T) {
	c := clientWith(&config.Credentials{Server: "s", Username: "u", Password: "p"}, "/nonexistent")
	defer c.Stop()

	got := c.GetCredentials()
	got.Password = "changed"
	if c.GetCredentials().Password != "p" {
		t.Error("mutating the returned credentials must not affect the client")
	}
}

func TestStopIsSafeAfterUse(t *testing.T) {
	c := clientWith(&config.Credentials{}, "/nonexistent")
	c.Stop()
}
