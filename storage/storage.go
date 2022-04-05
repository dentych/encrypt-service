package storage

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

var storageDirectory string

func Setup(storageDir string) {
	storageDirectory = storageDir
}

func FileExists(filename string) (bool, error) {
	if _, err := os.Stat(filename); err == nil {
		return true, nil
	} else if errors.Is(err, os.ErrNotExist) {
		return false, nil
	} else {
		return false, fmt.Errorf("failed to stat file: %w", err)
	}
}

func StreamToStorage(filename string, stream io.Reader) error {
	file, err := os.OpenFile(filepath.Join(storageDirectory, filename), os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file for writing: %s", err)
	}

	_, err = io.Copy(file, stream)
	if err != nil {
		return fmt.Errorf("failed to copy: %s", err)
	}

	return nil
}
