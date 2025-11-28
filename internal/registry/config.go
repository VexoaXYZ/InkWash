package registry

import (
	"os"
	"path/filepath"
	"runtime"
)

// GetDefaultConfigPath returns the default config directory path
func GetDefaultConfigPath() string {
	if runtime.GOOS == "windows" {
		appData := os.Getenv("APPDATA")
		if appData == "" {
			appData = filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Roaming")
		}
		return filepath.Join(appData, "inkwash")
	}

	// Linux/macOS
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "inkwash")
}

// GetDefaultCachePath returns the default cache directory path
func GetDefaultCachePath() string {
	if runtime.GOOS == "windows" {
		localAppData := os.Getenv("LOCALAPPDATA")
		if localAppData == "" {
			localAppData = filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Local")
		}
		return filepath.Join(localAppData, "inkwash", "cache", "fxserver")
	}

	// Linux/macOS
	home, _ := os.UserHomeDir()
	cacheHome := os.Getenv("XDG_CACHE_HOME")
	if cacheHome == "" {
		cacheHome = filepath.Join(home, ".cache")
	}
	return filepath.Join(cacheHome, "inkwash", "fxserver")
}

// GetDefaultDataPath returns the default data directory path
func GetDefaultDataPath() string {
	if runtime.GOOS == "windows" {
		localAppData := os.Getenv("LOCALAPPDATA")
		if localAppData == "" {
			localAppData = filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Local")
		}
		return filepath.Join(localAppData, "inkwash", "data")
	}

	// Linux/macOS
	home, _ := os.UserHomeDir()
	dataHome := os.Getenv("XDG_DATA_HOME")
	if dataHome == "" {
		dataHome = filepath.Join(home, ".local", "share")
	}
	return filepath.Join(dataHome, "inkwash")
}

// GetRegistryPath returns the path to the servers.json registry file
func GetRegistryPath() string {
	return filepath.Join(GetDefaultConfigPath(), "servers.json")
}

// GetConfigFilePath returns the path to the config.yaml file
func GetConfigFilePath() string {
	return filepath.Join(GetDefaultConfigPath(), "config.yaml")
}
