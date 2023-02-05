package cryptography

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"log"
	"math"
	"testing"

	"golang.org/x/crypto/argon2"
)

func check(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
	}
}

func TestChunking(t *testing.T) {
	key := argon2.Key([]byte("this is my awesome key"), []byte{}, 2, 1024, 2, 32)

	// prepare IV
	iv := make([]byte, aes.BlockSize)
	_, err := io.ReadFull(rand.Reader, iv)
	check(t, err)

	sourceBuf := make([]byte, aes.BlockSize)
	_, err = io.ReadFull(rand.Reader, sourceBuf)
	check(t, err)

	cip, err := aes.NewCipher(key)
	check(t, err)

	bm := cipher.NewCBCEncrypter(cip, iv)
	ebuf := make([]byte, bm.BlockSize())
	bm.CryptBlocks(ebuf, sourceBuf)

	bm = cipher.NewCBCDecrypter(cip, iv)
	dbuf := make([]byte, bm.BlockSize())
	bm.CryptBlocks(dbuf, ebuf)

	if !bytes.Equal(sourceBuf, dbuf) {
		t.Fatal("encryption and decryption didn't return wanted result")
	}
}

func TestEncryptDecrypt(t *testing.T) {
	// prepare ciphers
	key := argon2.Key([]byte("this is my awesome key"), []byte{}, 2, 1024, 2, 16)
	cip, err := aes.NewCipher(key)
	check(t, err)

	// prepare IV
	iv := make([]byte, aes.BlockSize)
	_, err = io.ReadFull(rand.Reader, iv)
	check(t, err)

	encryptSrc := make([]byte, 10000)
	_, err = io.ReadFull(rand.Reader, encryptSrc)
	check(t, err)

	// test encryption
	encryptDst := bytes.NewBuffer(nil)
	encryptionBlockMode := cipher.NewCBCEncrypter(cip, iv)
	encryptSrcReader := bytes.NewReader(encryptSrc)
	err = cryptLoop(encryptDst, encryptSrcReader, encryptionBlockMode)
	check(t, err)

	// test decryption
	decryptDst := bytes.NewBuffer(nil)
	decryptionBlockMode := cipher.NewCBCDecrypter(cip, iv)
	err = cryptLoop(decryptDst, encryptDst, decryptionBlockMode)
	check(t, err)

	// verify output
	bs := cip.BlockSize()
	chunksVerified := 0
	verifyChunk := func(decBuffer []byte, offset, read int) {
		sliceto := offset + read
		if len(encryptSrc) < sliceto {
			sliceto = len(encryptSrc) - 1
		}

		plaintextSlice := encryptSrc[offset:sliceto]
		if !bytes.Equal(plaintextSlice, decBuffer) {
			t.Error(plaintextSlice, "!=", decBuffer)
			t.Fatalf("%d out of %d chunks were verified before corruption.", chunksVerified, len(encryptSrc)/bs)
		}
		chunksVerified++
	}

	offset := 0
	for len(encryptSrc) > offset+1 {
		decBuffer := make([]byte, bs)
		bitsRead, err := decryptDst.Read(decBuffer)
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}

		verifyChunk(decBuffer, offset, bitsRead)
		offset += int(math.Min(float64(bitsRead), float64(len(encryptSrc))-1))
	}
}
