package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/fadinflame/simple-ac/config"
	"github.com/fadinflame/simple-ac/updater"
	"github.com/fadinflame/simple-ac/vpn"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct provides all the methods that can be called from the frontend
type App struct {
	mu            sync.RWMutex
	ctx           context.Context
	ConfigManager *config.Manager
	VPNClient     *vpn.Client
	isConnecting  bool
}

func NewApp() *App {
	return &App{
		ConfigManager: config.NewManager(),
	}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	runtime.WindowShow(ctx)
	runtime.WindowCenter(ctx)

	_ = a.ConfigManager.LoadConfig()
	a.Log("Application started")
}

// setVPNClient initializes or updates the VPN client with credentials and sets up logging
func (a *App) setVPNClient(creds *config.Credentials) {
	a.mu.Lock()
	defer a.mu.Unlock()

	// If there's an existing client, stop its reconciliation loop
	if a.VPNClient != nil {
		a.VPNClient.Stop()
	}

	a.VPNClient = vpn.NewClient(creds)
	a.VPNClient.OnLog = func(log string) {
		if a.ctx != nil {
			runtime.EventsEmit(a.ctx, "vpn-log", log)
		}
	}
}

// getVPNClient safely returns the current VPN client
func (a *App) getVPNClient() *vpn.Client {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.VPNClient
}

// SaveConfig encrypts credentials and saves them to the specified path
func (a *App) SaveConfig(creds config.Credentials, password string, path string) error {
	if err := a.ConfigManager.SaveCredentials(&creds, password, path); err != nil {
		return err
	}

	cfg := a.ConfigManager.GetConfig()
	cfg.CredentialsFilePath = path

	if err := a.ConfigManager.SaveConfig(cfg); err != nil {
		return err
	}

	a.setVPNClient(&creds)
	a.Log(fmt.Sprintf("Configuration saved to %s", path))
	return nil
}

// GetConfig returns the current application configuration
func (a *App) GetConfig() *config.Config {
	return a.ConfigManager.GetConfig()
}

// UpdateConfig updates the application configuration
func (a *App) UpdateConfig(cfg config.Config) error {
	err := a.ConfigManager.SaveConfig(&cfg)
	if err == nil {
		a.Log("Settings updated")
	}
	return err
}

// Unlock attempts to decrypt credentials with the master password
func (a *App) Unlock(password string) (bool, error) {
	if !a.ConfigManager.ConfigExists() {
		return false, fmt.Errorf("configuration not found")
	}

	if err := a.ConfigManager.LoadConfig(); err != nil {
		return false, err
	}

	cfg := a.ConfigManager.GetConfig()
	creds, err := a.ConfigManager.LoadCredentials(password, cfg.CredentialsFilePath)
	if err != nil {
		return false, fmt.Errorf("invalid master password or corrupted file")
	}

	a.setVPNClient(creds)
	a.Log("Master password unlocked successfully")
	return true, nil
}

// CheckConfig returns true if the main config exists
func (a *App) CheckConfig() bool {
	return a.ConfigManager.ConfigExists()
}

// VPN Wrapper actions

func (a *App) setIsConnecting(val bool) {
	a.mu.Lock()
	a.isConnecting = val
	a.mu.Unlock()
	if a.ctx != nil {
		runtime.EventsEmit(a.ctx, "connection-processing", val)
	}
}

func (a *App) IsConnecting() bool {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.isConnecting
}

func (a *App) Connect() (bool, error) {
	client := a.getVPNClient()
	if client == nil {
		return false, fmt.Errorf("VPN client not initialized")
	}

	a.setIsConnecting(true)
	defer a.setIsConnecting(false)

	a.Log("Connection request initiated")
	return client.Connect()
}

func (a *App) Disconnect() (bool, error) {
	client := a.getVPNClient()
	if client == nil {
		return false, fmt.Errorf("VPN client not initialized")
	}

	a.setIsConnecting(true)
	defer a.setIsConnecting(false)

	a.Log("Disconnect request initiated")
	return client.Disconnect()
}

func (a *App) IsConnected() (bool, error) {
	client := a.getVPNClient()
	if client == nil {
		return false, nil
	}
	return client.IsConnected()
}

func (a *App) SyncStatus() (bool, error) {
	client := a.getVPNClient()
	if client == nil {
		return false, nil
	}
	return client.SyncStatus()
}

// Lock clears sensitive data from memory and disconnects VPN
func (a *App) Lock() error {
	a.mu.Lock()
	client := a.VPNClient
	a.VPNClient = nil
	a.mu.Unlock()

	if client != nil {
		_, _ = client.Disconnect()
		client.Stop()
	}
	a.Log("Application locked")
	return nil
}

func (a *App) GetVPNLogs() string {
	client := a.getVPNClient()
	if client == nil {
		return ""
	}
	return client.GetLogs()
}

// Other Utilities

func (a *App) CheckForUpdates() (*updater.Release, error) {
	return updater.CheckForUpdates(Version)
}

func (a *App) SelectConfigPath() string {
	path, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title:           "Where to save the config?",
		DefaultFilename: "simplac_cfg.enc",
		Filters: []runtime.FileFilter{
			{DisplayName: "Encrypted files (*.enc)", Pattern: "*.enc"},
		},
	})

	if err != nil || path == "" {
		return ""
	}
	return path
}

func (a *App) CenterWindow() {
	if a.ctx != nil {
		runtime.WindowCenter(a.ctx)
	}
}

func (a *App) GetAppVersion() string {
	return Version
}

func (a *App) Log(message string) {
	if a.ctx != nil {
		runtime.EventsEmit(a.ctx, "app-log", message)
	}
}
