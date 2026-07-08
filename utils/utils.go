package utils

import (
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

var APP_CONF_DIR string

func init() {
	conf, err := os.UserConfigDir()
	if err != nil {
		log.Panicf("Error getting user config dir: %v\n", err)
	}
	APP_CONF_DIR = filepath.Join(conf, "riptvtime")
	log.Println(APP_CONF_DIR)
	err = os.MkdirAll(APP_CONF_DIR, 0755)
	if err != nil {
		log.Panicf("Error creating user config dir: %v - %v\n", conf, err)
	}
}

func IsGoRun() bool {
	execPath, err := os.Executable()
	if err != nil {
		return false
	}
	slog.Debug("Executable path", "path", execPath)
	return strings.HasPrefix(execPath, filepath.Join(os.TempDir(), "go-build"))
}
