package cryptography

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"io"
	"math"
	"reflect"
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
	key := argon2.Key([]byte("this is my awesome key"), []byte{}, 2, 1024, 2, 32)

	size := int(math.Pow(2, 20) * 4)
	encryptSrc := make([]byte, size)
	var totalDecryptBuffer []byte

	_, err := io.ReadFull(rand.Reader, encryptSrc)
	check(t, err)
	integrityOriginal := sha256.Sum256(encryptSrc)

	// test encryption
	encryptDst := bytes.NewBuffer(nil)
	encryptSrcReader := bytes.NewReader(encryptSrc)
	iv, err := EncryptStream(encryptDst, encryptSrcReader, key)
	check(t, err)

	// test decryption
	decryptDst := bytes.NewBuffer(nil)
	encryptDstReader := bytes.NewReader(encryptDst.Bytes())
	check(t, DecryptStream(decryptDst, encryptDstReader, key, iv))

	// verify output
	bs := aes.BlockSize
	chunksVerified := 0
	verifyChunk := func(buf []byte, offset, read int) {
		outset := offset + read
		if len(encryptSrc) < outset {
			outset = len(encryptSrc)
		}

		plaintextSlice := encryptSrc[offset:outset]
		buf = buf[:len(plaintextSlice)]
		if !bytes.Equal(plaintextSlice, buf) {
			t.Errorf("%v (%d) != %v (%d)", plaintextSlice, len(plaintextSlice), buf, len(buf))
			totalchunks := len(encryptSrc) / bs
			if len(encryptSrc)%bs != 0 {
				totalchunks++
			}
			t.Fatalf("%d out of %d chunks were verified before corruption was identified.", chunksVerified, totalchunks)
		}

		totalDecryptBuffer = append(totalDecryptBuffer, buf...)
		chunksVerified++
	}

	offset := 0
	for len(encryptSrc) > offset {
		decBuffer := make([]byte, bs)
		bitsRead, err := decryptDst.Read(decBuffer)
		if err != nil && !errors.Is(err, io.EOF) {
			t.Fatal(err)
		} else if errors.Is(err, io.EOF) {
			break
		}

		verifyChunk(decBuffer, offset, bitsRead)
		offset += int(math.Min(float64(bitsRead), float64(len(encryptSrc))-1))
	}

	integrityGot := sha256.Sum256(totalDecryptBuffer)

	if !reflect.DeepEqual(integrityGot, integrityOriginal) {
		t.Error("encryption integrity check failed.")
		t.Errorf("%v != %v", encryptSrc, totalDecryptBuffer)
	}
}
