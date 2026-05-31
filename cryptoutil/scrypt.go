package cryptoutil

import (
	"fmt"

	"golang.org/x/crypto/scrypt"
)

type scrypter struct {
	n, r, p, keyLen int
}

// NewScrypter creates a Scrypter with the given config.
// If config is zero-valued, DefaultScryptConfig is used.
// Returns ErrInvalidParams if N is not a power of two > 1.
func NewScrypter(cfg ScryptConfig) (Scrypter, error) {
	if cfg == (ScryptConfig{}) {
		cfg = DefaultScryptConfig()
	}
	if cfg.N <= 1 || cfg.N&(cfg.N-1) != 0 {
		return nil, fmt.Errorf("%w: N must be a power of two > 1", ErrInvalidParams)
	}
	return &scrypter{n: cfg.N, r: cfg.R, p: cfg.P, keyLen: cfg.KeyLen}, nil
}

func (s *scrypter) DeriveKey(password, salt []byte) ([]byte, error) {
	key, err := scrypt.Key(password, salt, s.n, s.r, s.p, s.keyLen)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrKeyDerivation, err)
	}
	return key, nil
}
