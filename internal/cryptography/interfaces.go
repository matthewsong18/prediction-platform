package cryptography

type CryptoService interface {
	Encrypt(plaintext string) (string, error)
	Decrypt(ciphertext string) (string, error)
	// GenerateBlindIndex creates a deterministic, non-reversible hash of the plaintext.
	// This is used for searching encrypted data in the database.
	GenerateBlindIndex(plaintext string) string
}