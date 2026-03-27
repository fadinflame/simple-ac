package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/fadinflame/simple-ac/security"
)

type Manager struct {
	mu  sync.RWMutex
	cfg *Config
}

func NewManager() *Manager {
	return &Manager{}
}

// configDir returns the directory where configurations are stored
func (m *Manager) configDir() string {
	if runtime.GOOS == "windows" {
		return filepath.Join(os.Getenv("APPDATA"), "simple-ac")
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".simple-ac")
}

// ConfigPath returns the path to the main application configuration file
func (m *Manager) ConfigPath() string {
	return filepath.Join(m.configDir(), "config.json")
}

// ConfigExists checks if the configuration file is present
func (m *Manager) ConfigExists() bool {
	_, err := os.Stat(m.ConfigPath())
	return err == nil
}

// LoadConfig reads the configuration from the file system
func (m *Manager) LoadConfig() error {
	data, err := os.ReadFile(m.ConfigPath())
	if err != nil {
		return err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return err
	}

	m.mu.Lock()
	m.cfg = &cfg
	m.mu.Unlock()
	return nil
}

// SaveConfig persists the configuration to disk and updates in-memory state
func (m *Manager) SaveConfig(cfg *Config) error {
	if err := os.MkdirAll(m.configDir(), 0700); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(m.ConfigPath(), data, 0600); err != nil {
		return err
	}

	m.mu.Lock()
	m.cfg = cfg
	m.mu.Unlock()
	return nil
}

// GetConfig returns the current configuration, loading it if necessary
func (m *Manager) GetConfig() *Config {
	m.mu.RLock()
	c := m.cfg
	m.mu.RUnlock()

	if c == nil {
		_ = m.LoadConfig()
		m.mu.RLock()
		c = m.cfg
		m.mu.RUnlock()
	}

	if c == nil {
		return &Config{}
	}
	return c
}

// --- Credentials Management ---

// SaveCredentials encrypts and saves VPN credentials to a file
func (m *Manager) SaveCredentials(creds *Credentials, password string, path string) error {
	data, err := json.Marshal(creds)
	if err != nil {
		return err
	}

	encrypted, err := security.Encrypt(data, password)
	if err != nil {
		return err
	}

	return os.WriteFile(path, encrypted, 0600)
}

// LoadCredentials reads and decrypts VPN credentials from a file
func (m *Manager) LoadCredentials(password string, path string) (*Credentials, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	decrypted, err := security.Decrypt(data, password)
	if err != nil {
		return nil, err
	}

	var creds Credentials
	if err := json.Unmarshal(decrypted, &creds); err != nil {
		return nil, err
	}

	return &creds, nil
}
