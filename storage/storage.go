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
	if _, err := os.Stat(filepath.Join(storageDirectory, filename)); err == nil {
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
	defer file.Close()

	_, err = io.Copy(file, stream)
	if err != nil {
		return fmt.Errorf("failed to copy: %s", err)
	}

	return nil
}

func StreamFromStorage(filename string) (*os.File, error) {
	file, err := os.Open(filepath.Join(storageDirectory, filename))
	if err != nil {
		return nil, fmt.Errorf("failed to open file for writing: %s", err)
	}

	return file, nil
}
