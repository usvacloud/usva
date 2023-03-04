package cryptography

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

var ErrProcessingFailed = errors.New("panic while processing data chunks")

// EncryptStream reads chunks from src and writes them as encrypted to dst.
// Encryption with AES Block cipher mode, returns a random initialization vector and error
func EncryptStream(dst io.Writer, src io.Reader, key []byte) ([]byte, error) {
	cip, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	iv := make([]byte, cip.BlockSize())
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	return iv, cryptLoop(dst, src, cipher.NewCBCEncrypter(cip, iv))
}

// DecryptStream reads chunks from src and writes decrypted chunks to dst.
// Decryption with AES Block cipher mode.
func DecryptStream(dst io.Writer, src io.Reader, key []byte, iv []byte) error {
	cip, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	return cryptLoop(dst, src, cipher.NewCBCDecrypter(cip, iv))
}

func cryptLoop(dst io.Writer, src io.Reader, bm cipher.BlockMode) error {
	for {
		plaintextChunk := make([]byte, aes.BlockSize)
		n, err := src.Read(plaintextChunk)
		if errors.Is(err, io.EOF) || n == 0 {
			break
		} else if err != nil {
			return err
		}

		if n < bm.BlockSize() {
			plaintextChunk = pad(plaintextChunk[:n], bm.BlockSize())
		}

		chunk := make([]byte, bm.BlockSize())
		bm.CryptBlocks(chunk, plaintextChunk)

		chunk = safeUnpad(chunk, bm.BlockSize())
		if _, err = dst.Write(chunk); err != nil {
			return err
		}
	}

	return nil
}
