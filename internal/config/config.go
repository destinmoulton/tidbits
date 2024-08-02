package config

import (
	"os"
	"path/filepath"
	"runtime"
)

func getConfigDir() string {
	// Check for the OS type
	switch runtime.GOOS {
	case "windows":
		appData := os.Getenv("APPDATA")
		if appData != "" {
			return filepath.Join(appData, "MyApp")
		}
	case "darwin":
		home := os.Getenv("HOME")
		if home != "" {
			return filepath.Join(home, "Library", "Application Support", "MyApp")
		}
	default: // Unix-like systems
		home := os.Getenv("HOME")
		if home != "" {
			return filepath.Join(home, ".config", "MyApp")
		}
	}

	// Fallback to current directory if no suitable directory is found
	return "."
}
