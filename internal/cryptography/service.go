package cryptography

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

type service struct {
	gcm       cipher.AEAD
	secretKey [32]byte
}

func NewService(secretKey [32]byte) (CryptoService, error) {
	// Uses AES-GCM for reversible, authenticated encryption. This protects user IDs
	// at rest while allowing decryption for display features like the leaderboard,
	// and ensures data integrity (tamper-proofing) by default.
	block, err := aes.NewCipher(secretKey[:])
	if err != nil {
		return nil, fmt.Errorf("failed to create aes cipher block: %w", err)
	}
	gcm, err := cipher.NewGCMWithRandomNonce(block)
	if err != nil {
		return nil, fmt.Errorf("failed to wrap block with gcm: %w", err)
	}

	return &service{
		gcm:       gcm,
		secretKey: secretKey,
	}, nil
}

// Encrypt implements [CryptoService].
// Returns a Base64 encoded string of the ciphertext.
func (s *service) Encrypt(plaintext string) (string, error) {
	ciphertext := s.gcm.Seal(nil, nil, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt implements [CryptoService].
// Expects a Base64 encoded string of the ciphertext.
func (s *service) Decrypt(encodedCiphertext string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encodedCiphertext)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64 ciphertext: %w", err)
	}

	plaintext, err := s.gcm.Open(nil, nil, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt the text: %w", err)
	}
	return string(plaintext), nil
}

// GenerateBlindIndex implements [CryptoService].
func (s *service) GenerateBlindIndex(plaintext string) string {
	h := hmac.New(sha256.New, s.secretKey[:])
	h.Write([]byte(plaintext))
	return hex.EncodeToString(h.Sum(nil))
}

var _ CryptoService = (*service)(nil)
