package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

const APP_NAME = "tidbits"

func GetConfigDir() string {
	home := os.Getenv("HOME")
	if home != "" {
		path := filepath.Join(home, ".config", APP_NAME)
		_, err := os.Stat(path)
		if os.IsNotExist(err) {
			// Make the dir if it doesn't exist
			err := os.Mkdir(path, 0770)
			if err != nil {
				fmt.Println("GetConfigDir :: cannot create config dir", err)
				os.Exit(2)
			}
		}
		return path
	} else {
		fmt.Println("GetConfigDir :: cannot locate $HOME directory")
		os.Exit(2)
	}
	return ""

}
