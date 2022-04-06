package api

import (
	"fmt"
	"github.com/dentych/encrypt-service/crypto"
	"github.com/dentych/encrypt-service/storage"
	"io"
)

var (
	ErrFileExists    = fmt.Errorf("file already exists")
	ErrFileNotExists = fmt.Errorf("file doesn't exist")
)

func SaveFile(filename string, content io.Reader) error {
	if exists, err := storage.FileExists(filename); err != nil {
		return fmt.Errorf("failed to check if file exists: %w", err)
	} else if exists {
		return ErrFileExists
	}

	encryptedContent, err := crypto.Encrypt([]byte(crypto.Secret), content)
	if err != nil {
		return fmt.Errorf("failed to encrypt: %w", err)
	}
	encryptedContent.Close()

	err = storage.StreamToStorage(filename, encryptedContent)
	if err != nil {
		return err
	}

	return nil
}

func RetrieveFile(filename string) (*io.PipeReader, error) {
	if exists, err := storage.FileExists(filename); err != nil {
		return nil, fmt.Errorf("failed to check if file exists: %w", err)
	} else if !exists {
		return nil, ErrFileNotExists
	}

	file, err := storage.StreamFromStorage(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to stream from storage: %w", err)
	}

	reader, err := crypto.Decrypt([]byte(crypto.Secret), file)
	if err != nil {
		return nil, fmt.Errorf("decrypting failed: %w", err)
	}

	return reader, nil
}
