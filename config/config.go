package config

import (
	"errors"
)

type Config struct {
	CredentialsFilePath   string `json:"credentials_file_path"`
	ShowConsoleLog        bool   `json:"show_console_log"`
}

func (c *Config) Validate() error {
	if c.CredentialsFilePath == "" {
		return errors.New("credentials_file_path is required")
	}
	return nil
}
