package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const defaultDirName = ".stashctl"

// Config holds application-level configuration.
type Config struct {
	StoreDir   string `json:"store_dir"`
	DefaultFmt string `json:"default_format"` // "table" | "detail"
}

// Default returns a Config populated with sensible defaults.
func Default() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	return &Config{
		StoreDir:   filepath.Join(home, defaultDirName),
		DefaultFmt: "table",
	}, nil
}

// Load reads config from path. Missing file returns Default.
func Load(path string) (*Config, error) {
	cfg, err := Default()
	if err != nil {
		return nil, err
	}
	f, err := os.Open(path)
	if os.IsNotExist(err) {
		return cfg, nil
	}
	if err != nil {
		return nil, err
	}
	defer f.Close()
	if err := json.NewDecoder(f).Decode(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

// Save writes cfg as JSON to path, creating parent directories as needed.
func Save(path string, cfg *Config) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(cfg)
}

// DefaultPath returns the default config file path.
func DefaultPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, defaultDirName, "config.json"), nil
}
