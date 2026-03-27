package vpn

import (
	"bytes"
	"context"
	"fmt"
	"io"
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

	stopReconcile chan struct{}
}

func NewClient(credentials *config.Credentials) *Client {
	c := &Client{
		creds:         credentials,
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
		return false, err
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
		return false, err
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

func (c *Client) executeCommand(args []string, input string, silent bool, timeout time.Duration) (string, error) {
	command := unixCLI
	if runtime.GOOS == "windows" {
		command = windowsCLI
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
		return "", err
	}

	var fullOutput strings.Builder
	buf := make([]byte, 1024)
	for {
		n, err := stdout.Read(buf)
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
