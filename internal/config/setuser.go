package config

import (
	"os"
	"path/filepath"
)

func (c *Config) SetUser(userName string) error {
	c.CurrentUserName = userName
	return nil
}

func GetConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	fullPath := filepath.Join(home, ConfigFileName)

	return fullPath, nil
}
