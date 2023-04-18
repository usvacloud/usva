package cryptography

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

var ErrProcessingFailed = errors.New("panic while processing data chunks")

func EncryptStream(ciphertext io.Writer, plaintext io.Reader, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return []byte{}, err
	}

	// Generate a random IV
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return iv, err
	}

	// Create a stream cipher using AES in CBC mode
	stream := cipher.NewCFBEncrypter(block, iv)

	// Encrypt the plaintext and write to the ciphertext
	writer := &cipher.StreamWriter{S: stream, W: ciphertext}
	if _, err := io.Copy(writer, plaintext); err != nil {
		return iv, err
	}

	return iv, nil
}

func DecryptStream(plaintext io.Writer, ciphertext io.Reader, key []byte, iv []byte) error {
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	// Create a stream cipher using AES in CBC mode
	stream := cipher.NewCFBDecrypter(block, iv)

	// Decrypt the ciphertext and write to the plaintext
	reader := &cipher.StreamReader{S: stream, R: ciphertext}
	if _, err := io.Copy(plaintext, reader); err != nil {
		return err
	}

	return nil
}
