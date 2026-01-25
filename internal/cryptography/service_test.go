package cryptography

import (
	"crypto/rand"
	"io"
	"log"
	"testing"
)

func TestNewService(t *testing.T) {
	var key [32]byte
	if _, err := io.ReadFull(rand.Reader, key[:]); err != nil {
		log.Fatal(err)
	}

	_, err := NewService(key)
	if err != nil {
		t.Fatalf("Failed to create service: %v", err)
	}
}

func TestEncryptDecrypt(t *testing.T) {
	var key [32]byte
	if _, err := io.ReadFull(rand.Reader, key[:]); err != nil {
		log.Fatal(err)
	}

	service, err := NewService(key)
	if err != nil {
		t.Fatalf("Failed to create service: %v", err)
	}

	original := "Hello, World!"
	encrypted, err := service.Encrypt(original)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	if encrypted == original {
		t.Error("Encrypted text should not match original text")
	}

	decrypted, err := service.Decrypt(encrypted)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	if decrypted != original {
		t.Errorf("Expected decrypted text to match original. Got %q, want %q", decrypted, original)
	}
}
