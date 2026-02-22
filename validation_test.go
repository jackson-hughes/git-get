package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidateGitExists(t *testing.T) {
	t.Run("git exists in PATH", func(t *testing.T) {
		mockLookPath := func(name string) (string, error) {
			return "/usr/bin/git", nil
		}

		err := validateGitExists(mockLookPath)
		if err != nil {
			t.Errorf("validateGitExists() should return nil when git is in PATH, got: %v", err)
		}
	})

	t.Run("git not in PATH", func(t *testing.T) {
		mockLookPath := func(name string) (string, error) {
			return "", fmt.Errorf("executable file not found in $PATH")
		}

		err := validateGitExists(mockLookPath)
		if err == nil {
			t.Error("validateGitExists() should return error when git is not in PATH")
		} else {
			expectedMsg := "git binary not found in PATH"
			if !strings.Contains(err.Error(), expectedMsg) {
				t.Errorf("error message should contain %q, got %q", expectedMsg, err.Error())
			}
		}
	})
}

func TestCheckDirectoryExists(t *testing.T) {
	t.Run("directory exists", func(t *testing.T) {
		tempDir := t.TempDir()

		exists, err := checkDirectoryExists(tempDir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !exists {
			t.Errorf("expected directory to exist, got exists=%v", exists)
		}
	})

	t.Run("directory does not exist", func(t *testing.T) {
		nonExistentPath := filepath.Join(t.TempDir(), "nonexistent")

		exists, err := checkDirectoryExists(nonExistentPath)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if exists {
			t.Errorf("expected directory to not exist, got exists=%v", exists)
		}
	})

	t.Run("path is a file not directory", func(t *testing.T) {
		tempDir := t.TempDir()
		tempFile := filepath.Join(tempDir, "testfile")
		if err := os.WriteFile(tempFile, []byte("test"), 0644); err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}

		exists, err := checkDirectoryExists(tempFile)
		if err == nil {
			t.Errorf("expected error when path is a file, got none")
		}
		if exists {
			t.Errorf("expected exists=false for file path, got exists=%v", exists)
		}
	})
}
