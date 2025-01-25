package internal

import (
	"copy-stamp/config"
	"os"
	"strings"
)

func EnsureDataDirectoryExists() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	templatesDir := strings.Join([]string{homeDir, config.TemplatesLocation}, "/")
	err = os.MkdirAll(templatesDir, 0744)

	return templatesDir, err
}
