package utils

import (
	"fmt"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
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

func BytesToMegabytesBinary(bytes int64) float64 {
	return float64(bytes) / 1048576
}

func AddCommas(n int) string {
	p := message.NewPrinter(language.English)
	return p.Sprintf("%d", n)
}
