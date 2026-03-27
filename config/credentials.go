package config

import "errors"

type Credentials struct {
	Server    string `json:"server"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Group     int    `json:"group"`
	OTPSecret string `json:"otp_secret"`
}

func (c *Credentials) Validate() error {
	if c.Server == "" {
		return errors.New("server is required")
	}
	if c.Username == "" {
		return errors.New("username is required")
	}
	if c.Password == "" {
		return errors.New("password is required")
	}
	return nil
}
