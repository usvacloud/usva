package cryptography

import (
	"bytes"
	"crypto/rand"
	"testing"
)

func BenchmarkAES(b *testing.B) {
	sourceBuffer := make([]byte, (2<<19)*128) // allocate a buffer of 128mb
	if _, err := rand.Read(sourceBuffer); err != nil {
		b.Error(err)
	}

	writer := bytes.NewBuffer([]byte{})
	key, err := DeriveBasicKey([]byte("moikka!"), 16)
	if err != nil {
		b.Error(err)
	}

	b.StartTimer()
	if _, err = EncryptStream(writer, bytes.NewReader(sourceBuffer), key); err != nil {
		b.Error(err)
	}
	b.StopTimer()
}
