package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func Read() (Config, error) {

	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "couldn't get home dir: %v\n", err)
		os.Exit(1)
		return Config{}, err

	}

	path := filepath.Join(home, ConfigFileName)

	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "couldn't read file %q: %v\n", path, err)
		os.Exit(1)
		return Config{}, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		fmt.Fprintf(os.Stderr, "couldn't parse json: %v\n", err)
		os.Exit(1)
		return Config{}, err
	}

	return cfg, nil

}

func Write(cfg Config) error {
	fullPath, err := GetConfigFilePath()
	if err != nil {
		return err
	}

	file, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(cfg)
	if err != nil {
		return err
	}

	return nil

}
