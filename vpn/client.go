package vpn

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/fadinflame/simple-ac/config"
	"github.com/pquerna/otp/totp"
)

const (
	unixCLI    = "/opt/cisco/anyconnect/bin/vpn"
	windowsCLI = `C:\Program Files (x86)\Cisco\Cisco AnyConnect Secure Mobility Client\vpncli.exe`
	maxLogSize = 256 * 1024 // 256KB limit

	statusTimeout     = 5 * time.Second
	connectTimeout    = 10 * time.Second
	disconnectTimeout = 5 * time.Second
)

// Client provides methods to interact with the Cisco AnyConnect CLI
type Client struct {
	mu        sync.RWMutex
	connected bool

	creds      *config.Credentials
	lastOutput string
	OnLog      func(string)

	cliPath       func() string
	stopReconcile chan struct{}
}

// DefaultCLIPath returns the standard Cisco CLI location for the current OS.
func DefaultCLIPath() string {
	if runtime.GOOS == "windows" {
		return windowsCLI
	}
	return unixCLI
}

// NewClient creates a client. cliPath returns a custom CLI path, or "" for the default; it may be nil.
func NewClient(credentials *config.Credentials, cliPath func() string) *Client {
	c := &Client{
		creds:         credentials,
		cliPath:       cliPath,
		stopReconcile: make(chan struct{}),
	}
	go c.reconcileLoop()
	return c
}

// IsConnected returns the cached connection state
func (c *Client) IsConnected() (bool, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.connected, nil
}

// GetCredentials returns a copy of the currently loaded credentials
func (c *Client) GetCredentials() *config.Credentials {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.creds == nil {
		return nil
	}
	credsCopy := *c.creds
	return &credsCopy
}

// updateStatus performs an actual CLI check and updates the cached state
func (c *Client) updateStatus() (bool, error) {
	output, err := c.executeCommand([]string{"status"}, "", true, statusTimeout)
	if err != nil {
		return false, fmt.Errorf("command failed: %w", err)
	}

	isConnected := strings.Contains(output, ">> state: Connected")
	c.mu.Lock()
	c.connected = isConnected
	c.mu.Unlock()

	return isConnected, nil
}

// SyncStatus forces a CLI check and updates state
func (c *Client) SyncStatus() (bool, error) {
	return c.updateStatus()
}

// reconcileLoop syncs the cached state every 30s
func (c *Client) reconcileLoop() {
	// Initial check
	_, _ = c.updateStatus()

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			_, _ = c.updateStatus()
		case <-c.stopReconcile:
			return
		}
	}
}

func (c *Client) Connect() (bool, error) {
	input := fmt.Sprintf("%d\n%s\n%s\n", c.creds.Group, c.creds.Username, c.creds.Password)
	if c.creds.OTPSecret != "" {
		code, err := totp.GenerateCode(c.creds.OTPSecret, time.Now())
		if err != nil {
			return false, err
		}
		input += fmt.Sprintf("%s\n", code)
	}
	input += "y\n"

	output, err := c.executeCommand([]string{"-s", "connect", c.creds.Server}, input, false, connectTimeout)
	if err != nil && !strings.Contains(output, ">> state: Connected") {
		return false, userError(err, "Connection failed. See Logs for details.")
	}

	ok := strings.Contains(output, ">> state: Connected")
	c.mu.Lock()
	c.connected = ok
	c.mu.Unlock()

	return ok, nil
}

func (c *Client) Disconnect() (bool, error) {
	output, err := c.executeCommand([]string{"disconnect"}, "", false, disconnectTimeout)
	if err != nil && !strings.Contains(output, ">> state: Disconnected") {
		return false, userError(err, "Disconnect failed. See Logs for details.")
	}

	ok := strings.Contains(output, ">> state: Disconnected")
	c.mu.Lock()
	c.connected = !ok // if disconnected correctly, ok is true, so connected becomes false
	c.mu.Unlock()

	return ok, nil
}

// GetLogs returns the current accumulated logs safely
func (c *Client) GetLogs() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.lastOutput
}

// Stop cleans up the background reconciliation goroutine
func (c *Client) Stop() {
	if c.stopReconcile != nil {
		close(c.stopReconcile)
	}
}

// CLIError means the CLI binary couldn't be started; its message is safe to show to the user.
type CLIError struct{ Msg string }

func (e *CLIError) Error() string { return e.Msg }

func startError(command string, err error) *CLIError {
	switch {
	case errors.Is(err, os.ErrNotExist):
		return &CLIError{fmt.Sprintf("VPN client not found at %s. Check the path in Settings.", command)}
	case errors.Is(err, os.ErrPermission):
		return &CLIError{fmt.Sprintf("No permission to run %s. Check that it is an executable.", command)}
	default:
		return &CLIError{fmt.Sprintf("Couldn't run %s: %v", command, err)}
	}
}

// userError keeps CLIError messages and replaces anything else with a generic one.
func userError(err error, fallback string) error {
	var cliErr *CLIError
	if errors.As(err, &cliErr) {
		return cliErr
	}
	return errors.New(fallback)
}

// VerifyCLI runs `status` on the binary and checks for the Cisco CLI's state output.
func VerifyCLI(path string) error {
	ctx, cancel := context.WithTimeout(context.Background(), statusTimeout)
	defer cancel()

	out, err := exec.CommandContext(ctx, path, "status").CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		return &CLIError{"The file didn't respond. It doesn't look like the VPN client."}
	}
	var exitErr *exec.ExitError
	if err != nil && !errors.As(err, &exitErr) {
		return startError(path, err)
	}
	if !strings.Contains(string(out), ">> state:") {
		return &CLIError{"This doesn't look like the VPN command-line client."}
	}
	return nil
}

func (c *Client) executeCommand(args []string, input string, silent bool, timeout time.Duration) (string, error) {
	command := DefaultCLIPath()
	if c.cliPath != nil {
		if p := c.cliPath(); p != "" {
			command = p
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, command, args...)
	cmd.Stdin = bytes.NewBufferString(input)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}
	cmd.Stderr = cmd.Stdout

	if err := cmd.Start(); err != nil {
		cliErr := startError(command, err)
		if !silent && c.OnLog != nil {
			c.OnLog("\n[ERROR] " + cliErr.Error() + "\n")
		}
		return "", cliErr
	}

	var fullOutput strings.Builder
	buf := make([]byte, 1024)
	for {
		n, err := stdout.Read(buf)
		if n == 0 && err == nil {
			// Prevent tight loop if stdout.Read somehow returns non-blocking 0 bytes
			time.Sleep(10 * time.Millisecond)
			continue
		}

		if n > 0 {
			chunk := string(buf[:n])
			fullOutput.WriteString(chunk)

			if !silent {
				c.mu.Lock()
				c.lastOutput += chunk
				if len(c.lastOutput) > maxLogSize {
					c.lastOutput = c.lastOutput[len(c.lastOutput)-maxLogSize:]
				}
				c.mu.Unlock()

				if c.OnLog != nil {
					c.OnLog(chunk)
				}
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			break
		}
	}

	err = cmd.Wait()
	if ctx.Err() == context.DeadlineExceeded {
		errText := fmt.Sprintf("\n[ERROR] Command timed out after %v. If this persists, try restarting the VPN process.\n", timeout)
		if c.OnLog != nil {
			c.OnLog(errText)
		}
		return fullOutput.String(), fmt.Errorf("command timed out after %v", timeout)
	}

	return fullOutput.String(), err
}
