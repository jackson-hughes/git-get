package main

import (
	"fmt"
	"os"
)

func validateGitExists(lookPath func(string) (string, error)) error {
	_, err := lookPath("git")
	if err != nil {
		return fmt.Errorf("git binary not found in PATH. Please install git to use git-get")
	}
	return nil
}

func checkDirectoryExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("error checking directory %s: %w", path, err)
	}

	if !info.IsDir() {
		return false, fmt.Errorf("path %s exists but is not a directory", path)
	}

	return true, nil
}
