package api

import (
	"fmt"
	"github.com/dentych/encrypt-service/crypto"
	"github.com/dentych/encrypt-service/storage"
	"io"
)

var (
	ErrFileExists = fmt.Errorf("file already exists")
)

type api struct {
}

func (a *api) SaveFile(filename string, content io.Reader) error {
	if exists, err := storage.FileExists(filename); err != nil {
		return fmt.Errorf("failed to check if file exists: %w", err)
	} else if exists {
		return ErrFileExists
	}

	encryptedContent, err := crypto.Encrypt([]byte(crypto.Secret), content)
	if err != nil {
		return fmt.Errorf("failed to encrypt: %w", err)
	}

	err = storage.StreamToStorage("testfile.txt", encryptedContent)
	if err != nil {
		return err
	}

	return nil
}
