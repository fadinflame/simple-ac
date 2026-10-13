package config

import (
	"errors"
)

type Config struct {
	CredentialsFilePath string `json:"credentials_file_path"`
	// CLIPath overrides the default Cisco CLI location; empty means default.
	CLIPath string `json:"cli_path,omitempty"`
}

func (c *Config) Validate() error {
	if c.CredentialsFilePath == "" {
		return errors.New("credentials_file_path is required")
	}
	return nil
}
