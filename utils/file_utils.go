package utils

import (
	"errors"
	"os"
	"strings"
)

// FileExists checks that a file exists
func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !errors.Is(err, os.ErrNotExist)
}

// FileContains checks that a file contains a string
func FileContains(filePath, match string) bool {
	if !FileExists(filePath) {
		return false
	}

	text, _ := os.ReadFile(filePath)

	return strings.Contains(string(text), match)
}
