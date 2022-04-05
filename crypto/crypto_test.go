package crypto

import (
	"io"
	"strings"
	"testing"
)

const TestSecret = "lk1j24lkj1lk4j1lk2j41lkj41lkj2kk"

func TestEncrypt(t *testing.T) {
	input := strings.NewReader("this is some string")
	encryptedData, err := Encrypt([]byte(TestSecret), input)
	if err != nil {
		t.Fatalf("failed to encrypt data: %s", err)
	}

	output, err := io.ReadAll(encryptedData)
	if err != nil {
		t.Fatalf("failed to read all: %s", err)
	}

	if len(output) != 35 {
		t.Fatalf("output should be 19 chars long, but was %d", len(output))
	}
}
