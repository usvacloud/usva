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
func EncryptStream(dst io.Writer, src io.ReadSeeker, key []byte) ([]byte, error) {
	cip, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	iv := make([]byte, cip.BlockSize())
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	return iv, cryptLoop(dst, src, cipher.NewCBCEncrypter(cip, iv), 0)
}

// DecryptStream reads chunks from src and writes decrypted chunks to dst.
// Decryption with AES Block cipher mode.
func DecryptStream(dst io.Writer, src io.ReadSeeker, key []byte, iv []byte) error {
	cip, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	return cryptLoop(dst, src, cipher.NewCBCDecrypter(cip, iv), 1)
}

func cryptLoop(dst io.Writer, src io.ReadSeeker, bm cipher.BlockMode, mode int) error {
	for {
		sourceChunk := make([]byte, aes.BlockSize)
		n, err := src.Read(sourceChunk)
		if errors.Is(err, io.EOF) || n == 0 {
			break
		} else if err != nil {
			return err
		}

		chunk := make([]byte, bm.BlockSize())
		if n < bm.BlockSize() && mode == 0 {
			sourceChunk = pad(sourceChunk[:n], bm.BlockSize())
		}

		bm.CryptBlocks(chunk, sourceChunk)

		_, err = src.Read(make([]byte, 1))
		if err != nil && mode == 1 {
			chunk = safeUnpad(chunk, bm.BlockSize())
		} else if err == nil {
			src.Seek(-1, io.SeekCurrent)
		}

		if _, err = dst.Write(chunk); err != nil {
			return err
		}
	}

	return nil
}
