package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"log"
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

	go func() {
		_, err = writer.Write(iv)
		if err != nil {
			log.Printf("failed to write iv: %s", err)
		}

		_, err = io.Copy(streamWriter, input)
		if err != nil {
			log.Printf("failed to copy: %s\n", err)
		}
		err = streamWriter.Close()
		if err != nil {
			log.Printf("failed to close streamwriter: %s\n", err)
		}
	}()

	return reader, nil
}

func Decrypt(secret []byte, input io.Reader) (io.Reader, error) {
	ciph, err := aes.NewCipher(secret)
	if err != nil {
		return nil, fmt.Errorf("failed to create new cipher: %w", err)
	}

	iv := make([]byte, aes.BlockSize)
	_, err = io.ReadFull(input, iv)
	if err != nil {
		return nil, fmt.Errorf("failed to read IV from input file: %w", err)
	}

	reader, writer := io.Pipe()

	ctr := cipher.NewCTR(ciph, iv)
	streamWriter := cipher.StreamWriter{
		S: ctr,
		W: writer,
	}

	go func() {
		_, err = io.Copy(streamWriter, input)
		if err != nil {
			log.Printf("failed to copy decrypted output: %s", err)
		}
		err = streamWriter.Close()
		if err != nil {
			log.Printf("failed to close streamwriter: %s", err)
		}
	}()

	return reader, nil
}
