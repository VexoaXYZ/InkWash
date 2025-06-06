package config

import (
	"os"
	"path/filepath"
)

type Config struct {
	DefaultServerPath string
	CacheDir          string
	TemplatesDir      string
}

func Load() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configDir := filepath.Join(homeDir, ".inkwash")
	
	return &Config{
		DefaultServerPath: filepath.Join(homeDir, "fivem-servers"),
		CacheDir:          filepath.Join(configDir, "cache"),
		TemplatesDir:      filepath.Join(configDir, "templates"),
	}, nil
}

func (c *Config) EnsureDirectories() error {
	dirs := []string{
		c.DefaultServerPath,
		c.CacheDir,
		c.TemplatesDir,
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	return nil
}