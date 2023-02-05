package cryptography

import (
	"errors"

	"golang.org/x/crypto/argon2"
)

var (
	ErrPasswordTooShort = errors.New("password is too short")
	ErrPasswordTooLong  = errors.New("password is too long")
)

func DeriveBasicKey(password []byte, length uint32) ([]byte, error) {
	if len(password) > 128 {
		return []byte{}, ErrPasswordTooLong
	} else if len(password) < 6 {
		return []byte{}, ErrPasswordTooShort
	}

	return argon2.Key(password, []byte{}, length*8, length*1024, 2, length), nil
}
