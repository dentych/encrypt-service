package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

const Secret = "rtj10cv824h19x124jkeh8d91hx2k5jf"

func Encrypt(secret []byte, input io.Reader) (io.Reader, error) {
	ciph, err := aes.NewCipher(secret)
	if err != nil {
		return nil, fmt.Errorf("failed to create new cipher: %w", err)
	}

	iv := make([]byte, aes.BlockSize)
	_, err = io.ReadFull(rand.Reader, iv)
	if err != nil {
		return nil, fmt.Errorf("failed to read random bytes into IV: %w", err)
	}

	reader, writer := io.Pipe()

	ctr := cipher.NewCTR(ciph, iv)

	streamWriter := cipher.StreamWriter{
		S: ctr,
		W: writer,
	}

	_, err = io.Copy(streamWriter, input)
	if err != nil {
		return nil, fmt.Errorf("failed to copy from input file to encrypted output: %w", err)
	}

	return reader, nil
}

func Decrypt(secret []byte, input io.Reader, output io.Writer) error {
	ciph, err := aes.NewCipher(secret)
	if err != nil {
		return fmt.Errorf("failed to create new cipher: %w", err)
	}

	iv := make([]byte, aes.BlockSize)
	_, err = io.ReadFull(input, iv)
	if err != nil {
		return fmt.Errorf("failed to read IV from input file: %w", err)
	}

	ctr := cipher.NewCTR(ciph, iv)
	writer := cipher.StreamWriter{
		S: ctr,
		W: output,
	}

	_, err = io.Copy(writer, input)
	if err != nil {
		return fmt.Errorf("failed to copy from input file to encrypted output: %w", err)
	}

	return nil
}
