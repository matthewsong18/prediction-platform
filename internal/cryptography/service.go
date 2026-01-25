package cryptography

import (
	"crypto/aes"
	"crypto/cipher"
	"log"
)

type service struct {
	gcm cipher.AEAD
}

func NewService(secretKey [32]byte) (CryptoService, error) {
	// Uses AES-GCM for reversible, authenticated encryption. This protects user IDs
	// at rest while allowing decryption for display features like the leaderboard,
	// and ensures data integrity (tamper-proofing) by default.
	block, err := aes.NewCipher(secretKey[:])
	if err != nil {
		log.Fatal(err)
	}
	gcm, err := cipher.NewGCMWithRandomNonce(block)
	if err != nil {
		log.Fatal(err)
	}

	return &service{
		gcm,
	}, nil
}

// Encrypt implements [CryptoService].
func (s *service) Encrypt(plaintext string) (string, error) {
	ciphertext := s.gcm.Seal(nil, nil, []byte(plaintext), nil)
	return string(ciphertext), nil
}

// Decrypt implements [CryptoService].
func (s *service) Decrypt(ciphertext string) (string, error) {
	plaintext, err := s.gcm.Open(nil, nil, []byte(ciphertext), nil)
	if err != nil {
		log.Fatal(err)
	}
	return string(plaintext), nil
}

var _ CryptoService = (*service)(nil)
