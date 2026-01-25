package cryptography

type service struct {
}

func NewService(secretKey [32]byte) (CryptoService, error) {
	panic("unimplemented")
}

// Encrypt implements [CryptoService].
func (s *service) Encrypt(plaintext string) (string, error) {
	panic("unimplemented")
}

// Decrypt implements [CryptoService].
func (s *service) Decrypt(ciphertext string) (string, error) {
	panic("unimplemented")
}

var _ CryptoService = (*service)(nil)
