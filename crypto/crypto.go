package crypto

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"os"
)

func EncryptFile(secret []byte, inputFile, outputFile string) error {
	ciph, err := aes.NewCipher(secret)
	if err != nil {
		return fmt.Errorf("failed to create new cipher: %w", err)
	}

	input, err := os.Open(inputFile)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer input.Close()

	iv := make([]byte, aes.BlockSize)
	_, err = io.ReadFull(rand.Reader, iv)
	if err != nil {
		return fmt.Errorf("failed to read random bytes into IV: %w", err)
	}

	ctr := cipher.NewCTR(ciph, iv)

	outputFp, err := os.OpenFile(outputFile, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open output file for writing: %w", err)
	}
	defer outputFp.Close()

	output := bufio.NewWriter(outputFp)
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

func DecryptFile(secret []byte, inputFile, outputFile string) error {
	ciph, err := aes.NewCipher(secret)
	if err != nil {
		return fmt.Errorf("failed to create new cipher: %w", err)
	}

	input, err := os.Open(inputFile)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer input.Close()

	outputFp, err := os.OpenFile(outputFile, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open output file for writing: %w", err)
	}
	defer outputFp.Close()

	output := bufio.NewWriter(outputFp)

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
